package main

import (
	"fmt"
	"pullyCat/cli"
)

func main() {
	menu := cli.NewMenu("Chose a colour")

	menu.AddItem("Red", "red")
	menu.AddItem("Blue", "blue")
	menu.AddItem("Green", "green")
	menu.AddItem("Yellow", "yellow")
	menu.AddItem("Cyan", "cyan")

	choice := menu.Display()

	fmt.Printf("Choice: %s\n", choice)
}
