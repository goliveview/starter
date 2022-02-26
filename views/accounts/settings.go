package accounts

import (
	"log"
	"time"

	"github.com/adnaan/authn"

	glv "github.com/goliveview/controller"
)

type SettingsView struct {
	glv.DefaultView
	Auth *authn.API
}

func (s *SettingsView) Content() string {
	return "./templates/views/accounts/settings"
}

func (s *SettingsView) Layout() string {
	return "./templates/layouts/app.html"
}

func (s *SettingsView) OnEvent(ctx glv.Context) error {
	ctx.DOM().RemoveClass("#profile-loading", "is-hidden")
	ctx.DOM().SetAttributes("#profile_inputs", glv.M{"disabled": "disabled"})
	defer func() {
		time.Sleep(1 * time.Second)
		ctx.DOM().AddClass("#profile-loading", "is-hidden")
		ctx.DOM().RemoveAttributes("#profile_inputs", []string{"disabled"})
	}()
	switch ctx.Event().ID {
	case "account/update":
		return s.UpdateProfile(ctx)
	case "account/delete":
		return s.DeleteAccount(ctx)
	default:
		log.Printf("warning:handler not found for event => \n %+v\n", ctx.Event())
	}
	return nil
}

func (s *SettingsView) OnMount(ctx glv.Context) (glv.Status, glv.M) {
	r := ctx.Request()
	if r.Method != "GET" {
		return glv.Status{Code: 405}, nil
	}
	userID, _ := r.Context().Value(authn.AccountIDKey).(string)
	acc, err := s.Auth.GetAccount(r.Context(), userID)
	if err != nil {
		return glv.Status{Code: 200}, nil
	}

	name := ""
	m := acc.Attributes().Map()
	if m != nil {
		name, _ = m.String("name")
	}

	return glv.Status{Code: 200}, glv.M{
		"is_logged_in": true,
		"email":        acc.Email(),
		"name":         name,
	}
}

func (s *SettingsView) UpdateProfile(ctx glv.Context) error {
	req := new(ProfileRequest)
	if err := ctx.Event().DecodeParams(req); err != nil {
		return err
	}
	rCtx := ctx.Request().Context()
	userID, _ := rCtx.Value(authn.AccountIDKey).(string)
	acc, err := s.Auth.GetAccount(rCtx, userID)
	if err != nil {
		return err
	}
	if err := acc.Attributes().Set(rCtx, "name", req.Name); err != nil {
		return err
	}
	if req.Email != "" && req.Email != acc.Email() {
		if err := acc.ChangeEmail(rCtx, req.Email); err != nil {
			return err
		}
		ctx.DOM().RemoveClass("#change_email", "is-hidden")
	}

	ctx.DOM().Morph("#account_form", "account_form", glv.M{
		"name":  req.Name,
		"email": acc.Email(),
	})

	return nil
}

func (s *SettingsView) DeleteAccount(ctx glv.Context) error {
	rCtx := ctx.Request().Context()
	userID, _ := rCtx.Value(authn.AccountIDKey).(string)
	acc, err := s.Auth.GetAccount(rCtx, userID)
	if err != nil {
		return err
	}
	if err := acc.Delete(rCtx); err != nil {
		return err
	}
	ctx.DOM().Reload()
	return nil
}
