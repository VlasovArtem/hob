package tui

import (
	"fmt"
	userModel "github.com/VlasovArtem/hob/src/user/model"
	"github.com/rivo/tview"
)

const SingUpPageName = "signup"

type SignUp struct {
	*tview.Form
	*NavigationBack
	request userModel.CreateUserRequest
}

func (s *SignUp) my(app *TerminalApp, parent *NavigationInfo) *NavigationInfo {
	return NewNavigationInfo(SingUpPageName, func() tview.Primitive { return NewSignUp(app, parent) })
}

func (s *SignUp) enrichNavigation(app *TerminalApp, parent *NavigationInfo) {
	s.MyNavigation = interface{}(s).(MyNavigation)
	s.enrich(app, parent)
}

func NewSignUp(app *TerminalApp, parent *NavigationInfo) *SignUp {
	f := &SignUp{
		Form:           tview.NewForm(),
		NavigationBack: NewNavigationBack(NewNavigation()),
	}
	f.enrichNavigation(app, parent)

	f.
		AddInputField("Email", "", 20, nil, func(text string) { f.request.Email = text }).
		AddInputField("Last Name", "", 20, nil, func(text string) { f.request.LastName = text }).
		AddInputField("First Name", "", 20, nil, func(text string) { f.request.FirstName = text }).
		AddPasswordField("Password", "", 20, '*', func(text string) { f.request.Password = text }).
		AddButton("Sign Up", func() {
			if userResponse, err := app.GetUserService().Add(f.request); err != nil {
				f.ShowErrorTo(err)
			} else {
				app.AuthorizedUser = userResponse
				f.ShowInfoReturnHome(fmt.Sprintf("Welcome, %s %s to 'House of Bills'!", userResponse.LastName, userResponse.FirstName))
			}
		}).
		AddButton("Cancel", func() {
			f.Back()
		})

	f.SetBorder(true).SetTitle("Sign In").SetTitleAlign(tview.AlignCenter).SetRect(150, 30, 60, 15)

	f.SetInputCapture(f.KeyboardFunc)

	return f
}
