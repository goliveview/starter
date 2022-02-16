package accounts

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/adnaan/authn"

	glv "github.com/goliveview/controller"
)

type SignupView struct {
	glv.DefaultView
	Auth *authn.API
}

func (s *SignupView) Content() string {
	return "./templates/views/accounts/signup"
}

func (s *SignupView) Layout() string {
	return "./templates/layouts/index.html"
}

func (s *SignupView) OnEvent(ctx glv.Context) error {
	switch ctx.Event().ID {
	case "auth/signup":
		return s.Signup(ctx)
	default:
		log.Printf("warning:handler not found for event => \n %+v\n", ctx.Event())
	}
	return nil
}

func (s *SignupView) OnMount(w http.ResponseWriter, r *http.Request) (glv.Status, glv.M) {
	if _, err := s.Auth.CurrentAccount(r); err != nil {
		return glv.Status{Code: 200}, nil
	}

	return glv.Status{Code: 200}, glv.M{
		"is_logged_in": true,
	}
}

type ProfileRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *SignupView) Signup(ctx glv.Context) error {
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

	if err := s.Auth.Signup(ctx.RequestContext(), r.Email, r.Password, attributes); err != nil {
		return err
	}
	ctx.DOM().Morph("#signup_container", "signup_container", glv.M{
		"sent_confirmation": true,
	})
	return nil
}
