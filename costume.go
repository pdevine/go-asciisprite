package sprite

import (
	"strings"

	tm "github.com/pdevine/go-asciisprite/termbox"
)

// Block holds a rune at (X,Y).
type Block struct {
	Char rune
	Fg   tm.Attribute
	Bg   tm.Attribute
	X    int
	Y    int
}

// Costume contains a slice of Block pointers.
type Costume struct {
	Blocks []*Block
	Width  int
	Height int
}

// NewCostume provides a new Costume from a string.
func NewCostume(t string, alpha rune) Costume {
	c := Costume{}
	c.ChangeCostume(t, alpha)
	return c
}

// ChangeCostume replaces a Costume with a new Costume from a string.
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

// TopEdge returns the Y value of the heighest Block in a Costume.
func (c *Costume) TopEdge() int {
	top := c.Blocks[0].Y
	for _, b := range c.Blocks[1:] {
		top = min(b.Y, top)
	}
	return top
}

// LeftEdge returns the X value of the furthest left Block in a Costume.
func (c *Costume) LeftEdge() int {
	left := c.Blocks[0].X
	for _, b := range c.Blocks[1:] {
		left = min(b.X, left)
	}
	return left
}

// RightEdge returns the X value of the furthest right Block in a Costume.
func (c *Costume) RightEdge() int {
	right := c.Blocks[0].X
	for _, b := range c.Blocks[1:] {
		right = max(b.X, right)
	}
	return right
}

// BottomEdge returns the Y value of the lowest Block in a Costume.
func (c *Costume) BottomEdge() int {
	bottom := c.Blocks[0].Y
	for _, b := range c.Blocks[1:] {
		bottom = max(b.Y, bottom)
	}
	return bottom
}

// LeftEdgeByRow returns a map of the left most X value in each row of a Costume.
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

// RightEdgeByRow returns a map of the right most X value in each row of a Costume.
func (c *Costume) RightEdgeByRow() map[int]int {
	t := make(map[int]int)
	for _, b := range c.Blocks {
		t[b.Y] = max(t[b.Y], b.X)
	}
	return t
}

// BottomEdgeByColumn returns a map of the lowest Y value in each column of a Costume.
func (c *Costume) BottomEdgeByColumn() map[int]int {
	t := make(map[int]int)
	for _, b := range c.Blocks {
		t[b.X] = max(t[b.X], b.Y)
	}
	return t
}

// TopEdgeByColumn returns a map of the heighest Y value in each column of a Costume.
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
