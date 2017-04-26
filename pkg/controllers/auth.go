package controllers

import (
	"errors"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/gu-io/midash/pkg/internals/auth"
	"github.com/gu-io/midash/pkg/internals/handlers"
	"github.com/gu-io/midash/pkg/internals/utils"
	"github.com/influx6/faux/sink"
	"github.com/influx6/faux/sink/sinks"
	"golang.org/x/oauth2"
)

// Auth defines an handler which provides authorization handling for
// a request, needing user authentication.
type Auth struct {
	handlers.BearerAuth
}

// CheckAuthorization handles receiving requests to verify user authorization.
/* Service API
HTTP Method: GET
Header:
		{
			"Authorization":"Bearer <TOKEN>",
		}

		WHERE: <TOKEN> = <USERID>:<SESSIONTOKEN>
*/
func (u Auth) CheckAuthorization(w http.ResponseWriter, r *http.Request, params map[string]string) error {
	defer u.Log.Emit(sinks.Info("Authenticate Authorization").WithFields(sink.Fields{
		"params": params,
		"remote": r.RemoteAddr,
		"path":   r.URL.Path,
	}).Trace("Auth.CheckAuthorization").End())

	// Retrieve authorization header.
	if err := u.BearerAuth.CheckAuthorization(r.Header.Get("Authorization")); err != nil {
		u.Log.Emit(sinks.Error(err).WithFields(sink.Fields{
			"path":   r.URL.Path,
			"remote": r.RemoteAddr,
			"params": params,
		}))

		http.Error(w, utils.ErrorMessage(http.StatusInternalServerError, "Invalid Auth: Failed to validate authorization", err), http.StatusInternalServerError)
		return err
	}

	return nil
}

//==================================================================================================================================================================

// OAuth defines a controller which handles the incoming request that it contains the giving "secret"
// within it's data.
type OAuth struct {
	Auth    *auth.Auth
	Options []oauth2.AuthCodeOption
	Log     sink.Sink
}

// Redirect attempts to redirect incoming request with the OAuth URL from the supplied OAuth
// structure and uses the giving secret state to generate the URL to redirect to.
func (u *OAuth) Redirect(secret string, w http.ResponseWriter, r *http.Request) (string, error) {
	defer u.Log.Emit(sinks.Info("Redirect Request to OAuth.URL").WithFields(sink.Fields{
		"remote": r.RemoteAddr,
		"path":   r.URL.Path,
	}).Trace("OAuth.Redirect").End())

	return u.Auth.LoginURL(secret, u.Options...), nil
}

// Validate that the giving SecretCode matches the incoming value of the request else returns an error.
func (u *OAuth) Validate(secret string, w http.ResponseWriter, r *http.Request) error {
	defer u.Log.Emit(sinks.Info("Validated OAuth Secret in Request").WithFields(sink.Fields{
		"remote": r.RemoteAddr,
		"path":   r.URL.Path,
	}).Trace("OAuth.Validate").End())

	stateSecret := r.FormValue("state")
	if stateSecret != secret {
		err := errors.New("Invalid OAuth secret")
		u.Log.Emit(sinks.Error("OAuth State Fails to match: %+q", err).WithFields(sink.Fields{
			"path":   r.URL.Path,
			"remote": r.RemoteAddr,
		}))

		return err
	}

	return nil
}

//==================================================================================================================================================================

// Guarded defines a struct which exposes a session secured request life cycle where a request made will be guarded
// with specific data from a underline session and will be validated when receiving response.
type Guarded struct {
	SessionName  string
	CookieName   string
	CookieSecret string
	Cookies      sessions.CookieStore
	Log          sink.Sink
}

// Guard attempts to added incoming request with a session which is stored in the outgoing response which
// then will be used to guard against other incoming request.
func (u Guarded) Guard(w http.ResponseWriter, r *http.Request) error {
	defer u.Log.Emit(sinks.Info("Guard Request").WithFields(sink.Fields{
		"remote": r.RemoteAddr,
		"path":   r.URL.Path,
	}).Trace("Guarded.Guard").End())

	defer context.Clear(r)

	session, err := u.Cookies.Get(r, u.SessionName)
	if err != nil {
		u.Log.Emit(sinks.Error("Cookie Retreival Failed: %+q", err).WithFields(sink.Fields{
			"path":   r.URL.Path,
			"remote": r.RemoteAddr,
		}))

		return err
	}

	session.Values[u.CookieName] = u.CookieSecret
	session.Save(r, w)

	return nil
}

// Validate attempts to authenticate incoming request with a sessio data expected from the request.
func (u Guarded) Validate(w http.ResponseWriter, r *http.Request) error {
	defer u.Log.Emit(sinks.Info("Validated Guarded Request").WithFields(sink.Fields{
		"remote": r.RemoteAddr,
		"path":   r.URL.Path,
	}).Trace("Guarded.Validate").End())

	defer context.Clear(r)

	session, err := u.Cookies.Get(r, u.SessionName)
	if err != nil {
		u.Log.Emit(sinks.Error("Cookie Retreival Failed: %+q", err).WithFields(sink.Fields{
			"path":   r.URL.Path,
			"remote": r.RemoteAddr,
		}))

		return err
	}

	// Attempt to retrieve specific Guard.CookieName in retrieved session.
	value, ok := session.Values[u.CookieName]
	if !ok {
		err := errors.New("Session cookie guard not found in request")
		u.Log.Emit(sinks.Error(err).WithFields(sink.Fields{
			"path":   r.URL.Path,
			"remote": r.RemoteAddr,
		}))

		return err
	}

	// Did value match expected guard secret?
	if value != u.CookieSecret {
		err := errors.New("Session cookie guard not found in request")
		u.Log.Emit(sinks.Error(err).WithFields(sink.Fields{
			"path":   r.URL.Path,
			"remote": r.RemoteAddr,
		}))

		return err
	}

	return nil
}
