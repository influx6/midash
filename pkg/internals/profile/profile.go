package session

import (
	"time"

	"github.com/fatih/structs"
)

const (
	tableName  = "profiles"
	timeFormat = "Mon Jan 2 15:04:05 -0700 MST 2006"
)

// Profile defines a struct which holds the the details of a giving user's profile.
type Profile struct {
	Address   string    `json:"address"`
	UserID    string    `json:"user_id"`
	PublicID  string    `json:"public_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Created   time.Time `json:"created_at"`
	Updated   time.Time `json:"updated_at"`
}

// Table returns the given table which the given struct corresponds to.
func (u Profile) Table() string {
	return tableName
}

// WithFields attempts to syncing the giving data within the provided
// map into it's own fields.
func (u *Profile) WithFields(fields map[string]interface{}) error {
	if user, ok := fields["user_id"].(string); ok {
		u.UserID = user
	}

	if public, ok := fields["public_id"].(string); ok {
		u.PublicID = public
	}

	if firstname, ok := fields["first_name"].(string); ok {
		u.FirstName = firstname
	}

	if lastname, ok := fields["last_name"].(string); ok {
		u.LastName = lastname
	}

	if address, ok := fields["address"].(string); ok {
		u.Address = address
	}

	if created, ok := fields["created_at"]; ok {
		switch co := created.(type) {
		case string:
			t, err := time.Parse(timeFormat, co)
			if err != nil {
				return err
			}

			u.Created = t.UTC()
		case time.Time:
			u.Created = co.UTC()
		}
	}

	if updated, ok := fields["updated_at"]; ok {
		switch co := updated.(type) {
		case string:
			t, err := time.Parse(timeFormat, co)
			if err != nil {
				return err
			}

			u.Updated = t.UTC()
		case time.Time:
			u.Updated = co.UTC()
		}
	}

	return nil
}

// Fields returns a map representing the data of the user.
func (u *Profile) Fields() map[string]interface{} {
	fields := structs.Map(u)

	if u.Created.IsZero() {
		delete(fields, "created_at")
	} else {
		fields["created_at"] = u.Created.UTC()
	}

	if u.Updated.IsZero() {
		delete(fields, "updated_at")
	} else {
		fields["updated_at"] = u.Updated.UTC()
	}

	return fields
}
