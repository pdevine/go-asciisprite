package sprite

import (
	"strings"

	tm "github.com/pdevine/go-asciisprite/termbox"
)

// A Background interface provides methods for initializing, updating, and rendering sprites.
type Background interface {
	Init()
	Update()
	Render()
	AddBackground(t string)
}

// A BaseBackground is a 2D background which sits behind Sprites.
type BaseBackground struct {
	X          int
	Y          int
	Height     int
	Width      int
	Tiled      bool
	Background []Block
}

// Init provides a hook for initializing a background.
func (s *BaseBackground) Init() {
	// Init things
}

// Update provides a hook for updating a background during the main loop.
func (s *BaseBackground) Update() {
	// Do things
}

// Render draws background to the buffer.
func (s *BaseBackground) Render() {
	for _, b := range s.Background {
		tm.SetCell(b.X+s.X, b.Y+s.Y, b.Char, b.Fg, b.Bg)
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
