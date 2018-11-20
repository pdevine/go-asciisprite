package main

import (
	"fmt"
	sprite "github.com/pdevine/go-asciisprite"
)

func main() {
	fmt.Printf("\n%s\n", sprite.BuildString("game over"))
	fmt.Printf("\n%s\n", sprite.BuildString("the quick brown fox jumped over the lazy dog!/.,"))
	fmt.Printf("\n%s\n", sprite.BuildString("0123456789"))
	fmt.Printf("\n%s\n", sprite.BuildString("score 0123456789"))
}

