package tui

import (
	"context"
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

func (s *SignInForm) my(app *TerminalApp, ctx context.Context) *NavigationInfo {
	return NewNavigationInfo(SignInPageName, func() tview.Primitive { return NewSignIn(app) })
}

func (s *SignInForm) enrichNavigation(app *TerminalApp, ctx context.Context) {
	s.MyNavigation = interface{}(s).(MyNavigation)
	s.enrich(app, ctx).addCustomPage(ctx, &SignUp{})
}

func NewSignIn(app *TerminalApp) *SignInForm {

	f := &SignInForm{
		Form:       tview.NewForm(),
		Navigation: NewNavigation(),
		app:        app,
	}
	f.MyNavigation = interface{}(f).(MyNavigation)
	f.enrichNavigation(app, nil)

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
