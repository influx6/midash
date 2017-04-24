package handlers

import (
	"errors"

	"github.com/gu-io/midash/pkg/db"
	"github.com/gu-io/midash/pkg/internals/models/user"
	"github.com/influx6/faux/sink"
	"github.com/influx6/faux/sink/sinks"
)

// Users exposes a central handle for which requests are served to all requests.
type Users struct {
	DB  db.DB
	Log sink.Sink
}

// Delete handles receiving requests to delete a user from the database.
func (u Users) Delete(id string) error {
	defer u.Log.Emit(sinks.Info("Get Existing User").With("user_id", id).Trace("handlers.Users.Create").End())

	dbi, err := u.DB.New()
	if err != nil {
		u.Log.Emit(sinks.Error(err).WithFields(sink.Fields{"public_id": id}))
		return err
	}

	defer dbi.Close()

	var nu user.User
	if err := db.Delete(u.Log, dbi, nu, "public_id", id); err != nil {
		u.Log.Emit(sinks.Error(err).WithFields(sink.Fields{"public_id": id}))
		return err
	}

	// Add user profile.
	profiles := Profiles{Log: u.Log, DB: u.DB}
	if err = profiles.DeleteByUser(id); err != nil {
		u.Log.Emit(sinks.Error(err).WithFields(sink.Fields{"public_id": id}))
		return err
	}

	return nil
}

// Get handles receiving requests to retrieve a user from the database.
func (u Users) Get(id string) (*user.User, error) {
	defer u.Log.Emit(sinks.Info("Get Existing User").With("user_id", id).Trace("handlers.Users.Create").End())

	dbi, err := u.DB.New()
	if err != nil {
		u.Log.Emit(sinks.Error(err).WithFields(sink.Fields{"public_id": id}))
		return nil, err
	}

	defer dbi.Close()

	var nu user.User

	if err := db.Get(u.Log, dbi, nu, &nu, "public_id", id); err != nil {
		u.Log.Emit(sinks.Error(err).WithFields(sink.Fields{"public_id": id}))
		return nil, err
	}

	// Add user profile.
	profiles := Profiles{Log: u.Log, DB: u.DB}
	nu.Profile, err = profiles.GetByUser(nu.PublicID)
	if err != nil {
		u.Log.Emit(sinks.Error(err).WithFields(sink.Fields{"public_id": id}))
		return nil, err
	}

	return &nu, nil
}

// GetByEmail handles receiving requests to retrieve a user with user's email from the database.
func (u Users) GetByEmail(email string) (*user.User, error) {
	defer u.Log.Emit(sinks.Info("Get Existing User").With("user_email", email).Trace("handlers.Users.Create").End())

	dbi, err := u.DB.New()
	if err != nil {
		u.Log.Emit(sinks.Error(err).WithFields(sink.Fields{"user_email": email}))
		return nil, err
	}

	defer dbi.Close()

	var nu user.User

	if err := db.Get(u.Log, dbi, nu, &nu, "email", email); err != nil {
		u.Log.Emit(sinks.Error(err).WithFields(sink.Fields{"user_email": email}))
		return nil, err
	}

	return &nu, nil
}

// UserRecords defines a struct which returns the total fields and page details
// used in retrieving the records.
type UserRecords struct {
	Total         int         `json:"total"`
	Page          int         `json:"page"`
	ResponserPage int         `json:"responserPerPage"`
	Records       []user.User `json:"records"`
}

// GetAll handles receiving requests to retrieve all user from the database.
func (u Users) GetAll(page, responsePerPage int) (UserRecords, error) {
	defer u.Log.Emit(sinks.Info("Get Existing User").WithFields(sink.Fields{
		"page":            page,
		"responsePerPage": responsePerPage,
	}).Trace("handlers.Users.Create").End())

	dbi, err := u.DB.New()
	if err != nil {
		u.Log.Emit(sinks.Error(err).WithFields(sink.Fields{
			"page":            page,
			"responsePerPage": responsePerPage,
		}))

		return UserRecords{}, err
	}

	defer dbi.Close()

	var nu user.User
	records, realTotalRecords, err := db.GetAllPerPage(u.Log, dbi, nu, "asc", "public_id", page, responsePerPage)
	if err != nil {
		u.Log.Emit(sinks.Error(err).WithFields(sink.Fields{
			"page":            page,
			"responsePerPage": responsePerPage,
		}))

		return UserRecords{}, err
	}

	var userRecords []user.User

	for _, record := range records {
		var nw user.User

		if err := nw.WithFields(record); err != nil {
			u.Log.Emit(sinks.Error(err).WithFields(sink.Fields{
				"page":            page,
				"responsePerPage": responsePerPage,
			}))
			return UserRecords{}, err
		}

		userRecords = append(userRecords, nw)
	}

	return UserRecords{
		Page:          page,
		Total:         realTotalRecords,
		ResponserPage: responsePerPage,
		Records:       userRecords,
	}, nil
}

// Create handles receiving requests to create a user from the server.
func (u Users) Create(nw user.NewUser) (*user.User, error) {
	defer u.Log.Emit(sinks.Info("Create New User").Trace("handlers.Users.Create").End())

	newUser, err := user.New(nw)
	if err != nil {
		u.Log.Emit(sinks.Error(err).WithFields(sink.Fields{"email": nw.Email}))
		return nil, err
	}

	dbi, err := u.DB.New()
	if err != nil {
		u.Log.Emit(sinks.Error(err).WithFields(sink.Fields{"email": nw.Email}))
		return nil, err
	}

	defer dbi.Close()

	if err := db.Save(u.Log, dbi, newUser); err != nil {
		u.Log.Emit(sinks.Error(err).WithFields(sink.Fields{"email": nw.Email}))
		return nil, err
	}

	dbi.Close()

	// Add user profile.
	profiles := Profiles{Log: u.Log, DB: u.DB}
	newUser.Profile, err = profiles.Create(newUser, nil)
	if err != nil {
		u.Log.Emit(sinks.Error(err).WithFields(sink.Fields{"email": nw.Email}))
		return nil, err
	}

	return newUser, nil
}

// UpdatePassword handles receiving requests to update a user identified by it's public_id.
func (u Users) UpdatePassword(nw user.UpdateUserPassword) error {
	defer u.Log.Emit(sinks.Info("Update User Password").With("user", nw.PublicID).Trace("handlers.Users.UpdatePassword").End())

	if nw.PublicID == "" {
		err := errors.New("JSON UpdateUserPassword.PublicID is empty")

		u.Log.Emit(sinks.Error(err).WithFields(sink.Fields{
			"user_id": nw.PublicID,
		}))

		return err
	}

	// TODO(influx6): Should we do some password validty checks.
	if nw.Password == "" {
		err := errors.New("JSON UpdateUserPassword.Password is empty")

		u.Log.Emit(sinks.Error(err).WithFields(sink.Fields{
			"user_id": nw.PublicID,
		}))

		return err
	}

	// TODO(influx6): Do we need to do this here.
	// if nw.Password != nw.PasswordConfirm {
	// 	err := errors.New("Invalid Confirmation Password")
	// 	u.Log.Emit(sinks.Error(err).WithFields(sink.Fields{
	//		"user_id":   nw.PublicID,
	// 	}))
	// 	http.Error(w, utils.ErrorMessage(http.StatusInternalServerError, "Failed to connect to database", err), http.StatusInternalServerError)
	// 	return
	// }

	dbi, err := u.DB.New()
	if err != nil {
		u.Log.Emit(sinks.Error(err).WithFields(sink.Fields{
			"user_id": nw.PublicID,
		}))

		return err
	}
	defer dbi.Close()

	var dbUser user.User

	if err := db.Get(u.Log, dbi, dbUser, &dbUser, "public_id", nw.PublicID); err != nil {
		u.Log.Emit(sinks.Error(err).WithFields(sink.Fields{
			"user_id": nw.PublicID,
		}))

		return err
	}

	if err := dbUser.ChangePassword(nw.Password); err != nil {
		u.Log.Emit(sinks.Error(err).WithFields(sink.Fields{
			"user_id": nw.PublicID,
		}))

		return err
	}

	if err := db.Update(u.Log, dbi, &dbUser, "public_id"); err != nil {
		u.Log.Emit(sinks.Error(err).WithFields(sink.Fields{
			"user_id": nw.PublicID,
		}))

		return err
	}

	return nil
}

// Update handles receiving requests to update a user identified by it's public_id.
func (u Users) Update(nw user.UpdateUser) error {
	defer u.Log.Emit(sinks.Info("Update User").With("user", nw.PublicID).Trace("handlers.Users.Update").End())

	if nw.PublicID == "" {
		err := errors.New("JSON User.PublicID is empty")
		u.Log.Emit(sinks.Error(err).WithFields(sink.Fields{
			"user_id": nw.PublicID,
			"email":   nw.Email,
		}))

		return err
	}

	dbi, err := u.DB.New()
	if err != nil {
		u.Log.Emit(sinks.Error(err).WithFields(sink.Fields{
			"user_id": nw.PublicID,
			"email":   nw.Email,
		}))

		return err
	}
	defer dbi.Close()

	if err := db.Update(u.Log, dbi, nw, "public_id"); err != nil {
		u.Log.Emit(sinks.Error(err).WithFields(sink.Fields{
			"user_id": nw.PublicID,
			"email":   nw.Email,
		}))

		return err
	}

	return nil
}
