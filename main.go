package main

import (
	"fmt"

	"github.com/Jakobaune/minyr/yr"
)

func main() {
	fmt.Println("Hvilken funksjon vil du kjøre?(average/konverter eller exit)")

	var input string
	fmt.Scanln(&input)

	if input == "average" {
		yr.Average()

	} else if input == "konverter" {
		yr.Konverter()
	} else if input == "exit" {
		return
	} else {
		fmt.Println("Ugyldig input")
	}
}
