package sprite

import (
	"strings"
)

var Blocks = map[int]rune{
   0: ' ',
   1: '▘',
   2: '▝' ,
   3: '▀',
   4: '▖',
   5: '▌',
   6: '▞',
   7: '▛',
   8: '▗',
   9: '▚',
  10: '▐',
  11: '▜',
  12: '▄',
  13: '▙',
  14: '▟',
  15: '█',
}

func Convert(s string) Costume {
	blocks := []*Block{}
	l := strings.Split(s, "\n")
	maxR := len(l) + len(l)%2

	// all block sprites must be even
	m := make([][]rune, maxR, maxR)

	var maxC int
	for _, r := range l {
		maxC = max(maxC, len(r) + len(r)%2)
	}

	for rcnt, r := range l {
		m[rcnt] = make([]rune, maxC, maxC)
		for ccnt, c := range r {
			if c != ' ' {
				m[rcnt][ccnt] = c
			}
		}	
	}

	// make certain we make a row for any added space
	if len(l) < maxR {
		m[maxR-1] = make([]rune, maxC, maxC)
	}

	for rcnt := 0; rcnt < len(m); rcnt+=2 {
		// XXX - needs to be max(len(m[rcnt]), len(m[rcnt+1]))
		// for ccnt := 0; ccnt < max(len(m[rcnt]), len(m[rcnt+1])); ccnt+=2 {
		for ccnt := 0; ccnt < len(m[rcnt]); ccnt+=2 {
			c := 0
			if m[rcnt][ccnt] != 0 {
				c += 1
			}
			if len(m[rcnt]) > ccnt+1 && m[rcnt][ccnt+1] != 0 {
				c += 2
			}
			if len(m) > rcnt+1 && m[rcnt+1][ccnt] != 0 {
				c += 4
			}
			if len(m) > rcnt+1 && len(m[rcnt]) > ccnt+1 && m[rcnt+1][ccnt+1] == 'X' {
				c += 8
			}

			if c > 0 {
				b := &Block{
					Char: Blocks[c],
					X:    ccnt/2,
					Y:    rcnt/2,
				}
				blocks = append(blocks, b)
			}
		}
	}

	costume := Costume{Blocks: blocks}

	return costume
}
