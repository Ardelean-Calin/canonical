The crux of this exercise is in the `shred` package. However, I have included a `main.go` which simply takes the first argument as filepath and shredds that file.
NOTE: I purposely didn't handle errors related to input sanitization, as it didn't fit in the scope and time-frame of this exercise.

Two functions are exported:
- **Shred** for shredding a file specified by the given path (assumes path exists and is a valid file)
- **ShredWithShredder** which is a more general implementation of shredding on io streams (can be files or not). Also takes a Shredder interface as argument as it allowed me to test shredding easier and allows for customization in shredding behavior.

To run shredding with a file, simply call `go run main.go <file_to_be_shredded>` from this directory. For example:
```
  go run main.go ~/Downloads/ubuntu-22.04.3-desktop-amd64.iso
```
or build the program with `go build` and then run `go-shred <file>`.

A `coverage.html` file was included to prove test coverage of the most critical code paths. Due to time constraint, some erronous behavior was not covered.
