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
	return 200, glv.M{}
}
