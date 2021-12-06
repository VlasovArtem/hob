package tui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type dataInfo struct {
	header string
	data   string
	color  tcell.Color
}

type dataContent struct {
	currentRow int
	content    []*dataInfo
}

type SessionInfo struct {
	*tview.Table
	app *TerminalApp
}

func (s *SessionInfo) refresh() {
	s.Table.Clear()
	s.init()
}

func (s *SessionInfo) init() {
	content := NewDataContent()

	s.SetSelectable(false, false)
	s.SetTitle("Info")
	if s.app.AuthorizedUser != nil {
		content.
			add(NewDataInfo().withHeader("User:").withDataf("%s %s", s.app.AuthorizedUser.LastName, s.app.AuthorizedUser.FirstName)).
			add(NewDataInfo().withHeader("Email:").withData(s.app.AuthorizedUser.Email))
	} else {
		content.add(NewDataInfo().withData("User not selected").withColor(tcell.ColorDefault))
	}
	if s.app.House != nil {
		content.
			add(NewDataInfo().withData("House Info").withColor(tcell.ColorDefault)).
			add(NewDataInfo().withHeader("Name:").withData(s.app.House.Name)).
			add(NewDataInfo().withHeader("Country:").withData(s.app.House.CountryCode)).
			add(NewDataInfo().withHeader("City:").withData(s.app.House.City))
	} else {
		content.add(NewDataInfo().withData("House not selected").withColor(tcell.ColorDefault))
	}

	for row, dInfo := range content.content {
		s.SetCellSimple(row, 0, dInfo.header).
			SetCell(row, 1, tview.NewTableCell(dInfo.data).SetTextColor(dInfo.color))
	}
}

func NewSessionInfo(app *TerminalApp) *SessionInfo {
	info := &SessionInfo{
		Table: tview.NewTable(),
		app:   app,
	}

	info.init()

	return info
}

func NewDataInfo() *dataInfo {
	return &dataInfo{
		color: tcell.ColorGreen,
	}
}

func (d *dataInfo) withHeader(header string) *dataInfo {
	d.header = header
	return d
}

func (d *dataInfo) withData(data string) *dataInfo {
	d.data = data
	return d
}

func (d *dataInfo) withDataf(format string, a ...interface{}) *dataInfo {
	d.data = fmt.Sprintf(format, a...)
	return d
}

func (d *dataInfo) withColor(color tcell.Color) *dataInfo {
	d.color = color
	return d
}

func NewDataContent() *dataContent {
	return new(dataContent)
}

func (c *dataContent) add(info *dataInfo) *dataContent {
	c.content = append(c.content, info)
	c.currentRow++
	return c
}
