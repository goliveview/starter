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
	r := new(ProfileRequest)
	if err := ctx.Event().DecodeParams(r); err != nil {
		return err
	}

	if err := f.Auth.Recovery(ctx.RequestContext(), r.Email); err != nil {
		return err
	}

	ctx.DOM().RemoveClass("#recovery_sent", "is-hidden")
	return nil
}
