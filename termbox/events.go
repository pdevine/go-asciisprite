// Package termbox is the terminal backend for go-asciisprite: a
// self-contained replacement for the vendored tcell/nsf-termbox stack
// this package used to shim.  It owns the tty, renders a cell diff to
// an xterm-compatible escape stream, and parses input inline — with
// first-class kitty keyboard protocol support via InitEnhancedKeys.
// Platforms: macOS, Linux; Windows fails cleanly (native console backend
// planned).  Debug: set TERMBOX_DEBUG=<path> at Init to trace input.

package termbox

// Event represents an event like a key press, mouse action, or window resize.
type Event struct {
	Type   EventType
	Mod    Modifier
	Key    Key
	Ch     rune
	Width  int
	Height int
	Err    error
	MouseX int
	MouseY int
	N      int
}

// EventType discriminates the kind of an Event.
type EventType uint8

// Event types.
const (
	EventNone EventType = iota
	// EventKey is a key press delivered in legacy mode (Init).  In
	// enhanced mode (InitEnhancedKeys) key events arrive as
	// EventKeyPress/EventKeyRepeat/EventKeyRelease and EventKey is
	// never emitted.
	EventKey
	EventResize
	EventMouse
	EventInterrupt
	EventError
	EventRaw
	// EventKeyPress, EventKeyRepeat, and EventKeyRelease are the
	// enhanced-mode key events (kitty keyboard protocol on unix;
	// native console input on Windows).
	EventKeyPress
	EventKeyRepeat
	EventKeyRelease
)

// Cell represents a single character cell on screen.
type Cell struct {
	Ch rune
	Fg Attribute
	Bg Attribute
}
