package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gu-io/midash/pkg/internals/db"
	"github.com/gu-io/midash/pkg/internals/models/user"
	"github.com/gu-io/midash/pkg/internals/utils"
	"github.com/influx6/faux/sink"
)

// Users exposes a central handle for which requests are served to all requests.
type Users struct {
	DB  db.DB
	Log sink.Sink
}

// CreateUser handles receiving requests to create a user from the server.
func (u Users) CreateUser(w http.ResponseWriter, r *http.Request, params map[string]string) {
	var nw user.NewUser

	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&nw); err != nil {
		http.Error(w, utils.ErrorMessage(http.StatusInternalServerError, "Failed to read body", err), http.StatusInternalServerError)
		return
	}

	newUser, err := user.New(nw)
	if err != nil {
		http.Error(w, utils.ErrorMessage(http.StatusInternalServerError, "Failed to create user", err), http.StatusInternalServerError)
		return
	}

	dbi, err := u.DB.New()
	if err != nil {
		http.Error(w, utils.ErrorMessage(http.StatusInternalServerError, "Failed to connect to database", err), http.StatusInternalServerError)
		return
	}

	if err := db.Save(dbi, newUser); err != nil {
		http.Error(w, utils.ErrorMessage(http.StatusInternalServerError, "Failed to save new user", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(newUser.SafeFields()); err != nil {
		http.Error(w, utils.ErrorMessage(http.StatusInternalServerError, "Failed to return new user data", err), http.StatusInternalServerError)
		return
	}
}
