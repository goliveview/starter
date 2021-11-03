package views

import (
	"log"
	"net/http"

	"github.com/adnaan/authn"

	glv "github.com/goliveview/controller"
)

type HandlerSignupView struct {
	Auth *authn.API
}

func (h *HandlerSignupView) EventHandler(ctx glv.Context) error {
	switch ctx.Event().ID {
	default:
		log.Printf("warning:handler not found for event => \n %+v\n", ctx.Event())
	}
	return nil
}

func (h *HandlerSignupView) OnMount(r *http.Request) (int, glv.M) {

	return 200, glv.M{
		"is_logged_in": false,
	}
}
