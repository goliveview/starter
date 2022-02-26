package accounts

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/adnaan/authn"

	glv "github.com/goliveview/controller"
)

type ConfirmMagicView struct {
	glv.DefaultView
	Auth *authn.API
}

func (c *ConfirmMagicView) Content() string {
	return "./templates/views/accounts/confirm_magic"
}

func (c *ConfirmMagicView) Layout() string {
	return "./templates/layouts/index.html"
}

func (c *ConfirmMagicView) OnMount(ctx glv.Context) (glv.Status, glv.M) {
	w, r := ctx.ResponseWriter(), ctx.Request()
	if r.Method != "GET" {
		return glv.Status{Code: 405}, nil
	}
	token := chi.URLParam(r, "token")
	err := c.Auth.LoginWithPasswordlessToken(w, r, token)
	if err != nil {
		return glv.Status{Code: 200}, nil
	}
	redirectTo := "/app"
	http.Redirect(w, r, redirectTo, http.StatusSeeOther)
	return glv.Status{Code: 200}, nil
}
