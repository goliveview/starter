package accounts

import (
	"errors"
	"fmt"
	"log"

	"github.com/go-chi/chi"

	"github.com/adnaan/authn"

	glv "github.com/goliveview/controller"
)

type ResetView struct {
	glv.DefaultView
	Auth *authn.API
}

func (rv *ResetView) Content() string {
	return "./templates/views/accounts/reset"
}

func (rv *ResetView) Layout() string {
	return "./templates/layouts/index.html"
}

func (rv *ResetView) OnEvent(ctx glv.Context) error {
	switch ctx.Event().ID {
	case "account/reset":
		return rv.Reset(ctx)
	default:
		log.Printf("warning:handler not found for event => \n %+v\n", ctx.Event())
	}
	return nil
}

func (rv *ResetView) OnMount(ctx glv.Context) (glv.Status, glv.M) {
	token := chi.URLParam(ctx.Request(), "token")
	return glv.Status{Code: 200}, glv.M{
		"token": token,
	}
}

type ResetReq struct {
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	Token           string `json:"token"`
}

func (rv *ResetView) Reset(ctx glv.Context) error {
	ctx.DOM().AddClass("#loading-modal", "is-active")
	defer func() {
		ctx.DOM().RemoveClass("#loading-modal", "is-active")
	}()
	r := new(ResetReq)
	if err := ctx.Event().DecodeParams(r); err != nil {
		return err
	}
	if r.ConfirmPassword != r.Password {
		return fmt.Errorf("%w", errors.New("passwords don't match"))
	}
	if err := rv.Auth.ConfirmRecovery(ctx.Request().Context(), r.Token, r.Password); err != nil {
		return err
	}
	ctx.DOM().AddClass("#new_password", "is-hidden")
	ctx.DOM().RemoveClass("#password_reset", "is-hidden")
	return nil
}
