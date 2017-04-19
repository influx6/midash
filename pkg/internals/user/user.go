package user

import (
	"golang.org/x/crypto/bcrypt"

	"time"

	"github.com/fatih/structs"
	uuid "github.com/satori/go.uuid"
)

const (
	tableName  = "users"
	timeFormat = "Mon Jan 2 15:04:05 -0700 MST 2006"
)

// User is a type defining the given user related fields for a given.
type User struct {
	Email     string    `json:"email"`
	PublicID  string    `json:"public_id"`
	PrivateID string    `json:"private_id"`
	Hash      string    `json:"hash"`
	Created   time.Time `json:"created_at"`
	Updated   time.Time `json:"updated_at"`
}

// NewUser defines the set of data received to create a new user.
type NewUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// New returns a new User instance based on the provided data.
func New(nw NewUser) (*User, error) {
	var u User
	u.Email = nw.Email
	u.PublicID = uuid.NewV4().String()
	u.PrivateID = uuid.NewV4().String()

	pass := []byte(u.PrivateID + ":" + nw.Password)
	hash, err := bcrypt.GenerateFromPassword(pass, 20)
	if err != nil {
		return nil, err
	}

	u.Hash = string(hash)
	u.Created = time.Now().UTC()
	u.Updated = time.Now().UTC()

	return &u, nil
}

// Authenticate attempts to authenticate the giving password to the provided user.
func (u *User) Authenticate(password string) error {
	pass := []byte(u.PrivateID + ":" + password)
	return bcrypt.CompareHashAndPassword([]byte(u.Hash), pass)
}

// Table returns the given table which the given struct corresponds to.
func (u User) Table() string {
	return tableName
}

// WithFields attempts to syncing the giving data within the provided
// map into it's own fields.
func (u *User) WithFields(fields map[string]interface{}) error {
	if email, ok := fields["email"].(string); ok {
		u.Email = email
	}

	if public, ok := fields["public_id"].(string); ok {
		u.PublicID = public
	}

	if private, ok := fields["private_id"].(string); ok {
		u.PrivateID = private
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

// SafeFields returns a map representing the data of the user with important
// security fields removed.
func (u *User) SafeFields() map[string]interface{} {
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

	delete(fields, "hash")
	delete(fields, "private_id")

	return fields
}

// Fields returns a map representing the data of the user.
func (u *User) Fields() map[string]interface{} {
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
