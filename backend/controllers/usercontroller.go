package controller

import (
	"encoding/json"
	"net/http"

	"github.com/adityanath-madishetti/todo/backend/middleware"
	model "github.com/adityanath-madishetti/todo/backend/models"
	"github.com/adityanath-madishetti/todo/backend/utils"
	"golang.org/x/crypto/bcrypt"
)

func ChangePassword(w http.ResponseWriter, r *http.Request) {
	userid, ok := r.Context().Value(middleware.ContextKeyUserID).(string)
	if !ok {
		utils.SendJSONError(w, http.StatusInternalServerError, "Something wrong with token, try logging in again")
		return
	}

	type passwordStatus struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	var reqBody passwordStatus
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		utils.SendJSONError(w, http.StatusBadRequest, "Invalid JSON object sent: "+err.Error())
		return
	}

	user, err := model.GetUserFromUserId(userid)
	if err != nil {
		utils.SendJSONError(w, http.StatusInternalServerError, "Error fetching user: "+err.Error())
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqBody.OldPassword)); err != nil {
		utils.SendJSONError(w, http.StatusUnauthorized, "Invalid old password: "+err.Error())
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reqBody.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		utils.SendJSONError(w, http.StatusInternalServerError, "Failed to hash new password")
		return
	}

	if err := model.UpdatePassword(userid, string(hashedPassword)); err != nil {
		utils.SendJSONError(w, http.StatusInternalServerError, "Failed to update password: "+err.Error())
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Password changed successfully"})
}
