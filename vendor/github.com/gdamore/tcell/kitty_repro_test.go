// Repro harness: drives the REAL collectEventsFromInput, in enhanced
// mode, over interleaved kitty/mouse/reply/paste streams, exactly the
// shape that panicked in parseRune at the user's terminal.

package tcell

import (
	"bytes"
	"testing"
	"time"

	"github.com/gdamore/encoding"
	"github.com/gdamore/tcell/terminfo"
	_ "github.com/gdamore/tcell/terminfo/x/xterm_kitty"
)

func reproScreen() *tScreen {
	enc := encoding.UTF8
	// Use the real xterm-kitty terminfo: its key table pre-registers
	// "\x1b[1;2A"-style shifted-arrow codes and Mouse is enabled,
	// which is exactly the configuration that panicked live.
	ti, err := terminfo.LookupTerminfo("xterm-kitty")
	if err != nil {
		panic(err)
	}
	t := &tScreen{
		ti:       ti,
		decoder:  enc.NewDecoder(),
		mouse:    []byte(ti.Mouse),
		evch:     make(chan Event, 64),
		charset:  "UTF-8",
	}
	t.keyexist = make(map[Key]bool)
	t.keycodes = make(map[string]*tKeyCode)
	t.prepareKeys()
	t.enhancedKeys = 11 // as our demo requests
	return t
}

func drive(t *testing.T, t_ *tScreen, s string) []Event {
	buf := bytes.NewBufferString(s)
	return t_.collectEventsFromInput(buf, true) // expire=true path, matching the panic site
}

func TestKittyReproPanic(t *testing.T) {
	for _, tc := range []string{
		"\x1b[?11u\x1b[97u",                     // reply then keypress
		"\x1b[97u\x1b[<0;5;5M\x1b[98u",          // key, mouse, key
		"\x1b[<65;3;4M\x1b[97;5u",               // mouse then ctrl+a
		"\x1b[?11u\x1b[?1;2c\x1b[97u",           // reply, DA, keypress
		"\x1b[",                                  // bare alt-prefix in enhanced mode
		"\x1b[\x1b[97u",                          // two prefixes
		"\r\x1b[97u	",                           // legacy bytes between kitty events
		"\x1b[27u\x1b[27u\x1b[27u",               // repeated esc (kitty shape)
		"\x1b",                                   // lone esc
		"\x1b\x1b\x1b",                           // alt-esc
		"\x1b[200~pasted\x1b[201~\x1b[97u",       // bracketed paste + key
		"\x9b97u",                                // 8-bit CSI prefix
	} {
		t.Run("", func(t *testing.T) {
			s := reproScreen()
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("PANIC on %q: %v", tc, r)
				}
			}()
			evs := drive(t, s, tc)
			t.Logf("%q -> %d events", tc, len(evs))
		})
	}
}

// TestKittyChunked feeds byte streams split at every possible boundary,
// calling collectEventsFromInput(expire=false) per chunk followed by one
// expire=true pass — simulating inputLoop's 128-byte reads plus the
// 50ms keytimer expiry, which is the delivery shape live terminals use.
func TestKittyChunked(t *testing.T) {
	streams := []string{
		"\x1b[?11u\x1b[97;5u\x1b[1;5A\x1b[119;1:3u",
		"\x1b[<65;3;4M\x1b[97;5u\x1b[3;2~",
		"\x1b[27u\x1b[97:65;2;65u",
	}
	for _, tc := range streams {
		for split := 1; split < len(tc); split++ {
			s := reproScreen()
			buf := &bytes.Buffer{}
			var evs []Event
			func() {
				defer func() {
					if r := recover(); r != nil {
						t.Fatalf("PANIC %q split@%d: %v", tc, split, r)
					}
				}()
				step := func(tag string, expire bool) bool {
					done := make(chan struct{})
					go func() {
						defer close(done)
						evs = append(evs, s.collectEventsFromInput(buf, expire)...)
					}()
					select {
					case <-done:
						return true
					case <-time.After(500 * time.Millisecond):
						t.Fatalf("WEDGE in %s %q split@%d (buf %q)", tag, tc, split, buf.String())
						return false
					}
				}
				buf.WriteString(tc[:split])
				if !step("first", false) {
					return
				}
				buf.WriteString(tc[split:])
				if !step("second", false) {
					return
				}
				step("expire", true)
			}()
		}
	}
}
