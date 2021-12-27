package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Info struct {
	*tview.TextView
}

func NewInfo(msg string, doneFunc func(key tcell.Key)) *Info {
	textVies := tview.NewTextView().
		SetText(msg).
		SetDoneFunc(doneFunc)
	textVies.SetBorder(true).
		SetBorderPadding(5, 0, 5, 0).
		SetRect(150, 30, 60, 15)

	return &Info{textVies}
}

func NewInfoWithError(err error, doneFunc func(key tcell.Key)) *Info {
	return NewInfo(err.Error(), doneFunc)
}
