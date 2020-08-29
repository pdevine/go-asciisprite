// Package sprite provides a framework for creating ASCII and Unicode based animations and games.
package sprite

import (
	tm "github.com/pdevine/go-asciisprite/termbox"
)

// A Sprite interface provides methods for initializing, updating, and rendering sprites.
type Sprite interface {
	Init()
	Update()
	Render()
	AddCostume(Costume)
	SetCostume(int)
	NextCostume()
	PrevCostume()
	TriggerEvent(string) bool
}

type Event struct {
	Callback func()
	Count    int
}

// A BaseSprite is a 2D sprite primitive.
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
	Events         map[string]*Event
}

// A SpriteGroup is a convenience method for holding groups of sprites.
type SpriteGroup struct {
	Sprites   []Sprite
	EventList []string
}

// NewBaseSprite creates a new BaseSprite from X and Y coordinates and a costume.
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

// AddCostume adds a costume to a BaseSprite.
func (s *BaseSprite) AddCostume(costume Costume) {
	s.Costumes = append(s.Costumes, &costume)
	if len(s.Costumes) == 1 {
		s.SetCostume(0)
	}
}

// SetCostume sets the current costume of a BaseSprite.
func (s *BaseSprite) SetCostume(c int) {
	// XXX - check for bounds here and return an error
	s.CurrentCostume = c
	s.Height = s.Costumes[s.CurrentCostume].Height
	s.Width = s.Costumes[s.CurrentCostume].Width
}

// Render draws the sprite to the screen buffer.
func (s *BaseSprite) Render() {
	if s.Visible {
		for _, b := range s.Costumes[s.CurrentCostume].Blocks {
			// call tcell screen.GetContent(b.X+s.X, b.Y+s.Y) here to see if we've already
			// written to the same location 
			tm.SetCell(b.X+s.X, b.Y+s.Y, b.Char, tm.Attribute(b.Fg), tm.Attribute(b.Bg))
		}
	}
}

// NextCostume changes a sprite's costume to the next costume.
func (s *BaseSprite) NextCostume() {
	// XXX - this should just call SetCostume()
	s.CurrentCostume++
	if s.CurrentCostume >= len(s.Costumes) {
		s.CurrentCostume = 0
	}
	s.Height = s.Costumes[s.CurrentCostume].Height
	s.Width = s.Costumes[s.CurrentCostume].Width
}

// PrevCostume changes a sprite's costume to the previous costume.
func (s *BaseSprite) PrevCostume() {
	s.CurrentCostume--
	if s.CurrentCostume < 0 {
		s.CurrentCostume = len(s.Costumes) - 1
	}
	s.Height = s.Costumes[s.CurrentCostume].Height
	s.Width = s.Costumes[s.CurrentCostume].Width
}

// Init provides a hook for initializing a sprite.
func (s *BaseSprite) Init() {
	// Init things
	s.Events = make(map[string]*Event)
}

// Update provides a hook for updating a sprite during the main loop.
func (s *BaseSprite) Update() {
	// Do things
}

// HitAtPoint reports whether a point on the screen intersects with this sprite.
func (s *BaseSprite) HitAtPoint(x, y int) bool {
	c := s.Costumes[s.CurrentCostume]
	if x >= s.X+c.LeftEdge() && x <= s.X+c.RightEdge() && y >= s.Y+c.TopEdge() && y <= s.Y+c.BottomEdge() {
		return true
	}
	return false
}

func (s *BaseSprite) RegisterEvent(name string, fn func()) {
	e := &Event{
		Callback: fn,
	}

	s.Events[name] = e
}

func (s *BaseSprite) TriggerEvent(name string) bool {
	e, ok := s.Events[name]
	if !ok {
		return false
	}
	e.Callback()
	return true
}

func (s *BaseSprite) RemoveEvent(name string) bool {
	_, ok := s.Events[name]
	if !ok {
		return false
	}
	s.Events[name] = nil
	return true
}

func (sg *SpriteGroup) TriggerEvent(name string) {
	sg.EventList = append(sg.EventList, name)
}

// Render draws each sprite in the SpriteGroup to the buffer.
func (sg *SpriteGroup) Render() {
	for _, s := range sg.Sprites {
		s.Render()
	}
	tm.Flush()
}

// Update updates each sprite in the SpriteGroup.
func (sg *SpriteGroup) Update() {
	// Consume any triggered events
	for _, e := range sg.EventList {
		for _, s := range sg.Sprites {
			s.TriggerEvent(e)
		}
	}
	sg.EventList = []string{}

	for _, s := range sg.Sprites {
		s.Update()
	}
}

// Remove removes a given sprite from the SpriteGroup.
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

// RemoveAll removes all sprites from the SpriteGroup.
func (sg *SpriteGroup) RemoveAll() {
	sg.Sprites = []Sprite{}
}
