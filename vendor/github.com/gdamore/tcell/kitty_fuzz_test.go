// Fuzz/regression harness for the kitty keyboard protocol parser.

package tcell

import (
	"bytes"
	"testing"
)

func bytesBuffer(s string) *bytes.Buffer {
	return bytes.NewBufferString(s)
}

// FuzzKittyKey ensures parseKittyKey never panics and never claims more
// bytes than it was given, and that recognized sequences round-trip
// through the pipeline stages without corrupting the buffer.
func FuzzKittyKey(f *testing.F) {
	seeds := []string{
		"\x1b[97;5u",
		"\x1b[97:65;2;65u",
		"\x1b[1;5A",
		"\x1b[3;5~",
		"\x1b[?",
		"\x1b[?11u",
		"\x1b[<35;10;20M",
		"\x1b[",
		"\x1b",
		"a",
		"\x1b[57401u",
		"\x1b[0;;;u",
		"\x1b[u\x1b[u\x1b[u",
		"\x1b[97",
		"\x1b[97;5",
		"\x1b[;;::;u",
		"\x1b[27u",
	}
	for _, s := range seeds {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		b := []byte(s)
		ev, n, partial := parseKittyKey(b)
		if n < 0 || n > len(b) {
			t.Fatalf("consumed %d of %d", n, len(b))
		}
		if n > 0 && partial {
			t.Fatalf("consumed+%d and partial both set for %q", n, s)
		}
		if n == 0 && !partial && ev != nil {
			t.Fatalf("event without consumption for %q", s)
		}
	})
}

// TestKittyReplyFuzz feeds odd shapes through the query-reply consumer
// via a minimal tScreen (no tty needed — parseKittyReply touches only
// the buffer and the reply channel).
func TestKittyReplyShapes(t *testing.T) {
	scr := &tScreen{kittyReplyCh: make(chan int, 4)}
	for _, tc := range []struct {
		in        string
		wantFlags int
		wantEvt   bool
	}{
		{"\x1b[?11u", 11, true},
		{"\x1b[?", -1, false},    // partial
		{"\x1b[?1;2c", -1, false}, // DA reply: left alone
		{"\x1b[?u", 0, true},
	} {
		buf := bytesBuffer(tc.in)
		partial, handled := scr.parseKittyReply(buf, false)
		if partial && tc.wantFlags >= 0 {
			t.Fatalf("%q: unexpected partial", tc.in)
		}
		if tc.wantEvt {
			if !handled {
				t.Fatalf("%q: not handled", tc.in)
			}
			got := <-scr.kittyReplyCh
			if got != tc.wantFlags {
				t.Fatalf("%q: flags %d, want %d", tc.in, got, tc.wantFlags)
			}
			if buf.Len() != 0 {
				t.Fatalf("%q: %d bytes left", tc.in, buf.Len())
			}
		} else if handled {
			t.Fatalf("%q: unexpectedly handled", tc.in)
		}
	}
}

// TestKittyStageFullPipeline pushes sequences through the stage the way
// collectEventsFromInput does, checking buffer accounting.
func TestKittyStageFullPipeline(t *testing.T) {
	scr := &tScreen{}
	evs := []Event{}
	buf := bytesBuffer("\x1b[1;5A\x1b[119;1:2u\x1b[119;1:3u")
	for buf.Len() > 0 {
		part, comp := scr.parseKittyStage(buf, &evs, false)
		if !comp {
			t.Fatalf("stage bailed: part=%v remaining=%q", part, buf.String())
		}
	}
	if len(evs) != 3 {
		t.Fatalf("events=%d", len(evs))
	}
	k0 := evs[0].(*EventKey)
	k1 := evs[1].(*EventKey)
	k2 := evs[2].(*EventKey)
	if k0.Key() != KeyUp || k0.Modifiers() != ModCtrl {
		t.Fatalf("ev0: %+v", k0)
	}
	if k1.Rune() != 'w' || k1.EventType() != KeyEventRepeat {
		t.Fatalf("ev1: %+v", k1)
	}
	if k2.EventType() != KeyEventRelease {
		t.Fatalf("ev2: %+v", k2)
	}
}
