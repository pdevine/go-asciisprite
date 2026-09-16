package termbox

// Public API.  Signatures and semantics preserved from the tcell-backed
// compatibility layer.

import (
	"errors"
	"os"
	"sync"
)

var errNotStarted = errors.New("termbox: not initialized")

var (
	tty     ttyState
	rend    *renderer
	front   screenBuf
	back    screenBuf
	outMode = Output256
	inited  bool
	mu      sync.Mutex

	eventq   = make(chan Event, 64)
	enhanced bool
)

// Init initializes the screen for use.
func Init() error {
	mu.Lock()
	defer mu.Unlock()
	enhanced = false
	return initLocked()
}

// InitEnhancedKeys initializes the screen and requests enhanced (kitty
// protocol) keyboard reporting: disambiguated escape codes, key
// repeat/release events, and press/release reporting for text keys.
// In enhanced mode key events arrive as EventKeyPress, EventKeyRepeat,
// and EventKeyRelease instead of EventKey.  The return value is the set
// of enhancement flags granted by the terminal; 0 means unsupported and
// the screen runs exactly as if Init had been called.
func InitEnhancedKeys() (int, error) {
	mu.Lock()
	defer mu.Unlock()
	enhanced = false
	if err := initLocked(); err != nil {
		return 0, err
	}
	flags := enableKittyLocked()
	enhanced = flags != 0
	return flags, nil
}

func initLocked() error {
	if err := tty.start(); err != nil {
		return err
	}
	if p := os.Getenv("TERMBOX_DEBUG"); p != "" {
		if f, err := os.Create(p); err == nil {
			debugLog = f
		}
	}
	w, h, err := tty.size()
	if err != nil {
		tty.stop()
		return err
	}
	front = newScreenBuf(w, h)
	back = newScreenBuf(w, h)
	rend = newRenderer(tty.out)
	// enter alt screen, hide cursor, clear, enable mouse (SGR + button)
	io_seq := "\x1b[?1049h\x1b[?25l\x1b[2J\x1b[?1002h\x1b[?1006h"
	tty.out.WriteString(io_seq)
	inited = true
	go inputMain()
	return nil
}

// Close cleans up the terminal, restoring the previous screen, cursor,
// colors, keyboard mode, and termios state.  The exit-sequence ordering
// (reset → cursor → mouse off → enhanced-keys off → leave alt screen)
// was derived from live testing: reorder at your peril.
func Close() {
	mu.Lock()
	if !inited {
		mu.Unlock()
		return
	}
	inited = false
	stopInput()
	seq := "\x1b[0m\x1b[?25h\x1b[?1006l\x1b[?1002l"
	if enhanced {
		seq += "\x1b[=0;1u\x1b[<u"
	}
	seq += "\x1b[?1049l"
	tty.out.WriteString(seq)
	tty.out.Sync()
	enhanced = false
	mu.Unlock()
	// settle before restoring termios, so the terminal has processed the
	// teardown bytes before input modes change underneath it
	settle()
	tty.stop()
}

// Flush updates the screen.
func Flush() error {
	mu.Lock()
	defer mu.Unlock()
	if !inited {
		return errNotStarted
	}
	rend.render(&front, &back)
	return nil
}

// Sync forces a full redraw.
func Sync() error {
	mu.Lock()
	defer mu.Unlock()
	if !inited {
		return errNotStarted
	}
	front.clear() // force every cell dirty
	rend.render(&front, &back)
	return nil
}

// SetCursor displays the terminal cursor at the given location.
func SetCursor(x, y int) {
	mu.Lock()
	defer mu.Unlock()
	if !inited {
		return
	}
	if x < 0 || y < 0 {
		tty.out.WriteString("\x1b[?25l")
		return
	}
	rend.cup(x, y)
	tty.out.WriteString("\x1b[?25h")
}

// HideCursor hides the terminal cursor.
func HideCursor() { SetCursor(-1, -1) }

// Size returns the screen size as width, height in character cells.
func Size() (int, int) {
	mu.Lock()
	defer mu.Unlock()
	if !inited {
		return 0, 0
	}
	return back.width, back.height
}

// SetCell sets the content of a cell in the back buffer.
func SetCell(x, y int, ch rune, fg, bg Attribute) {
	mu.Lock()
	defer mu.Unlock()
	if !inited || x < 0 || y < 0 || x >= back.width || y >= back.height {
		return
	}
	back.cells[y*back.width+x] = cell{ch: ch, fg: fg, bg: bg}
}

// GetCell returns the content of a cell from the back buffer.
func GetCell(x, y int) (rune, Attribute, Attribute) {
	mu.Lock()
	defer mu.Unlock()
	if !inited || x < 0 || y < 0 || x >= back.width || y >= back.height {
		return 0, ColorDefault, ColorDefault
	}
	c := back.cells[y*back.width+x]
	return c.ch, c.fg, c.bg
}

// Clear clears the back buffer with the given attributes.
func Clear(fg, bg Attribute) {
	mu.Lock()
	defer mu.Unlock()
	if !inited {
		return
	}
	for i := range back.cells {
		back.cells[i] = cell{ch: ' ', fg: fg, bg: bg}
	}
}

// PollEvent blocks until an event is ready.  On a consumer panic — or
// anywhere up the consumer's stack — the deferred restore here puts the
// terminal back before the panic continues, so a crashed app cannot
// leave the shell in raw mode.
func PollEvent() (ev Event) {
	defer func() {
		if r := recover(); r != nil {
			restoreTerminal()
			panic(r)
		}
	}()
	return pollEvent()
}

// Interrupt posts an interrupt event.
func Interrupt() {
	eventq <- Event{Type: EventInterrupt}
}

// SetInputMode toggles mouse input (the only mode with effect).
func SetInputMode(mode InputMode) InputMode {
	mu.Lock()
	defer mu.Unlock()
	if !inited {
		return InputCurrent
	}
	if mode&InputMouse != 0 {
		tty.out.WriteString("\x1b[?1002h\x1b[?1006h")
	} else {
		tty.out.WriteString("\x1b[?1006l\x1b[?1002l")
	}
	return InputEsc
}

// SetOutputMode sets the color output mode.
func SetOutputMode(mode OutputMode) OutputMode {
	mu.Lock()
	defer mu.Unlock()
	switch mode {
	case OutputCurrent:
		return outMode
	case OutputNormal, Output256, Output216, OutputGrayscale:
		outMode = mode
		rend.outMode = mode
		return mode
	default:
		return OutputCurrent
	}
}

// ParseEvent is not supported (raw byte parsing lives in the input layer).
func ParseEvent(data []byte) Event {
	return Event{Type: EventError, Err: errors.New("no raw events")}
}

// PollRawEvent is not supported.
func PollRawEvent(data []byte) Event {
	return Event{Type: EventError, Err: errors.New("no raw events")}
}

// restoreTerminal is the panic-path counterpart of Close: idempotent,
// lock-free (we may be panicking with the mutex held), best-effort.
func restoreTerminal() {
	tty.out.WriteString("\x1b[0m\x1b[?25h\x1b[?1006l\x1b[?1002l\x1b[=0;1u\x1b[<u\x1b[?1049l")
	tty.out.Sync()
	tty.stop()
}
