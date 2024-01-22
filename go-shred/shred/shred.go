package shred

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

// Defines a generic Shredder
type Shredder interface {
	// Retreive the current shredder chunk size
	ChunkSize() int
	// Retreive a randomness generator
	BytesBuffer() io.Reader
}

// RandomShredder is a Shredder that takes its input data from /dev/urandom
type RandomShredder struct {
	chunkSize int
	randGen   io.Reader
	urandDev  *os.File
}

func (s *RandomShredder) ChunkSize() int {
	return s.chunkSize
}

func (s *RandomShredder) BytesBuffer() io.Reader {
	return s.randGen
}

// NewRandomShredder creates a new shredding configuration with a chunksize of 4096
// and /dev/urandom as RNG source.
func NewRandomShredder() (*RandomShredder, error) {
	// Open /dev/urandom as a buffered reader. This will generate our random bytes.
	urandFd, err := os.Open("/dev/urandom")
	if err != nil {
		return nil, fmt.Errorf("error opening /dev/urandom")
	}
	randSource := bufio.NewReader(urandFd)
	return &RandomShredder{
		chunkSize: 4096,
		randGen:   randSource,
		urandDev:  urandFd,
	}, nil
}

// Close closes the /dev/urandom device
func (s *RandomShredder) Close() {
	s.urandDev.Close()
}

// Shred takes a filepath as input and overwrites its content 3 times, then deletes it.
func Shred(path string) error {
	// Extract file info. Nice side-effect is that this also throws an error if the file
	// doesn't exist.
	f, err := os.OpenFile(path, os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error when opening file %s: %s", path, err)
	}
	defer f.Close()

	// Create a shreder s that uses the Linux /dev/urandom special file
	s, err := NewRandomShredder()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	// Do the shredding operation 3 times.
	for i := 0; i < 3; i++ {
		log.Printf("Shredding... %d\n", i+1)
		ShredWithShredder(f, s)
	}
	f.Close()
	// Finally, as per requirement, remove the file
	log.Println("Done shredding. Removing file.")
	err = os.Remove(path)
	if err != nil {
		panic(err)
	}

	return nil
}

// ShredWithShredder shredds the given file with using a custom Shredder type
// object for added flexibility.
func ShredWithShredder(f *os.File, s Shredder) error {
	info, err := f.Stat()
	if err != nil {
		return err
	}
	fileSize := int(info.Size())

	err = shredRawWithShredder(f, fileSize, s)
	if err != nil {
		return err
	}

	return nil
}

// The raw shredding function. Operates on I/O traits, so should work on more than files.
// This is where the magic happens.
func shredRawWithShredder(dest io.WriteSeeker, destSize int, s Shredder) error {
	chunkSize := int(s.ChunkSize())
	src := s.BytesBuffer()

	// The algorithm is simple. I have split the shredding in chunks to balance
	// memory consumption and speed.
	chunk := make([]byte, chunkSize)
	noChunks := (destSize + chunkSize - 1) / chunkSize

	// 1. Go to beginning of the file/buffer
	dest.Seek(0, 0)

	// 2. Overwrite all bytes with randomness or some other data
	for i := 0; i < noChunks-1; i++ {
		_, err := src.Read(chunk)
		if err != nil {
			return err
		}
		_, err = dest.Write(chunk)
		if err != nil {
			return err
		}
	}

	// 3. The last chunk might be of a size less than chunkSize. Handle that
	lastChunkSize := destSize - (noChunks-1)*chunkSize
	_, err := src.Read(chunk)
	_, err = dest.Write(chunk[:lastChunkSize])
	if err != nil {
		return err
	}

	return nil
}
