package termbox

// Legacy input parser: a state machine over a pending byte buffer.
//
// Stages (matching tcell's collectEventsFromInput structure, which is
// proven against decades of terminal quirks):
//   kittyReply — CSI ? flags u enhancement responses (consumed always)
//   kittyKey   — CSI u / CSI ~ / CSI 1;mods X forms when enhanced
//   rune       — printable ASCII + UTF-8 runes
//   function   — legacy CSI/SS3/keycode table
//   mouse      — SGR and X10 mouse
//
// Each stage reports (partial, complete).  A stage that claims partial
// must drop the claim when expire fires, so the legacy byte-drain can
// deliver whatever was pending (the wedge regression of the tcell
// patch).

import (
	"bytes"
	"fmt"
	"os"
	"unicode/utf8"
)

// debugLog, when set (via TERmbox_DEBUG env var at Init), receives
// parser diagnostics for investigating live input issues.
var debugLog *os.File

// pendingInput holds undelivered bytes, delivered under keytimer expiry.
var pending bytes.Buffer
var escaped bool // saw a lone ESC; next key gets ModAlt

// parserEntry is a legacy escape-to-key table entry.
type parserEntry struct {
	key Key
	mod Modifier
}

// legacyKeys maps raw byte sequences to keys.  The table covers the
// xterm/kitty dialect (CSI and SS3 forms); see keys_legacy_test for the
// exact codes exercised.
var legacyKeys = map[string]parserEntry{
	"\x1b[A":    {KeyArrowUp, 0},
	"\x1b[B":    {KeyArrowDown, 0},
	"\x1b[C":    {KeyArrowRight, 0},
	"\x1b[D":    {KeyArrowLeft, 0},
	"\x1b[H":    {KeyHome, 0},
	"\x1b[F":    {KeyEnd, 0},
	"\x1bOA":    {KeyArrowUp, 0},
	"\x1bOB":    {KeyArrowDown, 0},
	"\x1bOC":    {KeyArrowRight, 0},
	"\x1bOD":    {KeyArrowLeft, 0},
	"\x1bOH":    {KeyHome, 0},
	"\x1bOF":    {KeyEnd, 0},
	"\x1bOP":    {KeyF1, 0},
	"\x1bOQ":    {KeyF2, 0},
	"\x1bOR":    {KeyF3, 0},
	"\x1bOS":    {KeyF4, 0},
	"\x1b[Z":    {KeyTab, ModShift}, // shift-tab
	"\x1b[2~":   {KeyInsert, 0},
	"\x1b[3~":   {KeyDelete, 0},
	"\x1b[5~":   {KeyPgup, 0},
	"\x1b[6~":   {KeyPgdn, 0},
	"\x1b[7~":   {KeyHome, 0},
	"\x1b[8~":   {KeyEnd, 0},
	"\x1b[11~":  {KeyF1, 0},
	"\x1b[12~":  {KeyF2, 0},
	"\x1b[13~":  {KeyF3, 0},
	"\x1b[14~":  {KeyF4, 0},
	"\x1b[15~":  {KeyF5, 0},
	"\x1b[17~":  {KeyF6, 0},
	"\x1b[18~":  {KeyF7, 0},
	"\x1b[19~":  {KeyF8, 0},
	"\x1b[20~":  {KeyF9, 0},
	"\x1b[21~":  {KeyF10, 0},
	"\x1b[23~":  {KeyF11, 0},
	"\x1b[24~":  {KeyF12, 0},
	"\x1b[1;2A": {KeyArrowUp, ModShift},
	"\x1b[1;2B": {KeyArrowDown, ModShift},
	"\x1b[1;2C": {KeyArrowRight, ModShift},
	"\x1b[1;2D": {KeyArrowLeft, ModShift},
	"\x1b[1;3A": {KeyArrowUp, ModAlt},
	"\x1b[1;3B": {KeyArrowDown, ModAlt},
	"\x1b[1;3C": {KeyArrowRight, ModAlt},
	"\x1b[1;3D": {KeyArrowLeft, ModAlt},
	"\x1b[1;5A": {KeyArrowUp, ModCtrl},
	"\x1b[1;5B": {KeyArrowDown, ModCtrl},
	"\x1b[1;5C": {KeyArrowRight, ModCtrl},
	"\x1b[1;5D": {KeyArrowLeft, ModCtrl},
	"\x1b[1;6A": {KeyArrowUp, ModShift | ModCtrl},
	"\x1b[1;6B": {KeyArrowDown, ModShift | ModCtrl},
	"\x1b[1;6C": {KeyArrowRight, ModShift | ModCtrl},
	"\x1b[1;6D": {KeyArrowLeft, ModShift | ModCtrl},
}

// ctrl mapping: bytes 1-26 are Ctrl+a..Ctrl+z.
func ctrlKey(b byte) Key { return Key(b) }

// drainLocked runs the stage pipeline over the pending buffer.
// expire=true flushes ambiguous prefixes byte-by-byte.
func drainLocked(expire bool) {
	if d := debugLog; d != nil {
		fmt.Fprintf(d, "drain expire=%v pending=%q\n", expire, pending.String())
		d.Sync()
	}
	for {
		b := pending.Bytes()
		if len(b) == 0 {
			return
		}
		partials := 0

		if part, done := parseKittyReply(&pending, expire); done {
			continue
		} else if part {
			partials++
		}
		if enhanced {
			if part, done := parseKittyKey(&pending, expire); done {
				continue
			} else if part {
				partials++
			}
		}
		if part, done := parseFunction(&pending); done {
			continue
		} else if part {
			partials++
		}
		if part, done := parseMouse(&pending); done {
			continue
		} else if part {
			partials++
		}
		if part, done := parseRuneLegacy(&pending); done {
			continue
		} else if part {
			partials++
		}

		if partials == 0 || expire {
			// no full match and nothing more coming: deliver byte-by-byte
			c, _ := pending.ReadByte()
			if c == 0x1b {
				if pending.Len() == 0 {
					post(Event{Type: EventKey, Key: KeyEsc})
					escaped = false
				} else {
					escaped = true
				}
				continue
			}
			mod := Modifier(0)
			if escaped {
				escaped = false
				mod = ModAlt
			}
			ev := Event{Type: EventKey, Mod: mod}
			switch {
			case c == '\r':
				ev.Key = KeyEnter
			case c == '\t':
				ev.Key = KeyTab
			case c == 0x7f:
				ev.Key = KeyBackspace2
			case c == 8:
				ev.Key = KeyBackspace
			case c == 0:
				ev.Key = KeySpace
				ev.Mod |= ModCtrl
			case c >= 1 && c <= 26:
				ev.Key = ctrlKey(c)
			case c == ' ':
				ev.Key = KeySpace
				ev.Ch = ' '
			default:
				ev.Ch = rune(c)
			}
			post(ev)
			continue
		}
		return // partial: await more bytes
	}
}

func parseBytes(b []byte) {
	inputMu.Lock()
	defer inputMu.Unlock()
	pending.Write(b)
	drainLocked(false)
}

func expireInput() {
	inputMu.Lock()
	defer inputMu.Unlock()
	if pending.Len() > 0 {
		drainLocked(true)
	}
}

func parseRuneLegacy(buf *bytes.Buffer) (bool, bool) {
	b := buf.Bytes()
	c := b[0]
	switch {
	case c >= ' ' && c < 0x7f:
		mod := Modifier(0)
		if escaped {
			mod = ModAlt
			escaped = false
		}
		buf.ReadByte()
		if c == ' ' {
			post(Event{Type: EventKey, Key: KeySpace, Ch: ' ', Mod: mod})
		} else {
			post(Event{Type: EventKey, Ch: rune(c), Mod: mod})
		}
		return false, true
	case c < 0x80:
		return false, false // control chars fall to the drain
	default:
		// UTF-8 multibyte
		r, n := utf8.DecodeRune(b)
		if r == utf8.RuneError {
			if !utf8.FullRune(b) {
				return true, false // partial
			}
			buf.ReadByte() // invalid byte: drop
			return false, true
		}
		mod := Modifier(0)
		if escaped {
			mod = ModAlt
			escaped = false
		}
		for i := 0; i < n; i++ {
			buf.ReadByte()
		}
		post(Event{Type: EventKey, Ch: r, Mod: mod})
		return false, true
	}
}

func parseFunction(buf *bytes.Buffer) (bool, bool) {
	b := buf.Bytes()
	partial := false
	for esc, k := range legacyKeys {
		if bytes.HasPrefix(b, []byte(esc)) {
			for range len(esc) {
				buf.ReadByte()
			}
			mod := k.mod
			if escaped {
				mod |= ModAlt
				escaped = false
			}
			post(Event{Type: EventKey, Key: k.key, Mod: mod})
			return false, true
		}
		if bytes.HasPrefix([]byte(esc), b) {
			partial = true
		}
	}
	return partial, false
}
