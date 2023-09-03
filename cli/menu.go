package cli

import (
	"fmt"
	tm "github.com/buger/goterm"
	"github.com/nsf/termbox-go"
)

type Menu struct {
	Prompt    string
	CursorPos int
	MenuItems []*MenuItem
}

type MenuItem struct {
	Text    string
	ID      string
	SubMenu *Menu
}

func NewMenu(prompt string) *Menu {
	return &Menu{
		Prompt:    prompt,
		MenuItems: make([]*MenuItem, 0),
	}
}

// AddItem will add a new menu option to the menu list
func (m *Menu) AddItem(option string, id string) *Menu {
	menuItem := &MenuItem{
		Text: option,
		ID:   id,
	}

	m.MenuItems = append(m.MenuItems, menuItem)
	return m
}

// renderMenuItems prints the menu item list.
// Setting redraw to true will re-render the options list with updated current selection.
func (m *Menu) renderMenuItems(redraw bool) {
	if redraw {
		// Move the cursor up n lines where n is the number of options, setting the new
		// location to start printing from, effectively redrawing the option list
		//
		// This is done by sending a VT100 escape code to the terminal
		// @see http://www.climagic.org/mirrors/VT100_Escape_Codes.html
		fmt.Printf("\033[%dA", len(m.MenuItems)-1)
	}

	for index, menuItem := range m.MenuItems {
		var newline = "\n"
		if index == len(m.MenuItems)-1 {
			// Adding a new line on the last option will move the cursor position out of range
			// For out redrawing
			newline = ""
		}

		menuItemText := menuItem.Text
		cursor := "  "
		if index == m.CursorPos {
			cursor = tm.Color("> ", tm.YELLOW)
			menuItemText = tm.Color(menuItemText, tm.YELLOW)
		}

		fmt.Printf("\r%s %s%s", cursor, menuItemText, newline)
	}
}

// Display will display the current menu options and awaits user selection
// It returns the users selected choice
func (m *Menu) Display() string {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer func() {
		termbox.Close()
		// Show cursor again.
		fmt.Printf("\033[?25h")
	}()

	fmt.Printf("%s\n", tm.Color(tm.Bold(m.Prompt)+":", tm.CYAN))

	m.renderMenuItems(false)

	// Turn the terminal cursor off
	fmt.Printf("\033[?25l")

	for {
		key := termbox.PollEvent().Key
		if key == termbox.KeyEsc {
			return ""
		} else if key == termbox.KeyEnter {
			menuItem := m.MenuItems[m.CursorPos]
			fmt.Println("\r")
			return menuItem.ID
		} else if key == termbox.KeyArrowUp {
			m.CursorPos = (m.CursorPos + len(m.MenuItems) - 1) % len(m.MenuItems)
			m.renderMenuItems(true)
		} else if key == termbox.KeyArrowDown {
			m.CursorPos = (m.CursorPos + 1) % len(m.MenuItems)
			m.renderMenuItems(true)
		}
	}
}
