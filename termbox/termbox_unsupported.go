//go:build !darwin && !linux

package termbox

// Unsupported platforms: Init fails cleanly.  Windows support lands as
// a dedicated backend (console records + VT output); see the plan.

import (
	"errors"
	"io"
)

var errUnsupported = errors.New("termbox: platform not supported by the native backend")

type ttyWriter struct{}

func (ttyWriter) Write(p []byte) (int, error)  { return 0, errUnsupported }
func (ttyWriter) Sync() error                  { return errUnsupported }
func (ttyWriter) WriteString(s string) (int, error) {
	return 0, errUnsupported
}

type ttyState struct {
	in  *osFileStub
	out ttyWriter
}

type osFileStub struct{}

func (t *ttyState) start() error            { return errUnsupported }
func (t *ttyState) stop()                   {}
func (t *ttyState) size() (int, int, error) { return 0, 0, errUnsupported }

// inputReadLoop never runs because start() fails first.
func inputReadLoop(f *osFileStub, quit chan struct{}) error { return errUnsupported }

var _ io.Writer = ttyWriter{}
