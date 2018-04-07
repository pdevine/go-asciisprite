package sprite

import (
	"strings"

	tm "github.com/nsf/termbox-go"
)

type Sprite interface {
	Update()
	Render()
	AddCostume(Costume)
}

type Costume struct {
	Text string
}

type BaseSprite struct {
	X              int
	Y              int
	Height         int
	Width          int
	Costumes       []Costume
	Alpha          rune
	CurrentCostume int
}

type SpriteGroup struct {
	Sprites []Sprite
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func NewBaseSprite(x, y int, costume Costume) *BaseSprite {
	s := &BaseSprite{
		X:              x,
		Y:              y,
		Height:         0,
		Width:          0,
		Costumes:       []Costume{},
		CurrentCostume: 0,
	}
	s.AddCostume(costume)
	return s
}

func (s *BaseSprite) AddCostume(costume Costume) {
	c := strings.Split(costume.Text, "\n")
	h := len(c)
	w := 0
	for _, l := range c {
		w = max(w, len(l))
	}
	s.Costumes = append(s.Costumes, costume)
	s.Height = h
	s.Width = w
}

func (s *BaseSprite) Render() {
	for y, line := range strings.Split(s.Costumes[s.CurrentCostume].Text, "\n") {
		for x, ch := range line {
			if ch != s.Alpha {
				tm.SetCell(s.X+x, s.Y+y, ch, tm.ColorWhite, tm.ColorBlack)
			}
		}
	}
}

func (s *BaseSprite) Update() {
	// Do things
}

func (sg *SpriteGroup) Render() {
	for _, s := range sg.Sprites {
		s.Render()
	}
	tm.Flush()
}

func (sg *SpriteGroup) Update() {
	for _, s := range sg.Sprites {
		s.Update()
	}
}
