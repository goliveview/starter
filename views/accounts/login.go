package accounts

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/adnaan/authn"

	glv "github.com/goliveview/controller"
)

type LoginView struct {
	glv.DefaultView
	Auth *authn.API
}

func (l *LoginView) Content() string {
	return "./templates/views/accounts/login"
}

func (l *LoginView) Layout() string {
	return "./templates/layouts/index.html"
}

func (l *LoginView) OnEvent(ctx glv.Context) error {
	ctx.DOM().AddClass("#loading-modal", "is-active")
	defer func() {
		ctx.DOM().RemoveClass("#loading-modal", "is-active")
	}()
	switch ctx.Event().ID {
	case "auth/magic-login":
		return l.MagicLogin(ctx)
	default:
		log.Printf("warning:handler not found for event => \n %+v\n", ctx.Event())
	}
	return nil
}

func (l *LoginView) OnMount(w http.ResponseWriter, r *http.Request) (glv.Status, glv.M) {
	if r.Method == "POST" {
		return l.LoginSubmit(w, r)
	}

	if _, err := l.Auth.CurrentAccount(r); err != nil {
		return glv.Status{Code: 200}, nil
	}

	return glv.Status{Code: 200}, glv.M{
		"is_logged_in": true,
	}
}

func (l *LoginView) LoginSubmit(w http.ResponseWriter, r *http.Request) (glv.Status, glv.M) {
	var email, password string
	_ = r.ParseForm()
	for k, v := range r.Form {
		if k == "email" && len(v) == 0 {
			return glv.Status{Code: 200}, glv.M{
				"error": "email is required",
			}
		}

		if k == "password" && len(v) == 0 {
			return glv.Status{Code: 200}, glv.M{
				"error": "password is required",
			}
		}

		if len(v) == 0 {
			continue
		}

		if k == "email" && len(v) > 0 {
			email = v[0]
			continue
		}

		if k == "password" && len(v) > 0 {
			password = v[0]
			continue
		}
	}
	if err := l.Auth.Login(w, r, email, password); err != nil {
		return glv.Status{Code: 200}, glv.M{
			"error": glv.UserError(err),
		}
	}
	redirectTo := "/app"
	from := r.URL.Query().Get("from")
	if from != "" {
		redirectTo = from
	}

	http.Redirect(w, r, redirectTo, http.StatusSeeOther)

	return glv.Status{Code: 200}, glv.M{}
}

func (l *LoginView) MagicLogin(ctx glv.Context) error {
	r := new(ProfileRequest)
	if err := ctx.Event().DecodeParams(r); err != nil {
		return err
	}
	if r.Email == "" {
		return fmt.Errorf("%w", errors.New("email is required"))
	}
	if err := l.Auth.SendPasswordlessToken(ctx.RequestContext(), r.Email); err != nil {
		return err
	}
	ctx.DOM().Morph("#signin_container", "signin_container", glv.M{"sent_magic_link": true})
	return nil
}
