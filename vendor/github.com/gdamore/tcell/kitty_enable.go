// Kitty keyboard protocol support — in-tree patch for go-asciisprite.
// Terminal detection and enable/disable plumbing for tScreen; implements
// the Screen.EnableEnhancedKeys/DisableEnhancedKeys API per
// https://sw.kovidgoyal.net/kitty/keyboard-protocol/

package tcell

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// Progressive enhancement flags for the kitty keyboard protocol, for use
// with Screen.EnableEnhancedKeys.
const (
	// KbdEnhDisambiguate makes Esc/alt/ctrl key combos arrive as
	// unambiguous CSI u sequences instead of legacy overlapping codes.
	KbdEnhDisambiguate = 1
	// KbdEnhEventTypes adds key repeat and key release events.
	KbdEnhEventTypes = 2
	// KbdEnhAlternateKeys reports shifted/base-layout key codepoints.
	KbdEnhAlternateKeys = 4
	// KbdEnhAllKeys reports text-producing keys as escape codes too,
	// enabling press/repeat/release for letter keys (game input).
	KbdEnhAllKeys = 8
	// KbdEnhAssociatedText embeds text codepoints in key events.
	KbdEnhAssociatedText = 16
)

// kittyLikelyUnsupported short-circuits the query on terminals known to
// lack the protocol, avoiding the detection round-trip timeout.
func kittyLikelyUnsupported() bool {
	term := os.Getenv("TERM")
	if strings.Contains(term, "kitty") || strings.Contains(term, "foot") ||
		strings.Contains(term, "wezterm") || strings.Contains(term, "ghostty") ||
		strings.Contains(term, "alacritty") || strings.Contains(term, "rio") ||
		strings.Contains(term, "contour") {
		return false
	}
	switch os.Getenv("TERM_PROGRAM") {
	case "iTerm.app", "WezTerm", "vscode", "ghostty", "Apple_Terminal":
		// Apple_Terminal definitely does not; the others do (modern
		// versions of iTerm2 and WezTerm implement the protocol).
		if os.Getenv("TERM_PROGRAM") == "Apple_Terminal" {
			return true
		}
		return false
	}
	// unknown terminal: let the query decide
	return false
}

// EnableEnhancedKeys implements Screen.EnableEnhancedKeys.
func (t *tScreen) EnableEnhancedKeys(flags int) int {
	if kittyLikelyUnsupported() {
		return 0
	}

	t.Lock()
	t.kittyReplyCh = make(chan int, 1)
	// push current mode, set requested flags, then query what took
	t.TPuts(fmt.Sprintf("\x1b[>%du", flags))
	t.TPuts(fmt.Sprintf("\x1b[=%d;1u", flags))
	t.TPuts("\x1b[?u")
	t.Unlock()

	// The query reply arrives on stdin and is consumed by
	// parseKittyReply inside the input parser, which signals us here.
	select {
	case granted := <-t.kittyReplyCh:
		if granted&flags == flags {
			t.Lock()
			t.enhancedKeys = flags
			t.Unlock()
			return granted
		}
		// terminal spoke the protocol but implemented only a subset;
		// restore and report legacy mode
		t.TPuts("\x1b[<u")
		return 0
	case <-time.After(300 * time.Millisecond):
		// no reply: unsupported; pop the mode we pushed
		t.TPuts("\x1b[<u")
		return 0
	}
}

// DisableEnhancedKeys implements Screen.DisableEnhancedKeys.
// The pop is written raw (bypassing TPuts, which is only for terminfo
// strings) so it is emitted exactly, and a small delay gives a slow
// terminal time to process it before the screen teardown bytes follow.
func (t *tScreen) DisableEnhancedKeys() {
	t.Lock()
	defer t.Unlock()
	if t.enhancedKeys != 0 {
		t.enhancedKeys = 0
		// Explicitly reset all flags first (no stack involvement), then
		// pop the mode we pushed at enable time.  Sending both makes the
		// restore robust against terminals whose push/pop stack semantics
		// are quirky or which interleave the pop with teardown codes.
		io.WriteString(t.out, "\x1b[=0;1u")
		io.WriteString(t.out, "\x1b[<u")
		time.Sleep(50 * time.Millisecond)
	}
}
