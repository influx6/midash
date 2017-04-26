package controllers

import (
	"net/http"

	"github.com/gu-io/midash/pkg/internals/handlers"
	"github.com/gu-io/midash/pkg/internals/utils"
	"github.com/influx6/faux/sink"
	"github.com/influx6/faux/sink/sinks"
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
func (u Auth) CheckAuthorization(w http.ResponseWriter, r *http.Request, params map[string]string) {
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
		return
	}
}
