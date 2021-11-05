package accounts

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/adnaan/authn"

	glv "github.com/goliveview/controller"
)

type HandlerResetView struct {
	Auth *authn.API
}

func (h *HandlerResetView) EventHandler(ctx glv.Context) error {
	switch ctx.Event().ID {
	case "account/reset":
		return h.Reset(ctx)
	default:
		log.Printf("warning:handler not found for event => \n %+v\n", ctx.Event())
	}
	return nil
}

func (h *HandlerResetView) OnMount(w http.ResponseWriter, r *http.Request) (int, glv.M) {
	token := chi.URLParam(r, "token")
	return 200, glv.M{
		"token": token,
	}
}

type ResetReq struct {
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	Token           string `json:"token"`
}

func (h *HandlerResetView) Reset(ctx glv.Context) error {
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
	if err := h.Auth.ConfirmRecovery(ctx.RequestContext(), r.Token, r.Password); err != nil {
		return err
	}
	ctx.DOM().RemoveClass("#password_reset", "is-hidden")
	return nil
}
