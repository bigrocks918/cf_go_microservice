package main

import (
	"fmt"
	"net/http"
)

func (app *config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := readJSON(w, r, &requestPayload)
	if err != nil {
		_ = errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// TODO validate against database
	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", requestPayload.Email),
		Data: User{
			ID:        1,
			FirstName: "Jack",
			LastName:  "Smith",
			Email:     "jack@smith.com",
			Active:    1,
		},
	}

	_ = writeJSON(w, http.StatusAccepted, payload)
}
