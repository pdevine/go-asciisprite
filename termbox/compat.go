// Copyright 2020 The TCell Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use file except in compliance with the License.
// You may obtain a copy of the license at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package termbox is a compatibility layer to allow tcell to emulate
// the github.com/nsf/termbox package.
package termbox

import (
	"errors"

	"github.com/gdamore/tcell"
)

var screen tcell.Screen
var outMode OutputMode

// Init initializes the screen for use.
func Init() error {
	outMode = OutputNormal
	//outMode = Output256
	if s, e := tcell.NewScreen(); e != nil {
		return e
	} else if e = s.Init(); e != nil {
		return e
	} else {
		s.EnableMouse()
		screen = s
		return nil
	}
}

// Close cleans up the terminal, restoring terminal modes, etc.
func Close() {
	screen.Fini()
}

// Flush updates the screen.
func Flush() error {
	screen.Show()
	return nil
}

// SetCursor displays the terminal cursor at the given location.
func SetCursor(x, y int) {
	screen.ShowCursor(x, y)
}

// HideCursor hides the terminal cursor.
func HideCursor() {
	SetCursor(-1, -1)
}

// Size returns the screen size as width, height in character cells.
func Size() (int, int) {
	return screen.Size()
}

// Attribute affects the presentation of characters, such as color, boldness,
// and so forth.
type Attribute uint16

// Colors first.  The order here is significant.
const (
	ColorDefault Attribute = iota
	ColorBlack
	ColorMaroon
	ColorGreen
	ColorOlive
	ColorNavy
	ColorPurple
	ColorTeal
	ColorSilver
	ColorGray
	ColorRed
	ColorLime
	ColorYellow
	ColorBlue
	ColorFuchsia
	ColorAqua
	ColorWhite
	Color16
	Color17
	Color18
	Color19
	Color20
	Color21
	Color22
	Color23
	Color24
	Color25
	Color26
	Color27
	Color28
	Color29
	Color30
	Color31
	Color32
	Color33
	Color34
	Color35
	Color36
	Color37
	Color38
	Color39
	Color40
	Color41
	Color42
	Color43
	Color44
	Color45
	Color46
	Color47
	Color48
	Color49
	Color50
	Color51
	Color52
	Color53
	Color54
	Color55
	Color56
	Color57
	Color58
	Color59
	Color60
	Color61
	Color62
	Color63
	Color64
	Color65
	Color66
	Color67
	Color68
	Color69
	Color70
	Color71
	Color72
	Color73
	Color74
	Color75
	Color76
	Color77
	Color78
	Color79
	Color80
	Color81
	Color82
	Color83
	Color84
	Color85
	Color86
	Color87
	Color88
	Color89
	Color90
	Color91
	Color92
	Color93
	Color94
	Color95
	Color96
	Color97
	Color98
	Color99
	Color100
	Color101
	Color102
	Color103
	Color104
	Color105
	Color106
	Color107
	Color108
	Color109
	Color110
	Color111
	Color112
	Color113
	Color114
	Color115
	Color116
	Color117
	Color118
	Color119
	Color120
	Color121
	Color122
	Color123
	Color124
	Color125
	Color126
	Color127
	Color128
	Color129
	Color130
	Color131
	Color132
	Color133
	Color134
	Color135
	Color136
	Color137
	Color138
	Color139
	Color140
	Color141
	Color142
	Color143
	Color144
	Color145
	Color146
	Color147
	Color148
	Color149
	Color150
	Color151
	Color152
	Color153
	Color154
	Color155
	Color156
	Color157
	Color158
	Color159
	Color160
	Color161
	Color162
	Color163
	Color164
	Color165
	Color166
	Color167
	Color168
	Color169
	Color170
	Color171
	Color172
	Color173
	Color174
	Color175
	Color176
	Color177
	Color178
	Color179
	Color180
	Color181
	Color182
	Color183
	Color184
	Color185
	Color186
	Color187
	Color188
	Color189
	Color190
	Color191
	Color192
	Color193
	Color194
	Color195
	Color196
	Color197
	Color198
	Color199
	Color200
	Color201
	Color202
	Color203
	Color204
	Color205
	Color206
	Color207
	Color208
	Color209
	Color210
	Color211
	Color212
	Color213
	Color214
	Color215
	Color216
	Color217
	Color218
	Color219
	Color220
	Color221
	Color222
	Color223
	Color224
	Color225
	Color226
	Color227
	Color228
	Color229
	Color230
	Color231
	Color232
	Color233
	Color234
	Color235
	Color236
	Color237
	Color238
	Color239
	Color240
	Color241
	Color242
	Color243
	Color244
	Color245
	Color246
	Color247
	Color248
	Color249
	Color250
	Color251
	Color252
	Color253
	Color254
	Color255
	ColorAliceBlue
	ColorAntiqueWhite
	ColorAquaMarine
	ColorAzure
	ColorBeige
	ColorBisque
	ColorBlanchedAlmond
	ColorBlueViolet
	ColorBrown
	ColorBurlyWood
	ColorCadetBlue
	ColorChartreuse
	ColorChocolate
	ColorCoral
	ColorCornflowerBlue
	ColorCornsilk
	ColorCrimson
	ColorDarkBlue
	ColorDarkCyan
	ColorDarkGoldenrod
	ColorDarkGray
	ColorDarkGreen
	ColorDarkKhaki
	ColorDarkMagenta
	ColorDarkOliveGreen
	ColorDarkOrange
	ColorDarkOrchid
	ColorDarkRed
	ColorDarkSalmon
	ColorDarkSeaGreen
	ColorDarkSlateBlue
	ColorDarkSlateGray
	ColorDarkTurquoise
	ColorDarkViolet
	ColorDeepPink
	ColorDeepSkyBlue
	ColorDimGray
	ColorDodgerBlue
	ColorFireBrick
	ColorFloralWhite
	ColorForestGreen
	ColorGainsboro
	ColorGhostWhite
	ColorGold
	ColorGoldenrod
	ColorGreenYellow
	ColorHoneydew
	ColorHotPink
	ColorIndianRed
	ColorIndigo
	ColorIvory
	ColorKhaki
	ColorLavender
	ColorLavenderBlush
	ColorLawnGreen
	ColorLemonChiffon
	ColorLightBlue
	ColorLightCoral
	ColorLightCyan
	ColorLightGoldenrodYellow
	ColorLightGray
	ColorLightGreen
	ColorLightPink
	ColorLightSalmon
	ColorLightSeaGreen
	ColorLightSkyBlue
	ColorLightSlateGray
	ColorLightSteelBlue
	ColorLightYellow
	ColorLimeGreen
	ColorLinen
	ColorMediumAquamarine
	ColorMediumBlue
	ColorMediumOrchid
	ColorMediumPurple
	ColorMediumSeaGreen
	ColorMediumSlateBlue
	ColorMediumSpringGreen
	ColorMediumTurquoise
	ColorMediumVioletRed
	ColorMidnightBlue
	ColorMintCream
	ColorMistyRose
	ColorMoccasin
	ColorNavajoWhite
	ColorOldLace
	ColorOliveDrab
	ColorOrange
	ColorOrangeRed
	ColorOrchid
	ColorPaleGoldenrod
	ColorPaleGreen
	ColorPaleTurquoise
	ColorPaleVioletRed
	ColorPapayaWhip
	ColorPeachPuff
	ColorPeru
	ColorPink
	ColorPlum
	ColorPowderBlue
	ColorRebeccaPurple
	ColorRosyBrown
	ColorRoyalBlue
	ColorSaddleBrown
	ColorSalmon
	ColorSandyBrown
	ColorSeaGreen
	ColorSeashell
	ColorSienna
	ColorSkyblue
	ColorSlateBlue
	ColorSlateGray
	ColorSnow
	ColorSpringGreen
	ColorSteelBlue
	ColorTan
	ColorThistle
	ColorTomato
	ColorTurquoise
	ColorViolet
	ColorWheat
	ColorWhiteSmoke
	ColorYellowGreen
)

// Other attributes.
const (
	AttrBold Attribute = 1 << (9 + iota)
	AttrUnderline
	AttrReverse
)

func fixColor(c tcell.Color) tcell.Color {
	if c == tcell.ColorDefault {
		return c
	}
	switch outMode {
	case OutputNormal:
		//c = tcell.PaletteColor(int(c) & 0xf)
	case Output256:
		//c = tcell.PaletteColor(int(c) & 0xff)
		c %= tcell.Color(256)
	case Output216:
		//c = tcell.PaletteColor(int(c)%216 + 16)
	case OutputGrayscale:
		//c %= tcell.PaletteColor(int(c)%24 + 232)
	default:
		c = tcell.ColorDefault
	}
	return c
}

func mkStyle(fg, bg Attribute) tcell.Style {
	st := tcell.StyleDefault

	//f := tcell.PaletteColor(int(fg)&0x1ff - 1)
	//b := tcell.PaletteColor(int(bg)&0x1ff - 1)
	f := tcell.Color(int(fg)&0x1ff) - 1
	b := tcell.Color(int(bg)&0x1ff) - 1


	f = fixColor(f)
	b = fixColor(b)
	st = st.Foreground(f).Background(b)
	if (fg|bg)&AttrBold != 0 {
		st = st.Bold(true)
	}
	if (fg|bg)&AttrUnderline != 0 {
		st = st.Underline(true)
	}
	if (fg|bg)&AttrReverse != 0 {
		st = st.Reverse(true)
	}
	return st
}

func mkColors(st tcell.Style) (Attribute, Attribute) {
	// XXX - figure out how to do attribs
	f, b, _ := st.Decompose()
	fg := Attribute(int(f)+1)
	bg := Attribute(int(b)+1)
	return fg, bg
}

// Clear clears the screen with the given attributes.
func Clear(fg, bg Attribute) {
	st := mkStyle(fg, bg)
	w, h := screen.Size()
	for row := 0; row < h; row++ {
		for col := 0; col < w; col++ {
			screen.SetContent(col, row, ' ', nil, st)
		}
	}
}

// InputMode is not used.
type InputMode int

// Unused input modes; here for compatibility.
const (
	InputEsc InputMode = 1 << iota
	InputAlt
	InputMouse
	InputCurrent InputMode = 0
)

// SetInputMode does not do much in this version.
func SetInputMode(mode InputMode) InputMode {
	if mode&InputMouse != 0 {
		screen.EnableMouse()
	} else {
		screen.DisableMouse()
	}
	return InputEsc
}

// OutputMode represents an output mode, which determines how colors
// are used.  See the termbox documentation for an explanation.
type OutputMode int

// OutputMode values.
const (
	OutputCurrent OutputMode = iota
	OutputNormal
	Output256
	Output216
	OutputGrayscale
)

// SetOutputMode is used to set the color palette used.
func SetOutputMode(mode OutputMode) OutputMode {
	if screen.Colors() < 256 {
		mode = OutputNormal
	}
	switch mode {
	case OutputCurrent:
		return outMode
	case OutputNormal, Output256, Output216, OutputGrayscale:
		outMode = mode
		return mode
	default:
		return outMode
	}
}

// Sync forces a resync of the screen.
func Sync() error {
	screen.Sync()
	return nil
}

// SetCell sets the character cell at a given location to the given
// content (rune) and attributes.
func SetCell(x, y int, ch rune, fg, bg Attribute) {
	st := mkStyle(fg, bg)
	screen.SetContent(x, y, ch, nil, st)
}

// GetCell
func GetCell(x, y int) (rune, Attribute, Attribute) {
	ch, _, st, _ := screen.GetContent(x, y)
	fg, bg := mkColors(st)
	return ch, fg, bg
}

// EventType represents the type of event.
type EventType uint8

// Modifier represents the possible modifier keys.
type Modifier tcell.ModMask

// Key is a key press.
type Key tcell.Key

// Event represents an event like a key press, mouse action, or window resize.
type Event struct {
	Type   EventType
	Mod    Modifier
	Key    Key
	Ch     rune
	Width  int
	Height int
	Err    error
	MouseX int
	MouseY int
	N      int
}

// Event types.
const (
	EventNone EventType = iota
	EventKey
	EventResize
	EventMouse
	EventInterrupt
	EventError
	EventRaw
)

// Keys codes.
const (
	KeyF1         = Key(tcell.KeyF1)
	KeyF2         = Key(tcell.KeyF2)
	KeyF3         = Key(tcell.KeyF3)
	KeyF4         = Key(tcell.KeyF4)
	KeyF5         = Key(tcell.KeyF5)
	KeyF6         = Key(tcell.KeyF6)
	KeyF7         = Key(tcell.KeyF7)
	KeyF8         = Key(tcell.KeyF8)
	KeyF9         = Key(tcell.KeyF9)
	KeyF10        = Key(tcell.KeyF10)
	KeyF11        = Key(tcell.KeyF11)
	KeyF12        = Key(tcell.KeyF12)
	KeyInsert     = Key(tcell.KeyInsert)
	KeyDelete     = Key(tcell.KeyDelete)
	KeyHome       = Key(tcell.KeyHome)
	KeyEnd        = Key(tcell.KeyEnd)
	KeyArrowUp    = Key(tcell.KeyUp)
	KeyArrowDown  = Key(tcell.KeyDown)
	KeyArrowRight = Key(tcell.KeyRight)
	KeyArrowLeft  = Key(tcell.KeyLeft)
	KeyCtrlA      = Key(tcell.KeyCtrlA)
	KeyCtrlB      = Key(tcell.KeyCtrlB)
	KeyCtrlC      = Key(tcell.KeyCtrlC)
	KeyCtrlD      = Key(tcell.KeyCtrlD)
	KeyCtrlE      = Key(tcell.KeyCtrlE)
	KeyCtrlF      = Key(tcell.KeyCtrlF)
	KeyCtrlG      = Key(tcell.KeyCtrlG)
	KeyCtrlH      = Key(tcell.KeyCtrlH)
	KeyCtrlI      = Key(tcell.KeyCtrlI)
	KeyCtrlJ      = Key(tcell.KeyCtrlJ)
	KeyCtrlK      = Key(tcell.KeyCtrlK)
	KeyCtrlL      = Key(tcell.KeyCtrlL)
	KeyCtrlM      = Key(tcell.KeyCtrlM)
	KeyCtrlN      = Key(tcell.KeyCtrlN)
	KeyCtrlO      = Key(tcell.KeyCtrlO)
	KeyCtrlP      = Key(tcell.KeyCtrlP)
	KeyCtrlQ      = Key(tcell.KeyCtrlQ)
	KeyCtrlR      = Key(tcell.KeyCtrlR)
	KeyCtrlS      = Key(tcell.KeyCtrlS)
	KeyCtrlT      = Key(tcell.KeyCtrlT)
	KeyCtrlU      = Key(tcell.KeyCtrlU)
	KeyCtrlV      = Key(tcell.KeyCtrlV)
	KeyCtrlW      = Key(tcell.KeyCtrlW)
	KeyCtrlX      = Key(tcell.KeyCtrlX)
	KeyCtrlY      = Key(tcell.KeyCtrlY)
	KeyCtrlZ      = Key(tcell.KeyCtrlZ)
	KeyBackspace  = Key(tcell.KeyBackspace)
	KeyBackspace2 = Key(tcell.KeyBackspace2)
	KeyTab        = Key(tcell.KeyTab)
	KeyEnter      = Key(tcell.KeyEnter)
	KeyEsc        = Key(tcell.KeyEscape)
	KeyPgdn       = Key(tcell.KeyPgDn)
	KeyPgup       = Key(tcell.KeyPgUp)
	MouseLeft     = Key(tcell.KeyF63) // arbitrary assignments
	MouseRight    = Key(tcell.KeyF62)
	MouseMiddle   = Key(tcell.KeyF61)
	MouseRelease  = Key(tcell.KeyF64)
	KeySpace      = Key(tcell.Key(' '))
)

// Modifiers.
const (
	ModAlt = Modifier(tcell.ModAlt)
)

func makeEvent(tev tcell.Event) Event {
	switch tev := tev.(type) {
	case *tcell.EventInterrupt:
		return Event{Type: EventInterrupt}
	case *tcell.EventResize:
		w, h := tev.Size()
		return Event{Type: EventResize, Width: w, Height: h}
	case *tcell.EventKey:
		k := tev.Key()
		ch := rune(0)
		if k == tcell.KeyRune {
			ch = tev.Rune()
			if ch == ' ' {
				k = tcell.Key(' ')
			} else {
				k = tcell.Key(0)
			}
		}
		mod := tev.Modifiers()
		return Event{
			Type: EventKey,
			Key:  Key(k),
			Ch:   ch,
			Mod:  Modifier(mod),
		}
	case *tcell.EventMouse:
		x, y := tev.Position()
		var k Key
		switch tev.Buttons() {
		case tcell.Button1:
			k = MouseLeft
			break
		case tcell.Button2:
			k = MouseMiddle
			break
		case tcell.Button3:
			k = MouseRight
			break
		case tcell.ButtonNone:
			k = MouseRelease
			break
		}

		return Event{
			Type:   EventMouse,
			MouseX: x,
			MouseY: y,
			Key:    k,
		}
	default:
		return Event{Type: EventNone}
	}
}

// ParseEvent is not supported.
func ParseEvent(data []byte) Event {
	// Not supported
	return Event{Type: EventError, Err: errors.New("no raw events")}
}

// PollRawEvent is not supported.
func PollRawEvent(data []byte) Event {
	// Not supported
	return Event{Type: EventError, Err: errors.New("no raw events")}
}

// PollEvent blocks until an event is ready, and then returns it.
func PollEvent() Event {
	ev := screen.PollEvent()
	return makeEvent(ev)
}

// Interrupt posts an interrupt event.
func Interrupt() {
	screen.PostEvent(tcell.NewEventInterrupt(nil))
}

// Cell represents a single character cell on screen.
type Cell struct {
	Ch rune
	Fg Attribute
	Bg Attribute
}
