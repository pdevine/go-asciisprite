package termbox

// Kitty keyboard protocol: query-reply consumption and key event
// parsing.  Port of the proven tcell patch (kitty_parse.go) with the
// (partial, complete) stage convention and expire-aware partial claims.

import (
	"bytes"
	"strconv"
	"unicode"
)

// parseKittyReply consumes CSI ? flags u responses.  (partial, complete)
func parseKittyReply(buf *bytes.Buffer, expire bool) (bool, bool) {
	b := buf.Bytes()
	if len(b) < 3 || b[0] != '\x1b' || b[1] != '[' || b[2] != '?' {
		return false, false
	}
	i := 3
	for ; i < len(b); i++ {
		if b[i] < '0' || b[i] > '9' {
			break
		}
	}
	if i == len(b) {
		return !expire, false // partial while more input may come
	}
	if b[i] == 'u' {
		var flags int
		if i > 3 {
			flags, _ = strconv.Atoi(string(b[3:i]))
		}
		kittyFlags(flags)
		for j := 0; j <= i; j++ {
			buf.ReadByte()
		}
		return false, true
	}
	return false, false // some other ?-reply (DA1); leave to legacy
}

// kittyParamTail reports whether b could still become a kitty event:
// ESC [ digits[;:[digits]...] with no final byte yet.  A lone "\x1b[" is
// a legacy alt-key prefix, not a kitty tail — claiming it wedges the
// parser (regression: TestKittyChunked).
func kittyParamTail(b []byte) bool {
	if len(b) < 3 || b[0] != '\x1b' || b[1] != '[' {
		return false
	}
	for _, c := range b[2:] {
		if (c >= '0' && c <= '9') || c == ';' || c == ':' {
			continue
		}
		return false
	}
	return true
}

// parseKittyKey decodes one CSI u / CSI ~ / CSI letter kitty event from
// the buffer front.  (partial, complete)
func parseKittyKey(buf *bytes.Buffer, expire bool) (bool, bool) {
	b := buf.Bytes()
	if len(b) < 3 || b[0] != '\x1b' || b[1] != '[' || b[2] == '<' {
		return false, false // non-kitty prefix, or SGR mouse
	}
	i := 2
	for ; i < len(b); i++ {
		c := b[i]
		if (c >= '0' && c <= '9') || c == ';' || c == ':' {
			continue
		}
		break
	}
	if i == len(b) {
		if expire {
			return false, false
		}
		return kittyParamTail(b), false
	}
	final := b[i]
	params := string(b[2:i])
	n := i + 1

	var ev Event
	known := true
	switch final {
	case 'u':
		ev, known = kittyDecodeU(params)
	case '~':
		ev, known = kittyDecodeTilde(params)
	case 'A', 'B', 'C', 'D', 'E', 'F', 'H':
		ev, known = kittyDecodeLetter(params, final)
	default:
		return false, false
	}
	if !known {
		return false, false
	}
	for j := 0; j < n; j++ {
		buf.ReadByte()
	}
	if ev.Type != EventNone {
		post(ev)
	}
	return false, true
}

// kittyFields splits "a;b:c" style params.
func kittyFields(s string) [][]string {
	var out [][]string
	for _, f := range bytes.Split([]byte(s), []byte{';'}) {
		var subs []string
		for _, sub := range bytes.Split(f, []byte{':'}) {
			subs = append(subs, string(sub))
		}
		out = append(out, subs)
	}
	return out
}

func fieldInt(fs [][]string, f, sub int) int {
	if f >= len(fs) || sub >= len(fs[f]) {
		return 0
	}
	n, _ := strconv.Atoi(fs[f][sub])
	return n
}

// kittyMod decodes "mods[:eventType]" (field value is 1+bitfield).
func kittyMod(fs [][]string) (Modifier, EventType) {
	m := fieldInt(fs, 1, 0)
	typ := fieldInt(fs, 1, 1)
	var mod Modifier
	bits := m - 1
	if bits&1 != 0 {
		mod |= ModShift
	}
	if bits&2 != 0 {
		mod |= ModAlt
	}
	if bits&4 != 0 {
		mod |= ModCtrl
	}
	if bits&32 != 0 {
		mod |= ModMeta
	}
	var et EventType
	switch typ {
	case 2:
		et = EventKeyRepeat
	case 3:
		et = EventKeyRelease
	default:
		et = EventKeyPress
	}
	return mod, et
}

// kittyPUA maps private-use-area functional codepoints to keys.
var kittyPUA = map[int]Key{
	27:    KeyEsc,
	13:    KeyEnter,
	9:     KeyTab,
	127:   KeyBackspace2,
	57361: keyPrint,
	57362: keyPause,
	57363: keyCancel, // menu
	57376: keyF13,    // F13..F24 sequential
	57377: keyF13 + 1,
	57378: keyF13 + 2,
	57379: keyF13 + 3,
	57380: keyF13 + 4,
	57381: keyF13 + 5,
	57382: keyF13 + 6,
	57383: keyF13 + 7,
	57384: keyF13 + 8,
	57385: keyF13 + 9,
	57386: keyF13 + 10,
	57387: keyF13 + 11,
}

var kittyTilde = map[int]Key{
	2: KeyInsert, 3: KeyDelete, 5: KeyPgup, 6: KeyPgdn,
	7: KeyHome, 8: KeyEnd,
	11: KeyF1, 12: KeyF2, 13: KeyF3, 14: KeyF4,
	15: KeyF5, 17: KeyF6, 18: KeyF7, 19: KeyF8,
	20: KeyF9, 21: KeyF10, 23: KeyF11, 24: KeyF12,
}

var kittyLetter = map[byte]Key{
	'A': KeyArrowUp, 'B': KeyArrowDown, 'C': KeyArrowRight,
	'D': KeyArrowLeft, 'H': KeyHome, 'F': KeyEnd,
	'E': keyCenter,
}

func kittyDecodeU(params string) (Event, bool) {
	fs := kittyFields(params)
	cp := fieldInt(fs, 0, 0)
	shifted := fieldInt(fs, 0, 1)
	mod, et := kittyMod(fs)
	ev := Event{Type: et}
	if cp == 0 {
		return Event{}, true // pure text event: consumed, dropped
	}
	if k, ok := kittyPUA[cp]; ok {
		ev.Key = k
		return ev, true
	}
	if cp >= 57344 && cp <= 63743 {
		return Event{}, true // unmapped PUA (keypad/media/modifier): drop
	}
	r := rune(cp)
	if mod&ModShift != 0 && mod&ModCtrl == 0 {
		if shifted != 0 {
			r = rune(shifted)
		} else {
			r = []rune(string(unicode.ToUpper(r)))[0]
		}
	}
	ev.Key = keyForCtrlRune(r, &mod, &et)
	if ev.Key != 0 {
		return ev, true
	}
	ev.Ch = r
	return ev, true
}

// keyForCtrlRune returns a Key for the special control runes
// (enter/tab/backspace) and ctrl+letter combos, or 0 for plain runes.
// May adjust mod (drops ModCtrl when folding into KeyCtrlX) and et.
func keyForCtrlRune(r rune, mod *Modifier, et *EventType) Key {
	switch r {
	case '\r':
		return KeyEnter
	case '\t':
		if *mod&ModShift != 0 {
			*mod &^= ModCtrl | ModShift
			return KeyTab | keyShiftTab
		}
		return KeyTab
	case 0x7f, '\b':
		return KeyBackspace2
	}
	if *mod&ModCtrl != 0 && r >= 'a' && r <= 'z' && *mod&ModAlt == 0 {
		*mod &^= ModCtrl
		return Key(r - 'a' + 1)
	}
	return 0
}

// keyShiftTab marks shift-tab distinct from plain tab.
const keyShiftTab Key = 0x4000

func kittyDecodeTilde(params string) (Event, bool) {
	fs := kittyFields(params)
	num := fieldInt(fs, 0, 0)
	k, ok := kittyTilde[num]
	if !ok {
		return Event{}, true // consumed, unknown
	}
	mod, et := kittyMod(fs)
	return Event{Type: et, Key: k, Mod: mod}, true
}

func kittyDecodeLetter(params string, final byte) (Event, bool) {
	k, ok := kittyLetter[final]
	if !ok {
		return Event{}, false
	}
	fs := kittyFields(params)
	mod, et := kittyMod(fs)
	return Event{Type: et, Key: k, Mod: mod}, true
}
