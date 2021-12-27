package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var promptPage = "prompt"

type ModalButton struct {
	Name   string
	Action func()
}

func ShowModal(p *tview.Pages, msg string, buttons []ModalButton) {
	var modalButtons []string
	buttonToAction := make(map[string]func())
	for _, button := range buttons {
		modalButtons = append(modalButtons, button.Name)
		buttonToAction[button.Name] = button.Action
	}
	modalButtons = append(modalButtons, "Cancel")
	buttonToAction["Cancel"] = func() { p.RemovePage(promptPage) }

	modal := tview.NewModal().
		AddButtons(modalButtons).
		SetTextColor(tcell.ColorDarkSlateGray).
		SetText(msg).
		SetDoneFunc(func(_ int, buttonLabel string) {
			buttonToAction[buttonLabel]()
		})

	p.AddPage(promptPage, modal, false, false)
	p.ShowPage(promptPage)
}
