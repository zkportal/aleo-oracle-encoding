package positionRecorder

import (
	"errors"
	"io"
)

var (
	ErrDataAlignment = errors.New("data is not aligned to block size")
)

type PositionInfo struct {
	// Index of the block where the write operation started
	Pos int
	// Number of blocks written in the write operation
	Len int
}

// PositionRecorder records the position and number of blocks written to the underlying data stream.
type PositionRecorder interface {
	io.Writer

	GetLastWrite() *PositionInfo
}

// PositionRecordingProxy is a wrapper around a data stream, which follows io.Writer interface and records
// positional information about the last write operation.
type PositionRecordingProxy struct {
	PositionRecorder

	writer    io.Writer
	blockSize int
	lastWrite *PositionInfo
}

// NewPositionRecorder creates a new position recorder. Block size must be an even number.
func NewPositionRecorder(writer io.Writer, blockSize int) *PositionRecordingProxy {
	if blockSize%2 != 0 {
		panic("block size must be an even number")
	}

	return &PositionRecordingProxy{
		writer:    writer,
		blockSize: blockSize,
		lastWrite: nil,
	}
}

// Writes p to the underlying writer and records successful writes. Returned values are io.Writer.Write return values.
// The information about the last write operation can be obtained using GetLastWrite.
func (r *PositionRecordingProxy) Write(p []byte) (n int, err error) {
	length := len(p)
	if length%r.blockSize != 0 {
		return 0, ErrDataAlignment
	}

	numBlocks := length / r.blockSize

	// write to writer to make sure we can write all of the data to make sure only full writes are recorded
	n, err = r.writer.Write(p)
	if err != nil || n != length {
		return n, err
	}

	// get the position where the current write starts
	lastPos := 0
	if r.lastWrite != nil {
		lastPos = r.lastWrite.Pos + r.lastWrite.Len
	}

	r.lastWrite = &PositionInfo{
		Pos: lastPos,
		Len: numBlocks,
	}

	return
}

// GetLastWrite returns information about the last complete write operation. The information is discarded and replaced
// every time Write was executed successfully
func (r *PositionRecordingProxy) GetLastWrite() *PositionInfo {
	return r.lastWrite
}
