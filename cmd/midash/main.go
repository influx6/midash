package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/dimfeld/httptreemux"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gu-io/midash/pkg/handlers" // loads up the go mysql driver.
	"github.com/influx6/faux/sink"
	"github.com/influx6/faux/sink/sinks"
	"github.com/jmoiron/sqlx"
)

// contains different environment flags for use to setting up
// a db connection.
const (
	PortEnv       = "PORT"
	APIVersionENV = "API_Version"
	DBPortEnv     = "MYSQL_PORT"
	DBIPEnv       = "MYSQL_IP"
	DBUserEnv     = "MYSQL_USER"
	DBDatabaseEnv = "MYSQL_DATABASE"
	DBUserPassEnv = "MYSQL_PASSWORD"
)

var (
	log = sink.New(sinks.Stdout{})
)

type djDB struct{}

// New returns a new instance of a sqlx.DB connected to the db with the provided
// credentials pulled from the host environment.
func (djDB) New() (*sqlx.DB, error) {
	user := strings.TrimSpace(os.Getenv(DBUserEnv))
	userPass := strings.TrimSpace(os.Getenv(DBUserPassEnv))
	port := strings.TrimSpace(os.Getenv(DBPortEnv))
	ip := strings.TrimSpace(os.Getenv(DBIPEnv))
	dbName := strings.TrimSpace(os.Getenv(DBDatabaseEnv))

	if ip == "" {
		ip = "0.0.0.0"
	}

	addr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, userPass, ip, port, dbName)
	db, err := sqlx.Connect("mysql", addr)
	if err != nil {
		log.Emit(sinks.Error("Failed to connect to SQLServer: %+q", err).WithFields(sink.Fields{
			"mysql_ip":   ip,
			"mysql_port": port,
			"dbName":     dbName,
			"user":       user,
		}))

		return nil, err
	}

	return db, nil
}

func main() {

	// Get API version.
	version := strings.TrimSpace(os.Getenv(APIVersionENV))
	if version == "" {
		version = "v1"
	}

	// Get the App port.
	port := strings.TrimSpace(os.Getenv(PortEnv))
	addr := fmt.Sprintf(":%s", port)

	var dj djDB

	tree := httptreemux.New()

	users := handlers.Users{DB: dj, Log: log}

	tree.Handle("GET", "/", index)
	tree.Handle("GET", fmt.Sprintf("/%s", version), welcome(version))

	tree.Handle("POST", fmt.Sprintf("/%s/users", version), users.CreateUser)

	cm := make(chan os.Signal, 1)
	signal.Notify(cm, os.Interrupt)

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
