package tui

import (
	"fmt"
	"github.com/rivo/tview"
)

const SignInPageName = "signin"

type SignInForm struct {
	*tview.Form
	*Navigation
	app      *TerminalApp
	email    string
	password string
}

func (s *SignInForm) NavigationInfo(app *TerminalApp, variables map[string]interface{}) *NavigationInfo {
	return NewNavigationInfo(SignInPageName, func() tview.Primitive { return NewSignIn(app) })
}

func (s *SignInForm) enrichNavigation(app *TerminalApp) {
	s.Navigation = NewNavigation(app, s.NavigationInfo(app, nil))
	s.AddCustomPage(&SignUp{})
}

func NewSignIn(app *TerminalApp) *SignInForm {

	f := &SignInForm{
		Form: tview.NewForm(),
		app:  app,
	}
	f.enrichNavigation(app)

	f.
		AddInputField("Email", "", 20, nil, func(text string) { f.email = text }).
		AddPasswordField("Password", "", 20, '*', func(text string) { f.password = text }).
		AddButton("Enter", func() {
			user, err := app.GetUserService().VerifyUser(f.email, f.password)
			if err != nil {
				f.ShowErrorTo(err)
			} else {
				app.AuthorizedUser = &user
				f.ShowInfoReturnHome(fmt.Sprintf("Welcome, %s %s to 'House of Bills'!", user.LastName, user.FirstName))
			}
		}).
		AddButton("Quit", func() {
			app.Stop()
		}).
		AddButton("Sign Up", func() {
			f.NavigateTo(SingUpPageName)
		})

	f.SetBorder(true).SetTitle("Sign In").SetTitleAlign(tview.AlignCenter).SetRect(150, 30, 60, 15)

	return f
}
