package accounts

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/adnaan/authn"

	glv "github.com/goliveview/controller"
)

type HandlerConfirmView struct {
	Auth *authn.API
}

func (h *HandlerConfirmView) EventHandler(ctx glv.Context) error {
	switch ctx.Event().ID {
	default:
		log.Printf("warning:handler not found for event => \n %+v\n", ctx.Event())
	}
	return nil
}

func (h *HandlerConfirmView) OnMount(w http.ResponseWriter, r *http.Request) (int, glv.M) {
	token := chi.URLParam(r, "token")
	err := h.Auth.ConfirmSignupEmail(r.Context(), token)
	if err != nil {
		log.Println("err confirm.onmount", err)
		return 200, nil
	}
	return 200, glv.M{
		"confirmed": true,
	}
}
