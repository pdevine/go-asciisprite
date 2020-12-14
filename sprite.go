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
	BlockRender(*Surface)
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
	BlockCostumes  []*Surface
	Alpha          rune
	Visible        bool
	CurrentCostume int
	Dead           bool
	Events         map[string]*Event
}

// A SpriteGroup is a convenience method for holding groups of sprites.
type SpriteGroup struct {
	Sprites    []Sprite
	EventList  []string
	BlockMode  bool
	Background tm.Attribute
	bg         *Surface
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
		BlockCostumes:  []*Surface{},
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
	if len(s.BlockCostumes) > 0 {
		s.Height = s.BlockCostumes[s.CurrentCostume].Height
		s.Width = s.BlockCostumes[s.CurrentCostume].Width
	} else if len(s.Costumes) > 0 {
		s.Height = s.Costumes[s.CurrentCostume].Height
		s.Width = s.Costumes[s.CurrentCostume].Width
	}
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

func (s *BaseSprite) BlockRender(bg *Surface) {
	if s.Visible {
		if len(s.BlockCostumes) > s.CurrentCostume {
			surf := s.BlockCostumes[s.CurrentCostume]
			bg.Blit(*surf, s.X, s.Y)
		}
	}
}

// NextCostume changes a sprite's costume to the next costume.
func (s *BaseSprite) NextCostume() {
	s.CurrentCostume++
	if len(s.BlockCostumes) > 0 {
		if s.CurrentCostume >= len(s.BlockCostumes) {
			s.CurrentCostume = 0
		}
	} else if len(s.Costumes) > 0 {
		if s.CurrentCostume >= len(s.Costumes) {
			s.CurrentCostume = 0
		}
	}
	// Set the Width/Height
	s.SetCostume(s.CurrentCostume)
}

// PrevCostume changes a sprite's costume to the previous costume.
func (s *BaseSprite) PrevCostume() {
	s.CurrentCostume--
	if s.CurrentCostume < 0 {
		if len(s.BlockCostumes) > 0 {
			s.CurrentCostume = len(s.BlockCostumes)-1
		} else {
			s.CurrentCostume = len(s.Costumes)-1
		}
	}
	// Set the Width/Height
	s.SetCostume(s.CurrentCostume)
}

// Init provides a hook for initializing a sprite.
func (s *BaseSprite) Init() {
	// Init things
	s.Events = make(map[string]*Event)
	if len(s.BlockCostumes) > 0 {
		s.Height = s.BlockCostumes[s.CurrentCostume].Height
		s.Width = s.BlockCostumes[s.CurrentCostume].Width
	} else if len(s.BlockCostumes) > 0 {
		s.Height = s.Costumes[s.CurrentCostume].Height
		s.Width = s.Costumes[s.CurrentCostume].Width
	}
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

func (s *BaseSprite) HitAtPointSurface(x, y int) bool {
	surf := s.BlockCostumes[s.CurrentCostume]
	if x >= s.X && x <= s.X+surf.Width && y >= s.Y && y <= s.Y+surf.Height {
		return true
	}
	return false
}


// RegisterEvent registers a callback function with a name in this sprite.
func (s *BaseSprite) RegisterEvent(name string, fn func()) {
	e := &Event{
		Callback: fn,
	}

	s.Events[name] = e
}

// TriggerEvent causes a previously registered callback function to be called.
func (s *BaseSprite) TriggerEvent(name string) bool {
	e, ok := s.Events[name]
	if !ok {
		return false
	}
	e.Callback()
	return true
}

// RemoveEvent removes an event with a given name from this sprite.
func (s *BaseSprite) RemoveEvent(name string) bool {
	_, ok := s.Events[name]
	if !ok {
		return false
	}
	s.Events[name] = nil
	return true
}

func (sg *SpriteGroup) Init(width, height int, blockMode bool) {
	if blockMode {
		sg.BlockMode = blockMode
		sg.Background = tm.ColorDefault
		surf := NewSurface(width, height, false)
		sg.bg = &surf

	}
}

func (sg *SpriteGroup) Resize(width, height int) {
	if sg.BlockMode {
		surf := NewSurface(width, height, false)
		sg.bg = &surf
	}
}

// TriggerEvent causes causes all events of this name to be called.
func (sg *SpriteGroup) TriggerEvent(name string) {
	sg.EventList = append(sg.EventList, name)
}

// Render draws each sprite in the SpriteGroup to the buffer.
func (sg *SpriteGroup) Render() {
	if sg.BlockMode {
		sg.bg.Clear()
		for _, s := range sg.Sprites {
			s.BlockRender(sg.bg)
		}
		c := sg.bg.ConvertToColorCostume(sg.Background)
		for _, b := range c.Blocks {
			tm.SetCell(b.X, b.Y, b.Char, tm.Attribute(b.Fg), tm.Attribute(b.Bg))
		}
	} else {
		for _, s := range sg.Sprites {
			s.Render()
		}

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
