package handlers

import (
	"errors"

	"github.com/gu-io/midash/pkg/db"
	"github.com/gu-io/midash/pkg/internals/models/profile"
	"github.com/gu-io/midash/pkg/internals/models/user"
	"github.com/influx6/faux/sink"
	"github.com/influx6/faux/sink/sinks"
)

// Profiles defines a handler which provides profile related methods.
type Profiles struct {
	DB  db.DB
	Log sink.Sink
}

// Create adds a new profile for the specified profile.
func (p Profiles) Create(nu *user.User, np *profile.NewProfile) (*profile.Profile, error) {
	defer p.Log.Emit(sinks.Info("Create New Profile").WithFields(sink.Fields{
		"user_email": nu.Email,
		"user_id":    nu.PublicID,
	}).Trace("Profiles.Create").End())

	if np != nil && np.UserID != nu.PublicID {
		err := errors.New("Invalid NewProfile.UserID: Does not match given user")
		p.Log.Emit(sinks.Error(err).WithFields(sink.Fields{"user_email": nu.Email, "user_id": nu.PublicID}))
		return nil, err
	}

	dbi, err := p.DB.New()
	if err != nil {
		p.Log.Emit(sinks.Error(err).WithFields(sink.Fields{"user_email": nu.Email, "user_id": nu.PublicID}))
		return nil, err
	}

	var newProfile profile.Profile

	// Attempt to retrieve profile from db if we still have an outstanding non-expired profile.
	if err := db.Get(p.Log, dbi, newProfile, &newProfile, profile.UniqueIndex, nu.PublicID); err == nil {
		return &newProfile, nil
	}

	newProfile = *profile.New(nu.PublicID)

	// If the provide data provided is not null, then add changes.
	if np != nil {
		newProfile.Address = np.Address
		newProfile.FirstName = np.FirstName
		newProfile.LastName = np.LastName
	}

	if err := db.Save(p.Log, dbi, &newProfile); err != nil {
		p.Log.Emit(sinks.Error(err).WithFields(sink.Fields{"user_email": nu.Email, "user_id": nu.PublicID}))
		return nil, err
	}

	return &newProfile, nil
}

// ProfileRecords defines a struct which returns the total fields and page details
// used in retrieving the records.
type ProfileRecords struct {
	Total         int               `json:"total"`
	Page          int               `json:"page"`
	ResponserPage int               `json:"responserPerPage"`
	Records       []profile.Profile `json:"records"`
}

// GetAll handles receiving requests to retrieve all profile from the database.
func (p Profiles) GetAll(page, responsePerPage int) (ProfileRecords, error) {
	defer p.Log.Emit(sinks.Info("Get Existing User").WithFields(sink.Fields{
		"page":            page,
		"responsePerPage": responsePerPage,
	}).Trace("handlers.Users.Create").End())

	dbi, err := p.DB.New()
	if err != nil {
		p.Log.Emit(sinks.Error(err).WithFields(sink.Fields{
			"page":            page,
			"responsePerPage": responsePerPage,
		}))
		return ProfileRecords{}, err
	}

	defer dbi.Close()

	var nu profile.Profile
	records, realTotalRecords, err := db.GetAllPerPage(p.Log, dbi, nu, "asc", "public_id", page, responsePerPage)
	if err != nil {
		p.Log.Emit(sinks.Error(err).WithFields(sink.Fields{
			"page":            page,
			"responsePerPage": responsePerPage,
		}))
		return ProfileRecords{}, err
	}

	var profileRecords []profile.Profile

	for _, record := range records {
		var nw profile.Profile

		if err := nw.WithFields(record); err != nil {
			p.Log.Emit(sinks.Error(err).WithFields(sink.Fields{
				"page":            page,
				"responsePerPage": responsePerPage,
			}))
			return ProfileRecords{}, err
		}

		profileRecords = append(profileRecords, nw)
	}

	return ProfileRecords{
		Page:          page,
		Total:         realTotalRecords,
		Records:       profileRecords,
		ResponserPage: responsePerPage,
	}, nil
}

// Get retrieves the profile associated with the giving profile_id.
func (p Profiles) Get(profileID string) (*profile.Profile, error) {
	defer p.Log.Emit(sinks.Info("Get Existing Profile").WithFields(sink.Fields{
		"profile_id": profileID,
	}).Trace("Profiles.Get").End())

	dbi, err := p.DB.New()
	if err != nil {
		p.Log.Emit(sinks.Error(err).WithFields(sink.Fields{"profile_id": profileID}))
		return nil, err
	}

	defer dbi.Close()

	var existingProfile profile.Profile

	// Attempt to retrieve profile from db if we still have an outstanding non-expired profile.
	if err := db.Get(p.Log, dbi, existingProfile, &existingProfile, "public_id", profileID); err != nil {
		p.Log.Emit(sinks.Error("Failed to retrieve profile from db: %+q", err).WithFields(sink.Fields{"profile_id": profileID}))
		return nil, err
	}

	return &existingProfile, nil
}

// GetByUser retrieves the profile associated with the giving UserID.
func (p Profiles) GetByUser(userID string) (*profile.Profile, error) {
	defer p.Log.Emit(sinks.Info("Get Existing Profile").WithFields(sink.Fields{
		"user_id": userID,
	}).Trace("Profiles.GetByUser").End())

	dbi, err := p.DB.New()
	if err != nil {
		p.Log.Emit(sinks.Error(err).WithFields(sink.Fields{"user_id": userID}))
		return nil, err
	}

	defer dbi.Close()

	var existingProfile profile.Profile

	// Attempt to retrieve profile from db if we still have an outstanding non-expired profile.
	if err := db.Get(p.Log, dbi, existingProfile, &existingProfile, profile.UniqueIndex, userID); err != nil {
		p.Log.Emit(sinks.Error("Failed to retrieve profile from db: %+q", err).WithFields(sink.Fields{"user_id": userID}))
		return nil, err
	}

	return &existingProfile, nil
}

// DeleteByUser removes an existing profile from the db for a specified profile.
func (p Profiles) DeleteByUser(userID string) error {
	defer p.Log.Emit(sinks.Info("Delete Existing Profile").WithFields(sink.Fields{
		"user_id": userID,
	}).Trace("Profiles.DeleteByUser").End())

	dbi, err := p.DB.New()
	if err != nil {
		p.Log.Emit(sinks.Error(err).WithFields(sink.Fields{"user_id": userID}))
		return err
	}

	defer dbi.Close()

	var ns profile.Profile

	// Delete this profile
	if err := db.Delete(p.Log, dbi, ns, profile.UniqueIndex, userID); err != nil {
		p.Log.Emit(sinks.Error("Failed to delete profile profile from db: %+q", err).WithFields(sink.Fields{"user_id": userID}))
		return err
	}

	return nil
}

// Delete removes an existing profile from the db for a specified profile by its id.
func (p Profiles) Delete(profileID string) error {
	defer p.Log.Emit(sinks.Info("Delete Existing Profile").WithFields(sink.Fields{
		"profile_id": profileID,
	}).Trace("Profiles.Delete").End())

	dbi, err := p.DB.New()
	if err != nil {
		p.Log.Emit(sinks.Error(err).WithFields(sink.Fields{"profile_id": profileID}))
		return err
	}

	defer dbi.Close()

	var ns profile.Profile

	// Delete this profile
	if err := db.Delete(p.Log, dbi, ns, "public_id", profileID); err != nil {
		p.Log.Emit(sinks.Error("Failed to delete profile from db: %+q", err).WithFields(sink.Fields{"profile_id": profileID}))
		return err
	}

	return nil
}

// Update handles receiving requests to update a profile identified by it's public_id.
func (p Profiles) Update(nw profile.UpdateProfile) error {
	defer p.Log.Emit(sinks.Info("Update User").With("profile_id", nw.PublicID).Trace("handlers.Users.Update").End())

	if nw.PublicID == "" {
		err := errors.New("JSON User.PublicID is empty")
		p.Log.Emit(sinks.Error(err).WithFields(sink.Fields{
			"profile_id": nw.PublicID,
		}))

		return err
	}

	dbi, err := p.DB.New()
	if err != nil {
		p.Log.Emit(sinks.Error(err).WithFields(sink.Fields{
			"profile_id": nw.PublicID,
		}))

		return err
	}

	defer dbi.Close()

	if err := db.Update(p.Log, dbi, nw, "public_id"); err != nil {
		p.Log.Emit(sinks.Error(err).WithFields(sink.Fields{
			"profile_id": nw.PublicID,
		}))

		return err
	}

	return nil
}
