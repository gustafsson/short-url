package service

import (
	"bytes"
	"io"
)

type bufferWriteCloser struct {
    *bytes.Buffer
}

func (bwc *bufferWriteCloser) Close() error {
    // No action needed for close, but method is required for io.WriteCloser
    return nil
}

func newBufferWriteCloser() io.WriteCloser {
    return &bufferWriteCloser{Buffer: new(bytes.Buffer)}
}
