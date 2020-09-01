package sprite

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

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// A Rect provides a rectangle starting at (X,Y) with W width and H height
type Rect struct {
	X int
	Y int
	W int
	H int
}

// A Point provides a point on the screen at (X,Y)
type Point struct {
	X int
	Y int
}
