package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Search struct {
	*tview.InputField
	app *TerminalApp
}

func NewSearch(app *TerminalApp) *Search {
	inputField := tview.NewInputField()
	inputField.SetFieldTextColor(tcell.Color(0))
	inputField.SetFieldBackgroundColor(tcell.Color(0))
	inputField.SetBorder(true)
	return &Search{
		InputField: inputField,
		app:        app,
	}
}
