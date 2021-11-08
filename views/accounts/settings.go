package accounts

import (
	"log"
	"net/http"
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

func (s *SettingsView) OnMount(w http.ResponseWriter, r *http.Request) (int, glv.M) {
	if r.Method != "GET" {
		return 405, nil
	}
	userID, _ := r.Context().Value(authn.AccountIDKey).(string)
	acc, err := s.Auth.GetAccount(r.Context(), userID)
	if err != nil {
		return 200, nil
	}

	name := ""
	m := acc.Attributes().Map()
	if m != nil {
		name, _ = m.String("name")
	}

	return 200, glv.M{
		"is_logged_in": true,
		"email":        acc.Email(),
		"name":         name,
	}
}

func (s *SettingsView) UpdateProfile(ctx glv.Context) error {

	ctx.DOM().RemoveClass("#profile-loading", "is-hidden")
	defer func() {
		time.Sleep(1 * time.Second)
		ctx.DOM().AddClass("#profile-loading", "is-hidden")
	}()
	r := new(ProfileRequest)
	if err := ctx.Event().DecodeParams(r); err != nil {
		return err
	}
	userID, _ := ctx.RequestContext().Value(authn.AccountIDKey).(string)
	acc, err := s.Auth.GetAccount(ctx.RequestContext(), userID)
	if err != nil {
		return err
	}

	if err := acc.Attributes().Set(ctx.RequestContext(), "name", r.Name); err != nil {
		return err
	}

	if r.Email != "" && r.Email != acc.Email() {
		if err := acc.ChangeEmail(ctx.RequestContext(), r.Email); err != nil {
			return err
		}
		ctx.DOM().RemoveClass("#change_email", "is-hidden")
	}

	return nil
}

func (s *SettingsView) DeleteAccount(ctx glv.Context) error {
	userID, _ := ctx.RequestContext().Value(authn.AccountIDKey).(string)
	acc, err := s.Auth.GetAccount(ctx.RequestContext(), userID)
	if err != nil {
		return err
	}
	if err := acc.Delete(ctx.RequestContext()); err != nil {
		return err
	}
	ctx.DOM().Reload()
	return nil
}
