package main

import (
	"fmt"
	"math"
	"strings"
	"time"

	sprite "github.com/pdevine/go-asciisprite"
	//tm "github.com/gdamore/tcell/termbox"
	tm "github.com/pdevine/go-asciisprite/termbox"
)

var allSprites sprite.SpriteGroup
var allBlocks []*Block
var Width int
var Height int

const BackgroundColor = 40

const big_mario = `
      RRRRR
    RRRRRRY
   RRRRRRYY
   RRRRRRRRRRR
   bbbttNttt
  bttbttNNtttt
  bttbbtttttttt
 bbttbbtttNtttt
 bbtttttNNNNNN
 bbbtttttNNNNN`

const mario = `   RRRRRR
  RRRRRRRRRR
  bbbtttNt
 btbttttNttt
 btbbttttNttt
 bbtttttNNNN
   tttttttt
  RRBRRRR
 RRRBRRBRRR
RRRRBBBBRRRR
ttRBYBBYBRtt
tttBBBBBBttt
ttBBBBBBBBtt
  BBB  BBB
 bbb    bbb
bbbb    bbbb`

const mario_walk1 = `     RRRRR
    RRRRRRRRR
    bbbttNt
   btbtttNttt
   btbbtttNttt
   bbttttNNNN
     ttttttt
    RRRRBR t
   tRRRRRRttt
  ttBRRRRRtt
  bbBBBBBBB
  bBBBBBBBB
 bbBBB BBB
 b    bbb
      bbbb`

const mario_walk2 = `     RRRRR
    RRRRRRRRR
    bbbttNt
   btbtttNttt
   btbbtttNttt
   bbttttNNNN
     ttttttt
    RRBRRR
   RRRRBBRR
   RRRBBYBBB
   RRRRBBBBB
    RRtttBBB
    BRttBBB
    BBBbbbb
    bbbbbbbb
    bbbb`

const mario_walk3 = `     RRRRR
    RRRRRRRRR
    bbbttNt
   btbtttNttt
   btbbtttNttt
   bbttttNNNN
     ttttttt
  RRRRBBRR
ttRRRRBBBRRRttt
ttt RRBYBBBRRtt
tt  BBBBBBB  b
   BBBBBBBBBbb
  BBBBBBBBBBbb
 bbBBB   BBBbb
 bbb
  bbb`

const mario_turnaround = `     RRRRR
   bRRRRRRRR
  bbbbbbtNt
 ttbttbtttttt
 ttbttbbttNNtt
  ttbttttttNN
   BBBRRRBtt
  BBtttRBBRRR
  BRtttRRRRRR
  BBBttRRRRRR
   BBBBBRRRR
   BbbbBBBB
    bbbbBBB
 b bBBbbbB
 bbbbbB
  bbbb`

const mario_jump = `             ttt
      RRRRR  ttt
     RRRRRRRRRtt
     bbbttNt RRR
    btbtttNttRRR
    btbbtttNtttR
    bbttttNNNNN
      tttttttR
  RRRRRBRRRBR
 RRRRRRRBRRRB  b
ttRRRRRRBBBBB  b
ttt BBRBBYBBBBbb
 t bBBBBBBBBBBbb
  bbbBBBBBBBBBbb
 bbbBBBBBBB
 b  BBBB`

const mushroom = `      oooo
     ooooOO
    ooooOOOO
   oooooOOOOO
  oooooooOOOoo
 ooOOOooooooooo
 oOOOOOoooooooo
ooOOOOOoooooOOoo
ooOOOOOoooooOOOo
oooOOOoooooooOOo
oooooooooooooooo
 oOOOwwwwwwOOOo
    wwwwwwww
    wwwwwwow
    wwwwwwow
     wwwwow`

const star = `       oo
      oooo
      oooo
     oooooo
 oooooooooooooo
 oooooOooOooooo
  ooooOooOoooo
   oooOooOooo
    oooooooo
    oooooooo
   oooooooooo
   oooooooooo
   oooo  oooo
  ooo      ooo
  oo        oo`

const flag = `
wwwwwwwwwwwwwwww
 wwwwwwwwGGGGGww
  wwwwwwGGwGwGGw
   wwwwwGwwGwwGw
    wwwwGwGGGwGw
     wwwGGGwGGGw
      wwGGGGGGGw
       wwwGGGwww
        wwwwwwww
         wwwwwww
          wwwwww
           wwwww
            wwww
             www
              ww
               w`

const flagpole_top = `







      NNNN
     NgGGGN
    NgGGGGGN
    NgGGGGGN
    NGGGGGGN
    NGGGGGGN
     NGGGGN
      NNNN`

const flagpole = `      gg
       gg
       gg
       gg
       gg
       gg
       gg
       gg
       gg
       gg
       gg
       gg
       gg
       gg
       gg
       gg`

type MarioState int

const (
	Standing MarioState = iota
	Walking
	Jumping
)

const FacingLeft = -1
const FacingRight = 1

type Mario struct {
	sprite.BaseSprite
	AX        float64
	VX        float64
	Timer     int
	TimeOut   int
	State     MarioState
	Direction int
}

func InitMario() *Mario {
	m := &Mario{BaseSprite: sprite.BaseSprite{
		Visible: true,
		X:       20,
		Y:       11 * 8},
		Direction: FacingRight,
		TimeOut:   2,
	}
	return m
}

func (s *Mario) Update() {
	s.VX = s.VX + s.AX
	s.AX = 0
	s.VX *= 0.85 // apply friction
	s.X += int(math.Round(s.VX))
	if s.X >= Width/2 {
		s.X = Width / 2
		for _, blk := range allBlocks {
			blk.X += -int(math.Round(s.VX))
		}
	}

	s.Timer++
	if s.Timer > s.TimeOut {
		s.CurrentCostume++
		if s.CurrentCostume >= len(s.Costumes) {
			s.CurrentCostume = 0
		}
		s.Timer = 0
	}
}

func (s *Mario) Jump() {
	s.State = Jumping
	s.Costumes = []*sprite.Costume{}
	s.AddCostume(sprite.ColorConvert(mario_jump, tm.Attribute(BackgroundColor)))
}

func (s *Mario) Walk(Direction int) {
	if s.State != Walking || s.Direction != Direction {
		s.State = Walking
		s.Direction = Direction
		s.Costumes = []*sprite.Costume{}
		//s.AddCostume(sprite.ColorConvert(mario_turnaround))
		bg := tm.Attribute(BackgroundColor)
		if Direction == FacingRight {
			s.AddCostume(sprite.ColorConvert(mario_walk1, bg))
			s.AddCostume(sprite.ColorConvert(mario_walk2, bg))
			s.AddCostume(sprite.ColorConvert(mario_walk3, bg))
			s.AddCostume(sprite.ColorConvert(mario_walk2, bg))
		} else {
			s.AddCostume(sprite.ColorConvert(reverseCostumeText(mario_walk1), bg))
			s.AddCostume(sprite.ColorConvert(reverseCostumeText(mario_walk2), bg))
			s.AddCostume(sprite.ColorConvert(reverseCostumeText(mario_walk3), bg))
			s.AddCostume(sprite.ColorConvert(reverseCostumeText(mario_walk2), bg))
		}
	}
}

func (s *Mario) MoveRight() {
	s.AX = 5
	s.VX = 0
	s.Walk(FacingRight)
}

func (s *Mario) MoveLeft() {
	s.AX = -5
	s.VX = 0
	s.Walk(FacingLeft)
}

func reverseCostumeText(s string) string {
	lines := strings.Split(s, "\n")
	var fstr string

	for _, l := range lines {
		chars := []rune(fmt.Sprintf("%-16v", l))
		for i, j := 0, len(chars)-1; i < j; i, j = i+1, j-1 {
			chars[i], chars[j] = chars[j], chars[i]
		}
		fstr += string(chars) + "\n"
	}
	return fstr
}

func main() {
	// XXX - Wait a bit until the terminal is properly initialized
	time.Sleep(500 * time.Millisecond)

	err := tm.Init()
	if err != nil {
		panic(err)
	}
	defer tm.Close()

	Width, Height = tm.Size()
	tm.SetOutputMode(tm.Output256)

	event_queue := make(chan tm.Event)
	go func() {
		for {
			event_queue <- tm.PollEvent()
		}
	}()

	bg := tm.Attribute(BackgroundColor)

	m := InitMario()
	m.Walk(FacingRight)

	allBlocks = ParseLevel(level1, bg)
	allSprites.Sprites = append(allSprites.Sprites, m)

mainloop:
	for {
		tm.Clear(tm.ColorDefault, bg)

		select {
		case ev := <-event_queue:
			if ev.Type == tm.EventKey {
				if ev.Key == tm.KeyEsc || ev.Ch == 'q' {
					break mainloop
				}
				if ev.Key == tm.KeyArrowRight {
					m.MoveRight()
				} else if ev.Key == tm.KeyArrowLeft {
					m.MoveLeft()
				}
				if ev.Ch == ' ' {
					if m.State == Walking {
						m.Jump()
					} else {
						m.Walk(m.Direction)
					}
				}
			} else if ev.Type == tm.EventResize {
				Width = ev.Width
				Height = ev.Height
			}
		default:
			allSprites.Update()
			allSprites.Render()
			time.Sleep(50 * time.Millisecond)
		}
	}
}
