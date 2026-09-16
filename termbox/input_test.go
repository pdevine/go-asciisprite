package termbox

// Parser tests: kitty forms + legacy interop + chunk-split regression
// (the wedge that wedged live terminals in the tcell patch).

import (
	"bytes"
	"testing"
	"time"
)

// drainOne runs the full drain pipeline over pending and returns events.
func drainOne(t *testing.T, s string, expire bool) []Event {
	t.Helper()
	inputMu.Lock()
	defer inputMu.Unlock()
	pending.Reset()
	pending.WriteString(s)
	// capture posts
	old := eventq
	eventq = make(chan Event, 64)
	defer func() { eventq = old }()
	drainLocked(expire)
	close(eventq)
	var evs []Event
	for ev := range eventq {
		evs = append(evs, ev)
	}
	return evs
}

func TestKittyPlainRuneNative(t *testing.T) {
	enhanced = true
	defer func() { enhanced = false }()
	evs := drainOne(t, "\x1b[97;1u", false)
	if len(evs) != 1 || evs[0].Ch != 'a' || evs[0].Type != EventKeyPress {
		t.Fatalf("%v", evs)
	}
}

func TestKittyCtrlANative(t *testing.T) {
	enhanced = true
	defer func() { enhanced = false }()
	evs := drainOne(t, "\x1b[97;5u", false)
	if len(evs) != 1 || evs[0].Key != KeyCtrlA {
		t.Fatalf("%v", evs)
	}
}

func TestKittyReleaseNative(t *testing.T) {
	enhanced = true
	defer func() { enhanced = false }()
	evs := drainOne(t, "\x1b[119;1:2u\x1b[119;1:3u", false)
	if len(evs) != 2 {
		t.Fatalf("%v", evs)
	}
	if evs[0].Type != EventKeyRepeat || evs[0].Ch != 'w' {
		t.Fatalf("%v", evs[0])
	}
	if evs[1].Type != EventKeyRelease || evs[1].Ch != 'w' {
		t.Fatalf("%v", evs[1])
	}
}

func TestKittyArrowsNative(t *testing.T) {
	enhanced = true
	defer func() { enhanced = false }()
	evs := drainOne(t, "\x1b[1;5A\x1b[3;5~", false)
	if len(evs) != 2 {
		t.Fatalf("%v", evs)
	}
	if evs[0].Key != KeyArrowUp || evs[0].Mod != ModCtrl {
		t.Fatalf("%v", evs[0])
	}
	if evs[1].Key != KeyDelete || evs[1].Mod != ModCtrl {
		t.Fatalf("%v", evs[1])
	}
}

// The wedge regression: a kitty prefix that never completes must drop
// its partial claim at expire so legacy drain delivers the bytes.
func TestKittyChunkedWedge(t *testing.T) {
	streams := []string{
		"\x1b[?11u\x1b[97;5u\x1b[1;5A\x1b[119;1:3u",
		"\x1b[<65;3;4M\x1b[97;5u\x1b[3;2~",
		"\x1b[27u\x1b[97:65;2;65u",
	}
	for _, tc := range streams {
		for split := 1; split < len(tc); split++ {
			enhanced = true
			func() {
				defer func() { enhanced = false }()
				done := make(chan struct{})
				go func() {
					drainOne(t, tc[:split], false)
					drainOne2(t, tc[split:], false)
					drainOne2(t, "", true)
					close(done)
				}()
				select {
				case <-done:
				case <-time.After(500 * time.Millisecond):
					t.Fatalf("WEDGE %q split@%d", tc, split)
				}
			}()
		}
	}
}

// drainOne2 appends to the existing pending buffer.
func drainOne2(t *testing.T, s string, expire bool) []Event {
	t.Helper()
	inputMu.Lock()
	defer inputMu.Unlock()
	pending.WriteString(s)
	old := eventq
	eventq = make(chan Event, 64)
	defer func() { eventq = old }()
	drainLocked(expire)
	close(eventq)
	var evs []Event
	for ev := range eventq {
		evs = append(evs, ev)
	}
	return evs
}

func TestLegacyKeys(t *testing.T) {
	cases := []struct {
		in   string
		want []Event
	}{
		{"a", []Event{{Type: EventKey, Ch: 'a'}}},
		{"\r", []Event{{Type: EventKey, Key: KeyEnter}}},
		{"\t", []Event{{Type: EventKey, Key: KeyTab}}},
		{"\x7f", []Event{{Type: EventKey, Key: KeyBackspace2}}},
		{"\x08", []Event{{Type: EventKey, Key: KeyBackspace}}},
		{"\x1b[A", []Event{{Type: EventKey, Key: KeyArrowUp}}},
		{"\x1bOA", []Event{{Type: EventKey, Key: KeyArrowUp}}},
		{"\x1b[3~", []Event{{Type: EventKey, Key: KeyDelete}}},
		{"\x1bab", []Event{{Type: EventKey, Ch: 'a', Mod: ModAlt}, {Type: EventKey, Ch: 'b'}}},
		{"\x03", []Event{{Type: EventKey, Key: KeyCtrlC}}},
	}
	for _, tc := range cases {
		got := drainOne(t, tc.in, true)
		if len(got) != len(tc.want) {
			t.Fatalf("%q: %v", tc.in, got)
		}
		for i, w := range tc.want {
			if got[i].Type != w.Type || got[i].Key != w.Key ||
				got[i].Ch != w.Ch || got[i].Mod != w.Mod {
				t.Fatalf("%q[%d]: got %v want %v", tc.in, i, got[i], w)
			}
		}
	}
}

func TestMouseSGR(t *testing.T) {
	evs := drainOne(t, "\x1b[<0;10;5M", true)
	if len(evs) != 1 || evs[0].Type != EventMouse || evs[0].MouseX != 9 || evs[0].MouseY != 4 {
		t.Fatalf("%v", evs)
	}
	var _ = bytes.MinRead
}

// A lone Esc pressed with no follow-up byte must NOT be delivered by the
// non-expire drain (it could still be an escape-sequence prefix) but MUST
// be delivered after the expiry timer fires.
func TestLoneEscNeedsExpire(t *testing.T) {
	evs := drainOne(t, "\x1b", false)
	if len(evs) != 0 {
		t.Fatalf("pre-expire: got %v", evs)
	}
	evs = drainOne2(t, "", true)
	if len(evs) != 1 || evs[0].Key != KeyEsc {
		t.Fatalf("post-expire: got %v", evs)
	}
}

// Esc followed quickly by another key = Alt+key, not Esc+key.
func TestEscPrefixIsAlt(t *testing.T) {
	evs := drainOne(t, "\x1bx", true)
	if len(evs) != 1 || evs[0].Ch != 'x' || evs[0].Mod&ModAlt == 0 {
		t.Fatalf("got %v", evs)
	}
}
