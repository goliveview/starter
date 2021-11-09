package views

import (
	"net/http"

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

func (l *LandingView) OnMount(w http.ResponseWriter, r *http.Request) (glv.Status, glv.M) {
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
