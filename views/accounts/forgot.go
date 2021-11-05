package accounts

import (
	"log"
	"net/http"

	"github.com/adnaan/authn"

	glv "github.com/goliveview/controller"
)

type HandlerForgotView struct {
	Auth *authn.API
}

func (h *HandlerForgotView) EventHandler(ctx glv.Context) error {
	switch ctx.Event().ID {
	case "account/forgot":
		return h.SendRecovery(ctx)
	default:
		log.Printf("warning:handler not found for event => \n %+v\n", ctx.Event())
	}
	return nil
}

func (h *HandlerForgotView) OnMount(w http.ResponseWriter, r *http.Request) (int, glv.M) {
	return 200, nil
}

func (h *HandlerForgotView) SendRecovery(ctx glv.Context) error {
	ctx.DOM().AddClass("#loading-modal", "is-active")
	defer func() {
		ctx.DOM().RemoveClass("#loading-modal", "is-active")
	}()
	r := new(ProfileRequest)
	if err := ctx.Event().DecodeParams(r); err != nil {
		return err
	}

	if err := h.Auth.Recovery(ctx.RequestContext(), r.Email); err != nil {
		return err
	}

	ctx.DOM().RemoveClass("#recovery_sent", "is-hidden")
	return nil
}
