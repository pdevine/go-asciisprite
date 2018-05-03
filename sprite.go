package sprite

import (
	"strings"

	tm "github.com/nsf/termbox-go"
)

type Block struct {
	Char rune
	X    int
	Y    int
}

type Costume struct {
	Blocks []Block
	Width  int
	Height int
}

func NewCostume(t string, alpha rune) Costume {
	c := Costume{}
	c.ChangeCostume(t, alpha)
	return c
}

func (c *Costume) ChangeCostume(t string, alpha rune) {
	c.Blocks = []Block{}
	var width int
	var height int

	for y, line := range strings.Split(t, "\n") {
		for x, ch := range line {
			if ch != alpha {
				b := Block{
					Char: ch,
					X:    x,
					Y:    y,
				}
				c.Blocks = append(c.Blocks, b)
				width = max(x, width)
			}
		}
		height = y
	}
	c.Width = width + 1
	c.Height = height
}

type Sprite interface {
	Update()
	Render()
	AddCostume(Costume)
}

type BaseSprite struct {
	X              int
	Y              int
	Height         int
	Width          int
	Costumes       []Costume
	Alpha          rune
	Visible        bool
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
		Visible:        true,
		Costumes:       []Costume{},
		CurrentCostume: 0,
	}
	s.AddCostume(costume)
	return s
}

func (s *BaseSprite) AddCostume(costume Costume) {
	s.Costumes = append(s.Costumes, costume)
	s.Height = s.Costumes[s.CurrentCostume].Height
	s.Width = s.Costumes[s.CurrentCostume].Width
}

func (s *BaseSprite) Render() {
	if s.Visible {
		for _, b := range s.Costumes[s.CurrentCostume].Blocks {
			tm.SetCell(b.X+s.X, b.Y+s.Y, b.Char, tm.ColorWhite, tm.ColorBlack)
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

func (sg *SpriteGroup) Remove(s Sprite) {
	var idx int
	for cnt, tSprite := range sg.Sprites {
		if s == tSprite {
			idx = cnt
			break
		}
	}
	copy(sg.Sprites[idx:], sg.Sprites[idx+1:])
	sg.Sprites[len(sg.Sprites)-1] = nil
	sg.Sprites = sg.Sprites[:len(sg.Sprites)-1]
}
