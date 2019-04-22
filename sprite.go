package sprite

import (
	//"fmt"
	//tm "github.com/gdamore/tcell/termbox"
	tm "github.com/pdevine/go-asciisprite/termbox"
)

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
	Costumes       []*Costume
	Alpha          rune
	Visible        bool
	CurrentCostume int
	Dead           bool
}

type SpriteGroup struct {
	Sprites []Sprite
}

func NewBaseSprite(x, y int, costume Costume) *BaseSprite {
	s := &BaseSprite{
		X:              x,
		Y:              y,
		Height:         0,
		Width:          0,
		Visible:        true,
		Costumes:       []*Costume{},
		CurrentCostume: 0,
		Dead:           false,
	}
	s.AddCostume(costume)
	return s
}

func (s *BaseSprite) AddCostume(costume Costume) {
	s.Costumes = append(s.Costumes, &costume)
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
			// call tcell screen.GetContent(b.X+s.X, b.Y+s.Y) here to see if we've already
			// written to the same location 
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

func (s *BaseSprite) HitAtPoint(x, y int) bool {
	c := s.Costumes[s.CurrentCostume]
	if x >= s.X+c.LeftEdge() && x <= s.X+c.RightEdge() && y >= s.Y+c.TopEdge() && y <= s.Y+c.BottomEdge() {
		return true
	}
	return false
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
	var found bool
	for cnt, tSprite := range sg.Sprites {
		if s == tSprite {
			idx = cnt
			found = true
			break
		}
	}
	if found {
		copy(sg.Sprites[idx:], sg.Sprites[idx+1:])
		sg.Sprites[len(sg.Sprites)-1] = nil
		sg.Sprites = sg.Sprites[:len(sg.Sprites)-1]
	}
}

func (sg *SpriteGroup) RemoveAll() {
	sg.Sprites = []Sprite{}
}
