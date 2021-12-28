package tui

import (
	"fmt"
	"github.com/VlasovArtem/hob/src/provider/model"
	"github.com/gdamore/tcell/v2"
	"github.com/google/uuid"
	"github.com/rivo/tview"
)

const UpdateProviderPageName = "provider-update-page"

type updateProviderReq struct {
	name, details string
}

type UpdateProvider struct {
	*FlexApp
	*Navigation
	app *TerminalApp
}

func (u *UpdateProvider) NavigationInfo(app *TerminalApp, variables map[string]any) *NavigationInfo {
	return NewNavigationInfo(UpdateProviderPageName, func() tview.Primitive { return NewUpdateProvider(app, variables["id"].(uuid.UUID)) })
}

func (u *UpdateProvider) enrichNavigation(app *TerminalApp, providerId uuid.UUID) {
	u.Navigation = NewNavigation(app, u.NavigationInfo(app, map[string]any{"id": providerId}))
}

func NewUpdateProvider(app *TerminalApp, providerId uuid.UUID) *UpdateProvider {
	f := &UpdateProvider{
		FlexApp: NewFlexApp(),
		app:     app,
	}
	f.enrichNavigation(app, providerId)
	f.InitFlexApp(app)
	f.bindKeys()

	providerDto, err := app.GetProviderService().FindById(providerId)
	if err != nil {
		f.ShowInfoReturnBack(err.Error())
	}

	var request updateProviderReq

	form := tview.NewForm().
		AddInputField("Name", providerDto.Name, 20, nil, func(text string) { request.name = text }).
		AddInputField("Details", providerDto.Details, 20, nil, func(text string) { request.details = text }).
		AddButton("Update", f.update(providerId, request)).
		AddButton("Cancel", f.BackFunc())

	form.SetBorder(true).SetTitle("Update Provider").SetTitleAlign(tview.AlignCenter).SetRect(150, 30, 60, 15)

	f.AddItem(form, 0, 8, true)

	f.SetInputCapture(f.KeyboardFunc)

	return f
}

func (u *UpdateProvider) bindKeys() {
	u.Actions = KeyActions{
		tcell.KeyEscape: NewKeyAction("Back", u.KeyBack),
	}
}

func (u *UpdateProvider) update(id uuid.UUID, update updateProviderReq) func() {
	return func() {
		request := model.UpdateProviderRequest{
			Name:    update.name,
			Details: update.details,
		}

		if err := u.app.GetProviderService().Update(id, request); err != nil {
			u.ShowErrorTo(err)
		} else {
			u.ShowInfoReturnBack(fmt.Sprintf("Provider %s (%s) successfully updated.", request.Name, id))
		}
	}
}
