package accounts

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

func (h *HandlerSignupView) OnMount(w http.ResponseWriter, r *http.Request) (int, glv.M) {
	if _, err := h.Auth.CurrentAccount(r); err != nil {
		return 200, glv.M{
			"is_logged_in": false,
		}
	}

	return 200, glv.M{
		"is_logged_in": true,
	}
}

type ProfileRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *HandlerSignupView) Signup(ctx glv.Context) error {
	ctx.DOM().AddClass("#loading-modal", "is-active")
	defer func() {
		ctx.DOM().RemoveClass("#loading-modal", "is-active")
	}()
	r := new(ProfileRequest)
	if err := ctx.Event().DecodeParams(r); err != nil {
		return err
	}

	if r.Email == "" {
		return fmt.Errorf("%w", errors.New("email is required"))
	}
	if r.Password == "" {
		return fmt.Errorf("%w", errors.New("password is required"))
	}

	attributes := make(map[string]interface{})
	attributes["name"] = r.Name

	if err := h.Auth.Signup(ctx.RequestContext(), r.Email, r.Password, attributes); err != nil {
		return err
	}
	ctx.DOM().Morph("#signup_container", "signup_container", glv.M{
		"sent_confirmation": true,
	})
	return nil
}
