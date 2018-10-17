package sprite

import (
	"strings"

	tm "github.com/gdamore/tcell/termbox"
)

type Attribute uint16

type Block struct {
	Char rune
	Fg   tm.Attribute
	Bg   tm.Attribute
	X    int
	Y    int
}

type Costume struct {
	Blocks []*Block
	Width  int
	Height int
}

func NewCostume(t string, alpha rune) Costume {
	c := Costume{}
	c.ChangeCostume(t, alpha)
	return c
}

func (c *Costume) ChangeCostume(t string, alpha rune) {
	c.Blocks = []*Block{}
	var width int
	var height int

	for y, line := range strings.Split(t, "\n") {
		for x, ch := range line {
			if ch != alpha {
				b := &Block{
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

func (c *Costume) TopEdge() int {
	top := c.Blocks[0].Y
	for _, b := range c.Blocks[1:] {
		top = min(b.Y, top)
	}
	return top
}

func (c *Costume) LeftEdge() int {
	left := c.Blocks[0].X
	for _, b := range c.Blocks[1:] {
		left = min(b.X, left)
	}
	return left
}

func (c *Costume) RightEdge() int {
	right := c.Blocks[0].X
	for _, b := range c.Blocks[1:] {
		right = max(b.X, right)
	}
	return right
}

func (c *Costume) BottomEdge() int {
	bottom := c.Blocks[0].Y
	for _, b := range c.Blocks[1:] {
		bottom = max(b.Y, bottom)
	}
	return bottom
}

func (c *Costume) LeftEdgeByRow() map[int]int {
	t := make(map[int]int)
	for _, b := range c.Blocks {
		if _, ok := t[b.Y]; ok == false {
			t[b.Y] = b.X
		}
		t[b.Y] = min(t[b.Y], b.X)
	}
	return t
}

func (c *Costume) RightEdgeByRow() map[int]int {
	t := make(map[int]int)
	for _, b := range c.Blocks {
		t[b.Y] = max(t[b.Y], b.X)
	}
	return t
}

func (c *Costume) BottomEdgeByColumn() map[int]int {
	t := make(map[int]int)
	for _, b := range c.Blocks {
		t[b.X] = max(t[b.X], b.Y)
	}
	return t
}

func (c *Costume) TopEdgeByColumn() map[int]int {
	t := make(map[int]int)
	for _, b := range c.Blocks {
		if _, ok := t[b.X]; ok == false {
			t[b.X] = b.Y
		}
		t[b.X] = min(t[b.X], b.Y)
	}
	return t
}

type Sprite interface {
	Init()
	Update()
	Render()
	AddCostume(Costume)
	SetCostume(int)
	NextCostume()
	PrevCostume()
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
	Dead           bool
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

func min(x, y int) int {
	if x < y {
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
		Dead:           false,
	}
	s.AddCostume(costume)
	return s
}

func (s *BaseSprite) AddCostume(costume Costume) {
	s.Costumes = append(s.Costumes, costume)
	if len(s.Costumes) == 1 {
		s.SetCostume(0)
	}
}

func (s *BaseSprite) SetCostume(c int) {
	s.CurrentCostume = c
	s.Height = s.Costumes[s.CurrentCostume].Height
	s.Width = s.Costumes[s.CurrentCostume].Width
}

func (s *BaseSprite) Render() {
	if s.Visible {
		for _, b := range s.Costumes[s.CurrentCostume].Blocks {
			tm.SetCell(b.X+s.X, b.Y+s.Y, b.Char, tm.Attribute(b.Fg), tm.Attribute(b.Bg))
		}
	}
}

func (s *BaseSprite) NextCostume() {
	s.CurrentCostume++
	if s.CurrentCostume >= len(s.Costumes) {
		s.CurrentCostume = 0
	}
	s.Height = s.Costumes[s.CurrentCostume].Height
	s.Width = s.Costumes[s.CurrentCostume].Width
}

func (s *BaseSprite) PrevCostume() {
	s.CurrentCostume--
	if s.CurrentCostume < 0 {
		s.CurrentCostume = len(s.Costumes) - 1
	}
	s.Height = s.Costumes[s.CurrentCostume].Height
	s.Width = s.Costumes[s.CurrentCostume].Width
}

func (s *BaseSprite) Init() {
	// Init things
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

func (sg *SpriteGroup) RemoveAll() {
	sg.Sprites = []Sprite{}
}
