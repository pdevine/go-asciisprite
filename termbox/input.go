package termbox

// Input layer: read goroutine, parser loop, event queue.

import (
	"sync"
	"time"
	"unicode/utf8"
)

var inputMu sync.Mutex

// post enqueues an event, dropping on full queue (input bursts can't
// block the parser).
func post(ev Event) {
	select {
	case eventq <- ev:
	default:
	}
}

// pollEvent blocks until an event is ready; called by PollEvent, which
// owns the panic-restore wrapper.
func pollEvent() Event {
	return <-eventq
}

// settle gives the terminal a moment to process teardown bytes before
// termios is restored.
func settle() { time.Sleep(50 * time.Millisecond) }

var (
	inputQuit  chan struct{}
	inputDoneq chan struct{}
)

// inputMain is the input goroutine: reads chunks, parses events, posts
// them.  Split into platform halves: readUnix vs record translation on
// Windows.
func inputMain() {
	inputQuit = make(chan struct{})
	inputDoneq = make(chan struct{})
	defer close(inputDoneq)
	_ = inputReadLoop(tty.in, inputQuit)
}

func stopInput() {
	if inputQuit != nil {
		close(inputQuit)
		select {
		case <-inputDoneq:
		case <-time.After(200 * time.Millisecond):
		}
	}
}

// decodeRune decodes one UTF-8 rune from b; returns (r, size).
func decodeRune(b []byte) (rune, int) {
	r, n := utf8.DecodeRune(b)
	return r, n
}
