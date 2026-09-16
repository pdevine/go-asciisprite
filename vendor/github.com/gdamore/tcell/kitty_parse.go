// Kitty keyboard protocol support — in-tree patch for go-asciisprite.
// Byte-level parser for the CSI u / CSI ~ / CSI <letter> event forms per
// https://sw.kovidgoyal.net/kitty/keyboard-protocol/

package tcell

import (
	"bytes"
	"strconv"
	"time"
	"unicode"
)

// NewEventKeyTyp is NewEventKey with an explicit event type,
// for kitty keyboard protocol events.
func NewEventKeyTyp(k Key, ch rune, mod ModMask, typ KeyEventType) *EventKey {
	ev := NewEventKey(k, ch, mod)
	ev.typ = typ
	return ev
}

// kittyParamTail reports whether b could still become a kitty event or
// kitty reply: ESC [ <?> digits[;:[digits]...] with no final byte yet.
// A lone "\x1b[" is a legacy alt-key prefix, not a kitty tail.  The
// kitty parse stage must not hold the buffer for that shape or the
// legacy table never gets its chance (escalating into a hung parser).
func kittyParamTail(b []byte) bool {
	if len(b) < 3 || b[0] != '\x1b' || b[1] != '[' {
		return false
	}
	i := 2
	if b[2] == '?' {
		i = 3
	}
	for ; i < len(b); i++ {
		c := b[i]
		if (c >= '0' && c <= '9') || c == ';' || c == ':' {
			continue
		}
		return false
	}
	return true
}

// parseKittyKey attempts to decode one kitty keyboard-protocol key event
// from the front of b.
//
// Returns (ev, consumed, partial):
//   - consumed > 0: that many input bytes were eaten. ev may be nil for
//     recognized-but-unmapped keys (keypad/media/modifier/lock keys) or
//     pure text events (key codepoint 0).
//   - partial == true: b is a prefix of a kitty sequence; wait for more
//     input (consumed is 0, ev nil).
//   - otherwise: not a kitty sequence at all.
func parseKittyKey(b []byte) (ev *EventKey, consumed int, partial bool) {
	if len(b) < 2 || b[0] != '\x1b' || b[1] != '[' {
		return nil, 0, false
	}
	// SGR mouse shape is not a kitty event (and not a kitty partial)
	if len(b) >= 3 && b[2] == '<' {
		return nil, 0, false
	}
	// scan for a final byte
	i := 2
	for ; i < len(b); i++ {
		c := b[i]
		if (c >= '0' && c <= '9') || c == ';' || c == ':' {
			continue
		}
		break
	}
	if i == len(b) {
		// possibly incomplete kitty sequence — but only hold the
		// buffer when it actually looks like one
		return nil, 0, kittyParamTail(b)
	}
	final := b[i]
	params := string(b[2:i])
	consumed = i + 1

	switch final {
	case 'u':
		return parseKittyU(params), consumed, false
	case '~':
		return parseKittyTilde(params), consumed, false
	case 'A', 'B', 'C', 'D', 'E', 'F', 'H':
		if _, ok := kittyLetterKeys[final]; !ok {
			return nil, consumed, false
		}
		// only "1[;mods]" shapes are kitty; other legacy shapes are
		// filtered by the caller (we are only invoked in enhanced mode)
		mod, typ := parseKittyModField(fields2(params, 1))
		key := kittyLetterKeys[final]
		return NewEventKeyTyp(key, 0, mod, typ), consumed, false
	default:
		// not a kitty event after all (e.g. SGR mouse, replies)
		return nil, 0, false
	}
}

// fields2 returns sub-field n of field 1 of a kitty param string,
// tolerating missing fields (empty string).
func fields2(params string, n int) string {
	fs := bytes.Split([]byte(params), []byte{';'})
	if len(fs) <= n {
		return ""
	}
	return string(fs[n])
}

// parseKittyModField decodes "mods[:eventType]" returning the tcell mask
// and event type.  An empty/absent field means no modifiers, press.
func parseKittyModField(s string) (ModMask, KeyEventType) {
	if s == "" {
		return ModNone, KeyEventPress
	}
	subs := bytes.Split([]byte(s), []byte{':'})
	m, _ := strconv.Atoi(string(subs[0]))
	var typ KeyEventType
	if len(subs) > 1 {
		if t, _ := strconv.Atoi(string(subs[1])); t == 2 {
			typ = KeyEventRepeat
		} else if t == 3 {
			typ = KeyEventRelease
		} else {
			typ = KeyEventPress
		}
	} else {
		typ = KeyEventPress
	}
	return kittyModMask(m - 1), typ
}

// parseKittyU decodes the CSI ... u form:
// codepoint[:shifted[:base]] ; mods[:eventType] ; text...
func parseKittyU(params string) *EventKey {
	fields := bytes.Split([]byte(params), []byte{';'})

	// field 0: key codepoint and alternates
	subs := bytes.Split(fields[0], []byte{':'})
	cp, _ := strconv.Atoi(string(subs[0]))
	shifted := 0
	if len(subs) > 1 {
		shifted, _ = strconv.Atoi(string(subs[1]))
	}

	mod, typ := ModNone, KeyEventPress
	if len(fields) > 1 {
		mod, typ = parseKittyModField(string(fields[1]))
	}

	if cp == 0 {
		// pure text event (no key info); we don't request associated
		// text, so this shouldn't occur — drop it.
		return nil
	}

	if key, ok := kittyFunctional[cp]; ok {
		return NewEventKeyTyp(key, 0, mod, typ)
	}
	if cp >= 57344 {
		// unmapped PUA functional key (keypad/media/modifier): consumed
		// intentionally but no event to report.
		return nil
	}

	r := rune(cp)
	// if shift produced a different codepoint, report the shifted rune
	// (but not for ctrl-based combos: keep the unshifted letter so
	// ctrl+shift+x matches the same KeyCtrl* value as ctrl+x)
	if mod&ModShift != 0 && shifted != 0 && mod&ModCtrl == 0 {
		r = rune(shifted)
	} else if mod&ModShift != 0 && mod&ModCtrl == 0 {
		r = []rune(string(unicode.ToUpper(r)))[0]
	}

	// enter/tab/backspace arrive as raw control codes in some kitty
	// shapes; normalize them to their tcell keys
	switch r {
	case '\r':
		return NewEventKeyTyp(KeyEnter, 0, mod, typ)
	case '\t':
		if mod&ModShift != 0 {
			return NewEventKeyTyp(KeyBacktab, 0, mod&^ModShift, typ)
		}
		return NewEventKeyTyp(KeyTAB, 0, mod, typ)
	case 0x7f, '\b':
		return NewEventKeyTyp(KeyBackspace, 0, mod, typ)
	}

	// map ctrl+letter to tcell's legacy KeyCtrl* where the letter is
	// lowercase ASCII, so existing Key-based checks keep matching
	if mod&ModCtrl != 0 && r >= 'a' && r <= 'z' && mod&^(ModCtrl|ModShift|ModAlt) == 0 && mod&ModAlt == 0 {
		return NewEventKeyTyp(Key(r-'a'+1), 0, mod&^ModCtrl, typ)
	}

	return NewEventKeyTyp(KeyRune, r, mod, typ)
}

// parseKittyTilde decodes the CSI number;mods ~ form.
func parseKittyTilde(params string) *EventKey {
	fields := bytes.Split([]byte(params), []byte{';'})
	num, _ := strconv.Atoi(string(fields[0]))
	key, ok := kittyTildeKeys[num]
	if !ok {
		return nil
	}
	mod, typ := ModNone, KeyEventPress
	if len(fields) > 1 {
		mod, typ = parseKittyModField(string(fields[1]))
	}
	return NewEventKeyTyp(key, 0, mod, typ)
}

// parseKittyReply consumes a kitty enhancement-flags query response,
// CSI ? flags u, and records support.  Returns (partial, handled) —
// matching the (part, comp) convention of the other parse stages.
// expire is set when the input timer fired: in that case partial
// claims are abandoned so the legacy fallback can drain the buffer.
func (t *tScreen) parseKittyReply(buf *bytes.Buffer, expire bool) (bool, bool) {
	b := buf.Bytes()
	if len(b) < 3 || b[0] != '\x1b' || b[1] != '[' || b[2] != '?' {
		return false, false
	}
	i := 3
	for ; i < len(b); i++ {
		c := b[i]
		if c >= '0' && c <= '9' {
			continue
		}
		break
	}
	if i == len(b) {
		return !expire, false // partial only while more input may come
	}
	if b[i] == 'u' {
		var flags int
		if i > 3 {
			flags, _ = strconv.Atoi(string(b[3:i]))
		}
		t.kittySupported = flags
		select {
		case t.kittyReplyCh <- flags:
		default:
		}
		for j := 0; j <= i; j++ {
			buf.ReadByte()
		}
		return false, true
	}
	// some other "?"-reply (e.g. DA1 response ends in 'c'); leave it
	return false, false
}

// parseKittyStage is parseKittyKey adapted to the buffer/pipeline shape
// used by collectEventsFromInput.  See parseKittyReply for expire.
func (t *tScreen) parseKittyStage(buf *bytes.Buffer, evs *[]Event, expire bool) (bool, bool) {
	b := buf.Bytes()
	if len(b) < 2 || b[0] != '\x1b' || b[1] != '[' {
		return false, false
	}
	// SGR mouse takes priority for its own shape
	if len(b) >= 3 && b[2] == '<' {
		return false, false
	}
	ev, n, partial := parseKittyKey(b)
	if partial && !expire {
		return true, false
	}
	if n == 0 {
		return false, false
	}
	for i := 0; i < n; i++ {
		buf.ReadByte()
	}
	if ev != nil {
		ev.t = time.Now()
		*evs = append(*evs, ev)
	}
	return false, true
}
