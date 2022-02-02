package main

import (
	"os"

	"image/png"
)

func main() {
	f, err := os.Open("dog.png")
	if err != nil {
		//
	}

	img, err := png.Decode(f)
	if err != nil {
		//
	}

}
