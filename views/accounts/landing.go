package accounts

import (
	"log"
	"net/http"

	"github.com/adnaan/authn"

	glv "github.com/goliveview/controller"
)

type HandlerLandingView struct {
	Auth *authn.API
}

func (h *HandlerLandingView) EventHandler(ctx glv.Context) error {
	switch ctx.Event().ID {
	default:
		log.Printf("warning:handler not found for event => \n %+v\n", ctx.Event())
	}
	return nil
}

func (h *HandlerLandingView) OnMount(w http.ResponseWriter, r *http.Request) (int, glv.M) {
	if r.Method != "GET" {
		return 405, nil
	}
	if _, err := h.Auth.CurrentAccount(r); err != nil {
		return 200, nil
	}

	return 200, glv.M{
		"is_logged_in": true,
	}
}
