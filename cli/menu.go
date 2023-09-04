package cli

import (
	"fmt"
	tm "github.com/buger/goterm"
	"github.com/nsf/termbox-go"
)

type MenuType int

const (
	MenuTypeSelect       MenuType = 1
	MenuTypeMultiSelect  MenuType = 2
	MenuTypeConfirmation MenuType = 3
)

type Menu struct {
	MenuType  MenuType
	Prompt    string
	CursorPos int
	Options   []MenuOption
	min       int
	max       int
}

type MenuOption struct {
	ID       string
	Text     string
	selected bool
	// subMenu  *Menu
}

func NewSelectMenu(prompt string, options ...MenuOption) *Menu {
	if options == nil {
		options = make([]MenuOption, 0)
	}
	return &Menu{MenuType: MenuTypeSelect, Prompt: prompt, CursorPos: 0, Options: options, min: 1, max: 1}
}

func NewMultiSelectMenu(prompt string, min int, max int, options ...MenuOption) *Menu {
	if options == nil {
		options = make([]MenuOption, 0)
	}
	return &Menu{MenuType: MenuTypeMultiSelect, Prompt: prompt, CursorPos: 0, Options: options, min: min, max: max}
}

func NewConfirmationMenu(prompt string, options ...MenuOption) *Menu {
	if options == nil {
		options = make([]MenuOption, 0)
	}
	return &Menu{MenuType: MenuTypeConfirmation, Prompt: prompt, CursorPos: 0, Options: options, min: 1, max: 1}
}

// AddOption will add a new menu option to the select menu list
func (m *Menu) AddOption(id string, text string) {
	optionItem := MenuOption{ID: id, Text: text, selected: false}
	m.Options = append(m.Options, optionItem)
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
		fmt.Printf("\033[%dA", len(m.Options)-1)
	}

	for index, menuItem := range m.Options {
		var newline = "\n"
		if index == len(m.Options)-1 {
			// Adding a new line on the last option will move the cursor position out of range
			// For out redrawing
			newline = ""
		}

		menuItemText := menuItem.Text
		cursor := "  "
		// single select menu
		if m.MenuType == MenuTypeSelect {
			if index == m.CursorPos {
				cursor = tm.Color("> ", tm.YELLOW)
				menuItemText = tm.Color(menuItemText, tm.YELLOW)
			}
			// multi select menu
		} else if m.MenuType == MenuTypeMultiSelect {
			if menuItem.selected {
				cursor = fmt.Sprintf("[%s] ", tm.Color("*", tm.BLUE))
			} else {
				cursor = "[ ] "
			}

			if index == m.CursorPos {
				menuItemText = tm.Color(menuItemText, tm.BLUE)
			}
			// confirmation menu
		} else if m.MenuType == MenuTypeConfirmation {
			// TODO
		} else {
			// TODO Error
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

		switch termbox.PollEvent().Key {
		case termbox.KeyEsc:
			return ""
		case termbox.KeyEnter:
			menuItem := m.Options[m.CursorPos]
			fmt.Println("\r")
			return menuItem.ID
		case termbox.KeySpace:
			if m.MenuType == MenuTypeMultiSelect {
				m.Options[m.CursorPos].selected = !m.Options[m.CursorPos].selected
				m.renderMenuItems(true)
			}
			//menuItem := m.Options[m.CursorPos]
			//fmt.Println("\r")
		case termbox.KeyArrowUp:
			m.CursorPos = (m.CursorPos + len(m.Options) - 1) % len(m.Options)
			m.renderMenuItems(true)
		case termbox.KeyArrowDown:
			m.CursorPos = (m.CursorPos + 1) % len(m.Options)
			m.renderMenuItems(true)
		}
	}
}
