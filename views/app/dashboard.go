package app

import (
	"log"
	"net/http"

	"github.com/adnaan/authn"
	glv "github.com/goliveview/controller"
)

type DashboardView struct {
	glv.DefaultView
	Auth *authn.API
}

func (d *DashboardView) Content() string {
	return "./templates/views/app"
}

func (d *DashboardView) Layout() string {
	return "./templates/layouts/app.html"
}

func (d *DashboardView) OnEvent(ctx glv.Context) error {
	switch ctx.Event().ID {
	default:
		log.Printf("warning:handler not found for event => \n %+v\n", ctx.Event())
	}
	return nil
}

func (d *DashboardView) OnMount(w http.ResponseWriter, r *http.Request) (int, glv.M) {
	return 200, glv.M{
		"is_logged_in": true,
	}
}
