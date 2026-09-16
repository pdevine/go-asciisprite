package termbox

// Input layer: read goroutine, parser loop, event queue.

import (
	"io"
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
	inputQuit   chan struct{}
	inputDoneq  chan struct{}
	resizeQuit  chan struct{}
	resizeDoneq chan struct{}
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

// resizeMain translates SIGWINCH into EventResize.  The channel is
// created by tty.start on unix; it is nil on platforms with native
// resize events (Windows), in which case this goroutine waits on quit
// only.
func resizeMain() {
	resizeQuit = make(chan struct{})
	resizeDoneq = make(chan struct{})
	defer close(resizeDoneq)
	for {
		select {
		case <-tty.winch:
			w, h, err := tty.size()
			if err != nil {
				continue
			}
			mu.Lock()
			if inited && (w != back.width || h != back.height) {
				front = newScreenBuf(w, h)
				back = newScreenBuf(w, h)
				io.WriteString(tty.out, "\x1b[2J")
				post(Event{Type: EventResize, Width: w, Height: h})
			}
			mu.Unlock()
		case <-resizeQuit:
			return
		}
	}
}

func stopResize() {
	if resizeQuit != nil {
		close(resizeQuit)
		select {
		case <-resizeDoneq:
		case <-time.After(200 * time.Millisecond):
		}
	}
}

// decodeRune decodes one UTF-8 rune from b; returns (r, size).
func decodeRune(b []byte) (rune, int) {
	r, n := utf8.DecodeRune(b)
	return r, n
}
