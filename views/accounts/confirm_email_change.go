package accounts

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/adnaan/authn"

	glv "github.com/goliveview/controller"
)

type ConfirmEmailChangeView struct {
	glv.DefaultView
	Auth *authn.API
}

func (c *ConfirmEmailChangeView) Content() string {
	return "./templates/views/accounts/confirm_email_change"
}

func (c *ConfirmEmailChangeView) Layout() string {
	return "./templates/layouts/app.html"
}

func (c *ConfirmEmailChangeView) OnMount(ctx glv.Context) (glv.Status, glv.M) {
	r := ctx.Request()
	if r.Method != "GET" {
		return glv.Status{Code: 405}, nil
	}
	token := chi.URLParam(r, "token")
	userID, _ := r.Context().Value(authn.AccountIDKey).(string)
	acc, err := c.Auth.GetAccount(r.Context(), userID)
	if err != nil {
		log.Printf("confirm change email: GetAccount err %v", err)
		return glv.Status{Code: 200}, nil
	}

	if err := acc.ConfirmEmailChange(r.Context(), token); err != nil {
		log.Printf("confirm change email: ConfirmEmailChange err %v", err)
		return glv.Status{Code: 200}, nil
	}

	redirectTo := "/account"
	http.Redirect(ctx.ResponseWriter(), r, redirectTo, http.StatusSeeOther)
	return glv.Status{Code: 200}, nil
}
