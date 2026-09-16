// Termbox compatibility constants — numeric compatibility with the
// values produced by the tcell-backed implementation is required: these
// were historically `Key(tcell.KeyX)` casts and consumers may embed the
// numbers.  Do not renumber.  Verified by compat_keys_test.go.

package termbox

// Modifier masks (same bit layout as tcell.ModMask).
type Modifier int16

const (
	ModShift Modifier = 1 << iota
	ModCtrl
	ModAlt
	ModMeta
)

// Key identifies a key press; printable runes arrive via Event.Ch with
// Key == 0, and KeySpace is a named alias for ' '.
type Key int16

// Numeric values preserve the tcell v1 assignment: KeyRune (unused as a
// public value here) is 256, followed by the named keys in tcell's
// declaration order.
const (
	keyBase Key = 256 // tcell KeyRune; not exported

	KeyArrowUp    Key = keyBase + 1
	KeyArrowDown  Key = keyBase + 2
	KeyArrowRight Key = keyBase + 3
	KeyArrowLeft  Key = keyBase + 4
	keyUpLeft     Key = keyBase + 5
	keyUpRight    Key = keyBase + 6
	keyDownLeft   Key = keyBase + 7
	keyDownRight  Key = keyBase + 8
	keyCenter     Key = keyBase + 9
	KeyPgup       Key = keyBase + 10
	KeyPgdn       Key = keyBase + 11
	KeyHome       Key = keyBase + 12
	KeyEnd        Key = keyBase + 13
	KeyInsert     Key = keyBase + 14
	KeyDelete     Key = keyBase + 15
	keyHelp       Key = keyBase + 16
	keyExit       Key = keyBase + 17
	keyClear      Key = keyBase + 18
	keyCancel     Key = keyBase + 19
	keyPrint      Key = keyBase + 20
	keyPause      Key = keyBase + 21
	keyBacktab    Key = keyBase + 22
	KeyF1         Key = keyBase + 23
	KeyF2         Key = keyBase + 24
	KeyF3         Key = keyBase + 25
	KeyF4         Key = keyBase + 26
	KeyF5         Key = keyBase + 27
	KeyF6         Key = keyBase + 28
	KeyF7         Key = keyBase + 29
	KeyF8         Key = keyBase + 30
	KeyF9         Key = keyBase + 31
	KeyF10        Key = keyBase + 32
	KeyF11        Key = keyBase + 33
	KeyF12        Key = keyBase + 34
	keyF13        Key = keyBase + 35
	keyF24        Key = keyBase + 46
	keyF61        Key = keyBase + 83
	keyF62        Key = keyBase + 84
	keyF63        Key = keyBase + 85
	keyF64        Key = keyBase + 86

	// Legacy control-key aliases (same values as tcell's).
	KeyBackspace  Key = 8
	KeyBackspace2 Key = 127
	KeyTab        Key = 9
	KeyEnter      Key = 13
	KeyEsc        Key = 27

	// Mouse pseudo-buttons: assigned from tcell's arbitrary F-key range
	// by the original compat layer.
	MouseLeft    Key = keyF63
	MouseRight   Key = keyF62
	MouseMiddle  Key = keyF61
	MouseRelease Key = keyF64

	KeySpace Key = Key(' ')
)

// Ctrl-letter keys — same values as tcell's KeyCtrl* (ASCII control codes).
const (
	KeyCtrlA Key = 1 + iota
	KeyCtrlB
	KeyCtrlC
	KeyCtrlD
	KeyCtrlE
	KeyCtrlF
	KeyCtrlG
	KeyCtrlH
	KeyCtrlI
	KeyCtrlJ
	KeyCtrlK
	KeyCtrlL
	KeyCtrlM
	KeyCtrlN
	KeyCtrlO
	KeyCtrlP
	KeyCtrlQ
	KeyCtrlR
	KeyCtrlS
	KeyCtrlT
	KeyCtrlU
	KeyCtrlV
	KeyCtrlW
	KeyCtrlX
	KeyCtrlY
	KeyCtrlZ
)
