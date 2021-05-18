package main

import (
	"fmt"
	sprite "github.com/pdevine/go-asciisprite"
)

func main() {
	f := sprite.NewPakuFont()
	fmt.Printf("\n%s\n", f.BuildString("game over"))
	fmt.Printf("\n%s\n", f.BuildString("the quick brown fox jumped over the lazy dog!/.,"))
	fmt.Printf("\n%s\n", f.BuildString("0123456789"))
	fmt.Printf("\n%s\n", f.BuildString("score 0123456789"))

	f2 := sprite.NewJRFont()
	fmt.Printf("\n%s\n", f2.BuildString("ABCDEFGHIJKLM"))
	fmt.Printf("\n%s\n", f2.BuildString("NOPQRSTUVWXYZ"))
	fmt.Printf("\n%s\n", f2.BuildString("abcdefghijklm"))
	fmt.Printf("\n%s\n", f2.BuildString("nopqrstuvwxyz"))
	fmt.Printf("\n%s\n", f2.BuildString("0123456789"))
	fmt.Printf("\n%s\n", f2.BuildString(".,;:+- */%!?"))
	fmt.Printf("\n%s\n", f2.BuildString("<>∨^[](){}\\"))
	fmt.Printf("\n%s\n", f2.BuildString("#&='@|ß☆_~$"))

	f3 := sprite.New90sFont()
	fmt.Printf("\n%s\n", f3.BuildString("the quick brown fox jumped over the lazy dog!/.,"))
	fmt.Printf("\n%s\n", f3.BuildString("THE QUICK BROWN FOX JUMPED OVER THE LAZY DOG!/.,"))
	fmt.Printf("\n%s\n", f3.BuildString("0123456789"))
}

