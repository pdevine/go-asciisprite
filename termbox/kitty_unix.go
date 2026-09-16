package termbox

// Kitty enable/disable on unix: push flags, verify via query reply, pop
// on close.  Ported from the proven tcell patch.

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// kittyLikelyUnsupported short-circuits the query on terminals known to
// lack the protocol, avoiding the detection round-trip timeout.
func kittyLikelyUnsupported() bool {
	term := os.Getenv("TERM")
	for _, s := range []string{"kitty", "foot", "wezterm", "ghostty", "alacritty", "rio", "contour"} {
		if strings.Contains(term, s) {
			return false
		}
	}
	if os.Getenv("TERM_PROGRAM") == "Apple_Terminal" {
		return true
	}
	return false
}

var kittyReply = make(chan int, 1)

// kittyFlags records the terminal's reply flags and wakes detection.
func kittyFlags(flags int) {
	select {
	case kittyReply <- flags:
	default:
	}
}

// enableKittyLocked requests enhancements.  Caller holds mu.
func enableKittyLocked() int {
	if kittyLikelyUnsupported() {
		return 0
	}
	const flags = 0b1011 // disambiguate | event types | all keys
	tty.out.WriteString(fmt.Sprintf("\x1b[>%du\x1b[=%d;1u\x1b[?u", flags, flags))
	select {
	case granted := <-kittyReply:
		if granted&flags == flags {
			return flags
		}
		tty.out.WriteString("\x1b[<u")
		return 0
	case <-time.After(300 * time.Millisecond):
		tty.out.WriteString("\x1b[<u")
		return 0
	}
}
