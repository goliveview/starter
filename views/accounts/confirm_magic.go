package accounts

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/adnaan/authn"

	glv "github.com/goliveview/controller"
)

type HandlerConfirmMagicView struct {
	Auth *authn.API
}

func (h *HandlerConfirmMagicView) EventHandler(ctx glv.Context) error {
	switch ctx.Event().ID {
	default:
		log.Printf("warning:handler not found for event => \n %+v\n", ctx.Event())
	}
	return nil
}

func (h *HandlerConfirmMagicView) OnMount(w http.ResponseWriter, r *http.Request) (int, glv.M) {
	if r.Method != "GET" {
		return 405, nil
	}
	token := chi.URLParam(r, "token")
	err := h.Auth.LoginWithPasswordlessToken(w, r, token)
	if err != nil {
		return 200, nil
	}
	redirectTo := "/app"
	http.Redirect(w, r, redirectTo, http.StatusSeeOther)
	return 200, nil
}
