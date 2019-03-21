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

type Rect struct {
	X int
	Y int
	W int
	H int
}

type Point struct {
	X int
	Y int
}
