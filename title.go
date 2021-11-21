package sprite

import (
	"math"
)

const frame_ul = `XXXXXXXX
XXXXXXXX
XX
XX
XX
XX      `

const frame_ur = `XXXXXXXX
XXXXXXXX
      XX
      XX
      XX
      XX`

const frame_ll = `XX
XX
XX
XX
XXXXXXXX
XXXXXXXX`

const frame_lr = `      XX
      XX
      XX
      XX
XXXXXXXX
XXXXXXXX`

const arrow_ul = `
XXXXX
XXX
X XX
X  XX
    XX `

const arrow_ur = `
   XXXXX
     XXX
    XX X
   XX  X
  XX    `

const arrow_ll = `
    XX
X  XX
X XX
XXX
XXXXX   `

const arrow_lr = `
  XX
   XX  X
    XX X
     XXX
   XXXXX`

type EdgeType int

const (
	UpperLeftEdge EdgeType = iota
	UpperRightEdge
	LowerLeftEdge
	LowerRightEdge
)

type TitleArrow struct {
	BaseSprite
	DX    float64
	DY    float64
	Angle float64
	Type  EdgeType
}

type TitleEdge struct {
	BaseSprite
}

type TitleScreen struct {
	Sprites SpriteGroup
}

// NewTitleEdge create an edge marker in the TitleScreen.
func NewTitleEdge(t EdgeType, r Rect) *TitleEdge {
	s := &TitleEdge{BaseSprite: BaseSprite{
		Visible: true},
	}

	switch t {
	case UpperLeftEdge:
		s.X = r.X + 1
		s.Y = r.Y + 1
		s.AddCostume(Convert(frame_ul))
	case UpperRightEdge:
		s.X = r.X + r.W - 1
		s.Y = r.Y + 1
		s.AddCostume(Convert(frame_ur))
	case LowerLeftEdge:
		s.X = r.X + 1
		s.Y = r.Y + r.H - 1
		s.AddCostume(Convert(frame_ll))
	case LowerRightEdge:
		s.X = r.X + r.W - 1
		s.Y = r.Y + r.H - 1
		s.AddCostume(Convert(frame_lr))
	}
	return s
}

// NewTitleArrow creates a TitleArrow in the TitleScreen.
func NewTitleArrow(t EdgeType, r Rect) *TitleArrow {
	s := &TitleArrow{BaseSprite: BaseSprite{
		Visible: true},
		Type: t,
	}
	switch s.Type {
	case UpperLeftEdge:
		s.DX = float64(r.X + 4)
		s.DY = float64(r.Y + 3)
		s.AddCostume(Convert(arrow_ul))
	case UpperRightEdge:
		s.DX = float64(r.X + r.W - 4)
		s.DY = float64(r.Y + 3)
		s.AddCostume(Convert(arrow_ur))
	case LowerLeftEdge:
		s.DX = float64(r.X + 4)
		s.DY = float64(r.Y + r.H - 3)
		s.AddCostume(Convert(arrow_ll))
	case LowerRightEdge:
		s.DX = float64(r.X + r.W - 4)
		s.DY = float64(r.Y + r.H - 3)
		s.AddCostume(Convert(arrow_lr))
	}
	return s
}

// Update moves a TitleArrow in a TitleScreen.
func (s *TitleArrow) Update() {
	s.Angle += 0.25

	d := math.Sin(s.Angle) * 0.2
	switch s.Type {
	case UpperLeftEdge:
		s.DX += d
		s.DY += d
	case UpperRightEdge:
		s.DX -= d
		s.DY += d
	case LowerLeftEdge:
		s.DX += d
		s.DY -= d
	case LowerRightEdge:
		s.DX -= d
		s.DY -= d
	}
	s.X = int(math.Round(s.DX))
	s.Y = int(math.Round(s.DY))

}

// InitTitleScreen creates a default TitleScreen.
func InitTitleScreen(r Rect) *TitleScreen {
	title := TitleScreen{
		Sprites: SpriteGroup{},
	}

	txt := "ADJUST YOUR TERMINAL TO SEE ALL OF THE EDGES OF THE PLAY AREA"
	adj_txt := &BaseSprite{
		X:       r.X + r.W/2 - len(txt)/2,
		Y:       22,
		Visible: true,
	}
	adj_txt.AddCostume(NewCostume(txt, '@'))

	txt = "Recommended Font:  Menlo 8pt (0.81 Line Spacing)"
	font_txt := &BaseSprite{
		X:       r.X + r.W/2 - len(txt)/2,
		Y:       24,
		Visible: true,
	}
	font_txt.AddCostume(NewCostume(txt, '@'))

	for _, et := range []EdgeType{UpperLeftEdge, UpperRightEdge, LowerLeftEdge, LowerRightEdge} {
		a := NewTitleArrow(et, r)
		e := NewTitleEdge(et, r)
		title.Sprites.Sprites = append(title.Sprites.Sprites, e)
		title.Sprites.Sprites = append(title.Sprites.Sprites, a)

	}

	title.Sprites.Sprites = append(title.Sprites.Sprites, adj_txt)
	title.Sprites.Sprites = append(title.Sprites.Sprites, font_txt)
	return &title
}

// Update updates all of the Sprites in a TitleScreen.
func (t *TitleScreen) Update() {
	for _, s := range t.Sprites.Sprites {
		s.Update()
	}
}

// Render renders all of the Sprites in a TitleScreen.
func (t *TitleScreen) Render() {
	for _, s := range t.Sprites.Sprites {
		s.Render()
	}
}
