package accounts

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/adnaan/authn"

	glv "github.com/goliveview/controller"
)

type HandlerConfirmEmailChangeView struct {
	Auth *authn.API
}

func (h *HandlerConfirmEmailChangeView) EventHandler(ctx glv.Context) error {
	switch ctx.Event().ID {
	default:
		log.Printf("warning:handler not found for event => \n %+v\n", ctx.Event())
	}
	return nil
}

func (h *HandlerConfirmEmailChangeView) OnMount(w http.ResponseWriter, r *http.Request) (int, glv.M) {
	if r.Method != "GET" {
		return 405, nil
	}
	token := chi.URLParam(r, "token")
	userID, _ := r.Context().Value(authn.AccountIDKey).(string)
	acc, err := h.Auth.GetAccount(r.Context(), userID)
	if err != nil {
		log.Printf("confirm change email: GetAccount err %v", err)
		return 200, nil
	}

	if err := acc.ConfirmEmailChange(r.Context(), token); err != nil {
		log.Printf("confirm change email: ConfirmEmailChange err %v", err)
		return 200, nil
	}

	redirectTo := "/account"
	http.Redirect(w, r, redirectTo, http.StatusSeeOther)
	return 200, nil
}
