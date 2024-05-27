# Position recorder

This package contains a utility for Aleo oracle encoding package. The goal of this utility is to record the start positions of
every write and the length of data that was written. The positions and the length is counted in blocks of specified number of bytes.

This package provides an interface `PositionRecorder`, which implements `io.Writer` interface and adds `GetLastWrite` method,
which returns information about the last write operation.

## Usage

Create a recorder using `NewPositionRecorder` function. It will panic if the block size is odd.
`Write` function on the recorder acts as a proxy for the underlying `io.Writer` and records the information about the writes in
its internal state.

**Important!** `Write` function will return a `ErrDataAlignment` error if the argument's length is not aligned to the block size.

```golang
import "bytes"

// it implements io.Writer so it works for the example
var buf bytes.Buffer
// create a recorder, pass an underlying writer and block size
recorder := NewPositionRecorder(&buf, 16)

exampleBlock := make([]byte, 16)
exampleBlockDouble := make([]byte, 32)

info := recorder.GetLastWrite()
// info is nil since there were no Write operations yet

// Write conforms to io.Writer interface
n, err := recorder.Write(exampleBlock)
info = recorder.GetLastWrite()
// info.Pos = 0
// info.Len = 1

// write 2 blocks in the next operation
recorder.Write(exampleBlockDouble)
info = recorder.GetLastWrite()
// info.Pos = 1
// info.Len = 2

// write another one block
n, err := recorder.Write(exampleBlock)
info = recorder.GetLastWrite()
// info.Pos = 3
// info.Len = 1

```
