package main

import (
	"fmt"
	"pullyCat/cli"
)

func main() {

	menu := cli.NewMultiSelectMenu("Chose a colour", -1, -1)

	menu.AddOption("Red", "red")
	menu.AddOption("Blue", "blue")
	menu.AddOption("Green", "green")
	menu.AddOption("Yellow", "yellow")
	menu.AddOption("Cyan", "cyan")

	choice := menu.Display()

	fmt.Printf("Choice: %s\n", choice)
}
