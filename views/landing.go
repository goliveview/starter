package views

import (
	"github.com/adnaan/authn"

	glv "github.com/goliveview/controller"
)

type LandingView struct {
	glv.DefaultView
	Auth *authn.API
}

func (l *LandingView) Content() string {
	return "./templates/views/landing"
}

func (l *LandingView) Layout() string {
	return "./templates/layouts/index.html"
}

func (l *LandingView) OnMount(ctx glv.Context) (glv.Status, glv.M) {
	r := ctx.Request()
	if r.Method != "GET" {
		return glv.Status{Code: 405}, nil
	}
	if _, err := l.Auth.CurrentAccount(r); err != nil {
		return glv.Status{Code: 200}, nil
	}
	return glv.Status{Code: 200}, glv.M{
		"is_logged_in": true,
	}
}
