package accounts

import (
	"log"

	"github.com/go-chi/chi"

	"github.com/adnaan/authn"

	glv "github.com/goliveview/controller"
)

type ConfirmView struct {
	glv.DefaultView
	Auth *authn.API
}

func (h *ConfirmView) Content() string {
	return "./templates/views/accounts/confirm"
}

func (h *ConfirmView) Layout() string {
	return "./templates/layouts/index.html"
}

func (h *ConfirmView) OnMount(ctx glv.Context) (glv.Status, glv.M) {
	r := ctx.Request()
	token := chi.URLParam(r, "token")
	err := h.Auth.ConfirmSignupEmail(r.Context(), token)
	if err != nil {
		log.Println("err confirm.onmount", err)
		return glv.Status{Code: 200}, nil
	}
	return glv.Status{Code: 200}, glv.M{
		"confirmed": true,
	}
}
