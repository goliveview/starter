package views

import (
	"log"
	"net/http"

	"github.com/adnaan/authn"

	glv "github.com/goliveview/controller"
)

type HandlerLoginView struct {
	Auth *authn.API
}

func (h *HandlerLoginView) EventHandler(ctx glv.Context) error {
	switch ctx.Event().ID {
	default:
		log.Printf("warning:handler not found for event => \n %+v\n", ctx.Event())
	}
	return nil
}

func (h *HandlerLoginView) OnMount(r *http.Request) (int, glv.M) {

	return 200, glv.M{
		"is_logged_in":      false,
		"confirmed":         false,
		"not_confirmed":     false,
		"confirmation_sent": false,
		"userErrs":          nil,
	}
}
