package app

import (
	"log"
	"net/http"

	"github.com/adnaan/authn"
	glv "github.com/goliveview/controller"
)

type HandlerDashboardView struct {
	Auth *authn.API
}

func (h *HandlerDashboardView) EventHandler(ctx glv.Context) error {
	switch ctx.Event().ID {
	default:
		log.Printf("warning:handler not found for event => \n %+v\n", ctx.Event())
	}
	return nil
}

func (h *HandlerDashboardView) OnMount(w http.ResponseWriter, r *http.Request) (int, glv.M) {
	return 200, glv.M{
		"is_logged_in": true,
	}
}
