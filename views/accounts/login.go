package accounts

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/adnaan/authn"

	glv "github.com/goliveview/controller"
)

type HandlerLoginView struct {
	Auth *authn.API
}

func (h *HandlerLoginView) EventHandler(ctx glv.Context) error {
	ctx.DOM().AddClass("#loading-modal", "is-active")
	defer func() {
		ctx.DOM().RemoveClass("#loading-modal", "is-active")
	}()
	switch ctx.Event().ID {
	case "auth/magic-login":
		return h.MagicLogin(ctx)
	default:
		log.Printf("warning:handler not found for event => \n %+v\n", ctx.Event())
	}
	return nil
}

func (h *HandlerLoginView) OnMount(w http.ResponseWriter, r *http.Request) (int, glv.M) {
	if r.Method == "POST" {
		return h.LoginSubmit(w, r)
	}
	return 200, glv.M{}
}

func (h *HandlerLoginView) LoginSubmit(w http.ResponseWriter, r *http.Request) (int, glv.M) {
	var email, password string
	_ = r.ParseForm()
	for k, v := range r.Form {
		if k == "email" && len(v) == 0 {
			return 200, glv.M{
				"error": "email is required",
			}
		}

		if k == "password" && len(v) == 0 {
			return 200, glv.M{
				"error": "password is required",
			}
		}

		if len(v) == 0 {
			continue
		}

		if k == "email" && len(v) > 0 {
			email = v[0]
			continue
		}

		if k == "password" && len(v) > 0 {
			password = v[0]
			continue
		}
	}
	if err := h.Auth.Login(w, r, email, password); err != nil {
		return 200, glv.M{
			"error": err.Error(),
		}
	}
	redirectTo := "/app"
	from := r.URL.Query().Get("from")
	if from != "" {
		redirectTo = from
	}

	http.Redirect(w, r, redirectTo, http.StatusSeeOther)

	return 200, glv.M{}
}

func (h *HandlerLoginView) MagicLogin(ctx glv.Context) error {
	r := new(AuthRequest)
	if err := ctx.Event().DecodeParams(r); err != nil {
		return err
	}
	if r.Email == "" {
		return fmt.Errorf("%w", errors.New("email is required"))
	}
	if err := h.Auth.SendPasswordlessToken(ctx.RequestContext(), r.Email); err != nil {
		return err
	}
	ctx.DOM().Morph("#signin_container", "signin_container", glv.M{"sent_magic_link": true})
	return nil
}
