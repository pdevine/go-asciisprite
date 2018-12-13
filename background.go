package sprite

import (
	"strings"

	//tm "github.com/gdamore/tcell/termbox"
	tm "github.com/pdevine/go-asciisprite/termbox"
)

type Background interface {
	Init()
	Update()
	Render()
	AddBackground(t string)
}

type BaseBackground struct {
	X          int
	Y          int
	Height     int
	Width      int
	Tiled      bool
	Background []Block
}

func (s *BaseBackground) Init() {
	// Init things
}

func (s *BaseBackground) Update() {
	// Do things
}

func (s *BaseBackground) Render() {
	for _, b := range s.Background {
		tm.SetCell(b.X+s.X, b.Y+s.Y, b.Char, tm.ColorWhite, tm.ColorBlack)
	}
}

func (s *BaseBackground) AddBackground(t string) {
	s.Background = []Block{}
	var width int
	var height int

	for y, line := range strings.Split(t, "\n") {
		for x, ch := range line {
			b := Block{
				Char: ch,
				X:    x,
				Y:    y,
			}
			s.Background = append(s.Background, b)
			width = max(x, width)
		}
		height = y
	}
	s.Width = width
	s.Height = height
}

func (s *BaseBackground) AddCostume(costume Costume) {
}
