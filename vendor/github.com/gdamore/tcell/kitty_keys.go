// Kitty keyboard protocol support — in-tree patch for go-asciisprite.
// Key code and modifier mapping tables per
// https://sw.kovidgoyal.net/kitty/keyboard-protocol/

package tcell

// kittyLetterKeys maps the final byte of the CSI 1;mods <letter> form.
var kittyLetterKeys = map[byte]Key{
	'A': KeyUp,
	'B': KeyDown,
	'C': KeyRight,
	'D': KeyLeft,
	'E': KeyCenter, // keypad begin legacy form
	'H': KeyHome,
	'F': KeyEnd,
}

// kittyTildeKeys maps the numeric field of the CSI number;mods ~ form.
var kittyTildeKeys = map[int]Key{
	2:  KeyInsert,
	3:  KeyDelete,
	5:  KeyPgUp,
	6:  KeyPgDn,
	7:  KeyHome,
	8:  KeyEnd,
	11: KeyF1,
	12: KeyF2,
	13: KeyF3,
	14: KeyF4,
	15: KeyF5,
	17: KeyF6,
	18: KeyF7,
	19: KeyF8,
	20: KeyF9,
	21: KeyF10,
	23: KeyF11,
	24: KeyF12,
	29: KeyCancel, // menu key; no dedicated KeyMenu in v1
}

// kittyFunctional maps PUA codepoints (57344+) to tcell keys.
// Keys with no tcell v1 equivalent (keypad, media, physical modifier
// keys) are intentionally unmapped; they are consumed but dropped.
var kittyFunctional = map[int]Key{
	27:    KeyEscape,
	13:    KeyEnter,
	9:     KeyTAB,
	127:   KeyBackspace,
	57358: KeyRune, // caps lock: placeholder rune key; ch=0 marks it
	57359: KeyRune, // scroll lock
	57360: KeyRune, // num lock
	57361: KeyPrint,
	57362: KeyPause,
	57363: KeyCancel, // menu key; no dedicated KeyMenu in v1
	57376: KeyF13,
	57377: KeyF14,
	57378: KeyF15,
	57379: KeyF16,
	57380: KeyF17,
	57381: KeyF18,
	57382: KeyF19,
	57383: KeyF20,
	57384: KeyF21,
	57385: KeyF22,
	57386: KeyF23,
	57387: KeyF24,
}

// kittyModMask converts a kitty modifier field (already minus the
// protocol's +1 offset, i.e. the raw bitfield) into a tcell ModMask.
// super/hyper/caps/num have no tcell v1 equivalent and are dropped;
// kitty meta maps onto ModMeta.
func kittyModMask(bits int) ModMask {
	var m ModMask
	if bits&1 != 0 { // shift
		m |= ModShift
	}
	if bits&2 != 0 { // alt
		m |= ModAlt
	}
	if bits&4 != 0 { // ctrl
		m |= ModCtrl
	}
	if bits&32 != 0 { // meta
		m |= ModMeta
	}
	return m
}
