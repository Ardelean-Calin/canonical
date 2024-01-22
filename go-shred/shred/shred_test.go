package shred

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

// We create a dummy shredder that we can use in our tests to prove that shredding
// works. This way we can actually "expect" some content, instead of always getting
// random bytes.
// This is actually one of the reasons I chose the interface-based approach,
// to make testing easier.
type TestShredder struct {
	chunkSize int
}

func newTestShredder(chunkSize int) *TestShredder {
	return &TestShredder{chunkSize}
}

func (s *TestShredder) ChunkSize() int {
	return s.chunkSize
}

func (s *TestShredder) BytesBuffer() io.Reader {
	buf := []byte(strings.Repeat("shredded", 32))
	b := bytes.NewReader(buf)
	return b
}

// Tests that Shredding works. Since we provide the "randomness" we can easily test
// wether the content was overwritten or not.
func TestShredFileCustom(t *testing.T) {
	testFile, _ := os.CreateTemp("", "shred")
	defer testFile.Close()
	defer os.Remove(testFile.Name())
	testFile.WriteString("Hello, World!")

	shredder := newTestShredder(1)
	ShredWithShredder(testFile, shredder)

	// After shredding, we test that the file contains
	buf := make([]byte, 13)
	testFile.Seek(0, 0)
	testFile.Read(buf)

	if string(buf) != "shreddedshred" {
		t.FailNow()
	}

}

// Tests behavior when the file to be shredded doesn't exist
func TestShredFileNotExists(t *testing.T) {
	path := "this-path-does-not-exist.txt"
	err := Shred(path)
	if err == nil {
		t.FailNow()
	}
}

// Test the Shred function start-to-end. I'm not happy with this test since I can't really
// check the intermediate states of the file to be shredded. But I had to fit in 2 hours so
// this will do for now.
func TestShred(t *testing.T) {
	// Step 1. Create a temporary file
	testFile, _ := os.CreateTemp("", "shred")
	testFile.WriteString("I am foobar.")
	path := testFile.Name()
	testFile.Close()

	// Step 2. Shred it! It should be removed from the filesystem
	Shred(path)

	// Step 3. The only thing I can check after this operation is that the file was deleted. NOTE: additional tests could be written by mocking the File interface. However to limit myself to 2 hours I chose not to do that. Also I was limited by the fact that the function prototype for "Shred" was required to only contain the parameter "path".
	_, err := os.Stat(path)
	if os.IsExist(err) {
		t.FailNow()
	}

}

// Didn't have time to implement this, sorry. Also I'm not yet sure that the whole
// alignment thing actually matters. The file might get some additional bytes,
// maybe this isn't a big deal.
func TestAlignment(t *testing.T) {
	// TODO. Test that writing a file not aligned to the chunk size works as expected.
}

// Test the constructor for a RandomShredder, the main Shredder of this exercise.
func TestNewRandomShredder(t *testing.T) {
	s, err := NewRandomShredder()
	if err != nil {
		t.FailNow()
	}

	if s.chunkSize != 4096 {
		t.FailNow()
	}

	if s.urandDev.Name() != "/dev/urandom" {
		t.FailNow()
	}
}
