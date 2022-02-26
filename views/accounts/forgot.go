package accounts

import (
	"log"

	"github.com/adnaan/authn"

	glv "github.com/goliveview/controller"
)

type ForgotView struct {
	glv.DefaultView
	Auth *authn.API
}

func (f *ForgotView) Content() string {
	return "./templates/views/accounts/forgot"
}

func (f *ForgotView) Layout() string {
	return "./templates/layouts/index.html"
}

func (f *ForgotView) OnEvent(ctx glv.Context) error {
	switch ctx.Event().ID {
	case "account/forgot":
		return f.SendRecovery(ctx)
	default:
		log.Printf("warning:handler not found for event => \n %+v\n", ctx.Event())
	}
	return nil
}

func (f *ForgotView) SendRecovery(ctx glv.Context) error {
	ctx.DOM().AddClass("#loading-modal", "is-active")
	defer func() {
		ctx.DOM().RemoveClass("#loading-modal", "is-active")
	}()
	req := new(ProfileRequest)
	if err := ctx.Event().DecodeParams(req); err != nil {
		return err
	}

	if err := f.Auth.Recovery(ctx.Request().Context(), req.Email); err != nil {
		return err
	}

	ctx.DOM().AddClass("#input_email", "is-hidden")
	ctx.DOM().RemoveClass("#recovery_sent", "is-hidden")
	return nil
}
