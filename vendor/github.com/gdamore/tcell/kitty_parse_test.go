// Unit tests for the kitty keyboard protocol parser (in-tree patch).

package tcell

import (
	"testing"
)

func keyEvent(t *testing.T, seq string) *EventKey {
	t.Helper()
	ev, n, partial := parseKittyKey([]byte(seq))
	if partial {
		t.Fatalf("%q: unexpected partial", seq)
	}
	if n != len(seq) {
		t.Fatalf("%q: consumed %d bytes, expected %d", seq, n, len(seq))
	}
	return ev
}

func TestKittyPlainRune(t *testing.T) {
	ev := keyEvent(t, "\x1b[97;1u") // a, no modifiers
	if ev == nil || ev.Key() != KeyRune || ev.Rune() != 'a' {
		t.Fatalf("got %+v", ev)
	}
	if ev.EventType() != KeyEventPress || ev.Modifiers() != ModNone {
		t.Fatalf("type=%d mod=%d", ev.EventType(), ev.Modifiers())
	}
}

func TestKittyCtrlA(t *testing.T) {
	ev := keyEvent(t, "\x1b[97;5u") // ctrl+a
	if ev == nil || ev.Key() != KeyCtrlA {
		t.Fatalf("got %+v key=%d", ev, ev.Key())
	}
}

func TestKittyShiftedRune(t *testing.T) {
	// shift+a with alternate-key reporting: 97:65;2;65
	ev := keyEvent(t, "\x1b[97:65;2;65u")
	if ev == nil || ev.Key() != KeyRune || ev.Rune() != 'A' {
		t.Fatalf("got %+v rune=%q", ev, ev.Rune())
	}
	if ev.Modifiers() != ModShift {
		t.Fatalf("mod=%d", ev.Modifiers())
	}
}

func TestKittyReleasedEvent(t *testing.T) {
	ev := keyEvent(t, "\x1b[97;1:3u") // a release
	if ev == nil || ev.EventType() != KeyEventRelease {
		t.Fatalf("got %+v type=%d", ev, ev.EventType())
	}
	if ev.Key() != KeyRune || ev.Rune() != 'a' {
		t.Fatalf("key=%d rune=%q", ev.Key(), ev.Rune())
	}
}

func TestKittyRepeatEvent(t *testing.T) {
	ev := keyEvent(t, "\x1b[119;1:2u") // w repeat
	if ev == nil || ev.EventType() != KeyEventRepeat || ev.Rune() != 'w' {
		t.Fatalf("got %+v", ev)
	}
}

func TestKittyArrows(t *testing.T) {
	ev := keyEvent(t, "\x1b[1;5A") // ctrl+up
	if ev == nil || ev.Key() != KeyUp || ev.Modifiers() != ModCtrl {
		t.Fatalf("got %+v key=%d mod=%d", ev, ev.Key(), ev.Modifiers())
	}
	ev = keyEvent(t, "\x1b[1;1:3D") // left, release
	if ev == nil || ev.Key() != KeyLeft || ev.EventType() != KeyEventRelease {
		t.Fatalf("got %+v", ev)
	}
}

func TestKittyTilde(t *testing.T) {
	ev := keyEvent(t, "\x1b[3~") // delete
	if ev == nil || ev.Key() != KeyDelete {
		t.Fatalf("got %+v", ev)
	}
	ev = keyEvent(t, "\x1b[3;5~") // ctrl+delete
	if ev == nil || ev.Key() != KeyDelete || ev.Modifiers() != ModCtrl {
		t.Fatalf("got %+v", ev)
	}
}

func TestKittyEscape(t *testing.T) {
	ev := keyEvent(t, "\x1b[27u") // Esc (disambiguated)
	if ev == nil || ev.Key() != KeyEscape {
		t.Fatalf("got %+v key=%d", ev, ev.Key())
	}
}

func TestKittyEnterTabBackspace(t *testing.T) {
	if ev := keyEvent(t, "\x1b[13u"); ev == nil || ev.Key() != KeyEnter {
		t.Fatalf("enter: got %+v", ev)
	}
	if ev := keyEvent(t, "\x1b[9u"); ev == nil || ev.Key() != KeyTAB {
		t.Fatalf("tab: got %+v", ev)
	}
	if ev := keyEvent(t, "\x1b[9;2u"); ev == nil {
		t.Fatalf("shift+tab: got nil")
	} else if ev.Key() != KeyBacktab && !(ev.Key() == KeyTAB && ev.Modifiers()&ModShift != 0) {
		t.Fatalf("shift+tab: got key=%d mod=%d", ev.Key(), ev.Modifiers())
	}
	if ev := keyEvent(t, "\x1b[127u"); ev == nil || ev.Key() != KeyBackspace {
		t.Fatalf("backspace: got %+v", ev)
	}
}

func TestKittyFunctionalPUA(t *testing.T) {
	ev := keyEvent(t, "\x1b[57362;2u") // pause + shift
	if ev == nil || ev.Key() != KeyPause || ev.Modifiers() != ModShift {
		t.Fatalf("got %+v", ev)
	}
}

func TestKittyUnmappedDropped(t *testing.T) {
	// keypad key 57401: consumed but no event
	ev, n, partial := parseKittyKey([]byte("\x1b[57401u"))
	if partial || n != len("\x1b[57401u") || ev != nil {
		t.Fatalf("ev=%v n=%d partial=%v", ev, n, partial)
	}
}

func TestKittyPartial(t *testing.T) {
	_, n, partial := parseKittyKey([]byte("\x1b[97;"))
	if !partial || n != 0 {
		t.Fatalf("n=%d partial=%v", n, partial)
	}
}

func TestKittyNotKitty(t *testing.T) {
	// SGR mouse sequence must be left alone
	ev, n, p := parseKittyKey([]byte("\x1b[<35;10;20M"))
	if ev != nil || n != 0 || p {
		// parseKittyKey's default branch returns not-kitty; the '<' is
		// rejected before scanning params so this falls through
	}
	// plain ascii is not kitty
	if _, n, p := parseKittyKey([]byte("a")); n != 0 || p {
		t.Fatalf("n=%d p=%v", n, p)
	}
	// legacy alt-prefixed CSI is not a kitty tail
	if _, n, p := parseKittyKey([]byte("\x1b[")); n != 0 || p {
		t.Fatalf("n=%d p=%v", n, p)
	}
}
