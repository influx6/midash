package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/dimfeld/httptreemux"
	_ "github.com/go-sql-driver/mysql" // loads up the go mysql driver.
	"github.com/influx6/backoffice/db"
	"github.com/influx6/backoffice/db/sql"
	"github.com/influx6/backoffice/handlers"
	"github.com/influx6/backoffice/migrations/sqltables"
	"github.com/influx6/backoffice/resources"
	"github.com/influx6/faux/naming"
	"github.com/influx6/faux/sink"
	"github.com/influx6/faux/sink/sinks"
)

// contains different environment flags for use to setting up
// a db connection.
const (
	PortEnv            = "PORT"
	APIVersionENV      = "API_Version"
	DBPortEnv          = "MYSQL_PORT"
	DBIPEnv            = "MYSQL_IP"
	DBUserEnv          = "MYSQL_USER"
	DBDatabaseEnv      = "MYSQL_DATABASE"
	DBUserPassEnv      = "MYSQL_PASSWORD"
	GoogleClientID     = "GOOGLE_CLIENT_ID"
	GoogleClientSecret = "GOOGLE_CLIENT_SECRET"
)

var (
	user     = strings.TrimSpace(os.Getenv(DBUserEnv))
	userPass = strings.TrimSpace(os.Getenv(DBUserPassEnv))
	ip       = strings.TrimSpace(os.Getenv(DBIPEnv))
	dbName   = strings.TrimSpace(os.Getenv(DBDatabaseEnv))
)

func main() {
	log := sink.New(sinks.Stdout{})

	// Get API version.
	version := strings.TrimSpace(os.Getenv(APIVersionENV))
	if version == "" {
		version = "v1"
	}

	port, _ := strconv.Atoi(strings.TrimSpace(os.Getenv(DBPortEnv)))

	// Created Namer and sql connection creator.
	appNamer := naming.NewNamer("%s_%s", naming.PrefixNamer{Prefix: "midash"})

	sqlDB := sql.New(log, sql.Conn{
		User:     user,
		Password: userPass,
		Addr:     ip,
		Port:     port,
		Database: dbName,
		Log:      log,
		Driver:   "mysql",
	}, sqltables.BasicTables(appNamer)...)

	// Create table namers which tell the models which database table to use.
	usersTable := db.TableName{Name: appNamer.New("users")}
	profilesTable := db.TableName{Name: appNamer.New("profiles")}
	sessionsTable := db.TableName{Name: appNamer.New("sessions")}

	// Create the resources model handlers we need to handle.
	sessionModels := handlers.SessionsFactory(log, sqlDB, 72*time.Hour, sessionsTable)
	profileModels := handlers.ProfilesFactory(log, sqlDB, profilesTable)
	userModels := handlers.UsersFactory(log, sqlDB, usersTable, profilesTable)
	auth := handlers.BearerAuth{Users: userModels, Sessions: sessionModels}

	// Create API resources request handlers
	users := resources.Users{Users: userModels}
	sessions := resources.Sessions{Sessions: sessionModels, Users: userModels}
	profiles := resources.Profiles{Users: userModels, Sessions: sessionModels, Profiles: profileModels}

	tree := httptreemux.New()
	tree.Handle("GET", "/", index)
	tree.Handle("GET", fmt.Sprintf("/%s", version), welcome(version))

	// Set up sessions related routes.
	tree.Handle("POST", fmt.Sprintf("/%s/%s", version, "sessions/login"), sessions.Login)
	tree.Handle("POST", fmt.Sprintf("/%s/%s", version, "sessions/logout"), sessions.Logout)

	// Set up users related routes.
	tree.Handle("GET", fmt.Sprintf("/%s/%s", version, "users"), users.GetAll)
	tree.Handle("GET", fmt.Sprintf("/%s/%s", version, "users/:total/:page"), users.GetAll)

	tree.Handle("POST", fmt.Sprintf("/%s/%s", version, "users"), users.Create)
	tree.Handle("GET", fmt.Sprintf("/%s/%s", version, "users/:public_id"), users.GetLimited)

	// Set up profiles related routes.
	tree.Handle("PUT", fmt.Sprintf("/%s/%s", version, "profiles"), resources.Auth{
		BearerAuth: auth,
		Next:       profiles.Update,
	}.CheckAuthorization)

	tree.Handle("GET", fmt.Sprintf("/%s/%s", version, "profiles/:public_id"), resources.Auth{
		BearerAuth: auth,
		Next:       profiles.Get,
	}.CheckAuthorization)

	tree.Handle("GET", fmt.Sprintf("/%s/%s", version, "profiles/users/:user_id"), resources.Auth{
		BearerAuth: auth,
		Next:       profiles.GetForUser,
	}.CheckAuthorization)

	cm := make(chan os.Signal, 1)
	signal.Notify(cm, os.Interrupt)

	// Get the App port.
	appPort := strings.TrimSpace(os.Getenv(PortEnv))
	addr := fmt.Sprintf(":%s", appPort)

	srv := &http.Server{Addr: addr, Handler: tree}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Emit(sinks.Error("Failed to start server: %+q", err).With("addr", addr))
		}
	}()

	log.Emit(sinks.Info("HTTP server started").With("addr", addr))

	<-cm
	log.Emit(sinks.Info("Shutting down server").With("addr", addr))

	// shut down gracefully, but wait no longer than 5 seconds before halting
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()
	srv.Shutdown(ctx)

	log.Emit(sinks.Info("Server gracefully stopped").With("addr", addr))
}

func welcome(version string) httptreemux.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Welcome to midash version " + version))
	}
}

func index(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome to midash"))
}
