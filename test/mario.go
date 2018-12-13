package main

import (
        "time"

        sprite "github.com/pdevine/go-asciisprite"
        //tm "github.com/gdamore/tcell/termbox"
        tm "github.com/pdevine/go-asciisprite/termbox"
)

var allSprites sprite.SpriteGroup
var Width int
var Height int

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
   bbbbttNt
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
 

type MarioState int

const (
	Walking MarioState = iota
	Jumping
)

type Mario struct {
	sprite.BaseSprite
	Timer   int
	TimeOut int
	State   MarioState
}

func (s *Mario) Update() {
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
	s.AddCostume(sprite.ColorConvert(mario_jump))
}

func (s *Mario) Walk() {
	s.State = Walking
	s.Costumes = []*sprite.Costume{}
	//s.AddCostume(sprite.ColorConvert(mario_turnaround))
	s.AddCostume(sprite.ColorConvert(mario_walk1))
	s.AddCostume(sprite.ColorConvert(mario_walk2))
	s.AddCostume(sprite.ColorConvert(mario_walk3))
	s.AddCostume(sprite.ColorConvert(mario_walk2))
}


/*
func (s *Mario) Reverse() {
	for _, c := range s.Costumes {
	}
}
*/

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

	m := &Mario{BaseSprite: sprite.BaseSprite{
		Visible: true},
		TimeOut: 2,
	}
	//m.AddCostume(sprite.ColorConvert(mario))
	m.Walk()
	allSprites.Sprites = append(allSprites.Sprites, m)

mainloop:
        for {
                tm.Clear(tm.ColorDefault, tm.ColorDefault)
		//tm.Clear(tm.ColorBlack, tm.ColorWhite)

                select {
                case ev := <-event_queue:
                        if ev.Type == tm.EventKey {
                                if ev.Key == tm.KeyEsc || ev.Ch == 'q' {
                                        break mainloop
                                }
				if ev.Ch == ' ' {
					if m.State == Walking {
						m.Jump()
					} else {
						m.Walk()
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

