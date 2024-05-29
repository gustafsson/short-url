//go:build test
// +build test

package service

import (
	"io"
	"testing"
)

func TestBufferWriteCloser(t *testing.T) {
    bwc := newBufferWriteCloser()

    // Verify that bwc implements io.WriteCloser
    if _, ok := bwc.(io.WriteCloser); !ok {
        t.Errorf("newBufferWriteCloser() did not return an io.WriteCloser")
    }

    // Test the Close method
    err := bwc.Close()
    if err != nil {
        t.Errorf("Close() returned an error: %v", err)
    }
}
