package views

import (
	"errors"
	"fmt"
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
	case "auth/signup":
		return h.Signup(ctx)
	default:
		log.Printf("warning:handler not found for event => \n %+v\n", ctx.Event())
	}
	return nil
}

func (h *HandlerSignupView) OnMount(r *http.Request) (int, glv.M) {
	return 200, glv.M{
		"is_logged_in":      false,
		"sent_confirmation": false,
	}
}

type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *HandlerSignupView) Signup(ctx glv.Context) error {
	r := new(SignupRequest)
	if err := ctx.Event().DecodeParams(r); err != nil {
		return err
	}

	if r.Email == "" {
		return fmt.Errorf("%w", errors.New("email is required"))
	}
	if r.Password == "" {
		return fmt.Errorf("%w", errors.New("password is required"))
	}
	//
	//if err := h.Auth.Signup(ctx.RequestContext(), r.Email, r.Password, nil); err != nil {
	//	return err
	//}
	ctx.DOM().Morph("#signup_container", "signup_container", glv.M{
		"is_logged_in":      false,
		"sent_confirmation": true,
	})

	return nil
}
