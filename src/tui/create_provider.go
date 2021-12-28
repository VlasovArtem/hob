package tui

import (
	"github.com/VlasovArtem/hob/src/provider/model"
	"github.com/gdamore/tcell/v2"
	"github.com/google/uuid"
	"github.com/rivo/tview"
)

const CreateProviderPageName = "create-provider"

type createProviderReq struct {
	name, details string
}

type CreateProvider struct {
	*FlexApp
	*Navigation
	app *TerminalApp
}

func (c *CreateProvider) NavigationInfo(app *TerminalApp, variables map[string]any) *NavigationInfo {
	return NewNavigationInfo(CreateProviderPageName, func() tview.Primitive { return NewCreateProvider(app, variables["userId"].(uuid.UUID)) })
}

func (c *CreateProvider) enrichNavigation(app *TerminalApp, userId uuid.UUID) {
	c.Navigation = NewNavigation(app, c.NavigationInfo(app, map[string]any{"userId": userId}))
}

func NewCreateProvider(app *TerminalApp, userId uuid.UUID) *CreateProvider {
	f := &CreateProvider{
		app:     app,
		FlexApp: NewFlexApp(),
	}
	f.bindKeys()
	f.InitFlexApp(app)
	f.enrichNavigation(app, userId)

	var request createProviderReq

	form := tview.NewForm().
		AddInputField("Name", "", 20, nil, func(text string) { request.name = text }).
		AddInputField("Details", "", 20, nil, func(text string) { request.details = text }).
		AddButton("Create", f.create(userId, request)).
		AddButton("Cancel", f.BackFunc())

	form.SetBorder(true).SetTitle("Add Provider").SetTitleAlign(tview.AlignCenter).SetRect(150, 30, 60, 15)

	f.AddItem(form, 0, 8, true)

	f.SetInputCapture(f.KeyboardFunc)

	return f
}

func (c *CreateProvider) bindKeys() {
	c.Actions = KeyActions{
		tcell.KeyEscape: NewKeyAction("Back", c.KeyBack),
	}
}

func (c *CreateProvider) create(userId uuid.UUID, request createProviderReq) func() {
	return func() {
		paymentRequest := model.CreateProviderRequest{
			UserId:  userId,
			Name:    request.name,
			Details: request.details,
		}

		if _, err := c.app.GetProviderService().Add(paymentRequest); err != nil {
			c.ShowErrorTo(err)
		} else {
			c.ShowInfoReturnBack("Payment for the user with id '%s' successfully added.", userId)
		}
	}
}
