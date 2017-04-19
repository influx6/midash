package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/dimfeld/httptreemux"
	_ "github.com/go-sql-driver/mysql" // loads up the go mysql driver.
	"github.com/gu-io/midash/pkg/controllers"
	"github.com/jmoiron/sqlx"
)

// contains different environment flags for use to setting up
// a db connection.
const (
	PortEnv       = "PORT"
	APIVersionENV = "API_Version"
	DBPortEnv     = "MYSQL_PORT"
	DBUserEnv     = "MYSQL_USER"
	DBDatabaseEnv = "MYSQL_DATABASE"
	DBUserPassEnv = "MYSQL_PASSWORD"
)

type djDB struct{}

// New returns a new instance of a sqlx.DB connected to the db with the provided
// credentials pulled from the host environment.
func (djDB) New() (*sqlx.DB, error) {
	user := os.Getenv(DBUserEnv)
	userPass := os.Getenv(DBUserPassEnv)
	// port := os.Getenv(DBPortEnv)
	dbName := os.Getenv(DBDatabaseEnv)

	addr := fmt.Sprintf("%s:%s@/%s", user, userPass, dbName)
	db, err := sqlx.Connect("mysql", addr)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func main() {

	// Get API version.
	version := os.Getenv(APIVersionENV)
	if version == "" {
		version = "v1"
	}

	// Get the App port.
	port := os.Getenv(PortEnv)

	var dj djDB

	users := controllers.Users{DB: dj}

	tree := httptreemux.New()

	tree.Handle("POST", fmt.Sprintf("%s/user", version), users.CreateUser)

	cm := make(chan os.Signal, 1)
	signal.Notify(cm, os.Interrupt)

	srv := &http.Server{Addr: ":" + port, Handler: tree}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()

	<-cm
	log.Println("Shutting down server...")

	// shut down gracefully, but wait no longer than 5 seconds before halting
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()
	srv.Shutdown(ctx)

	log.Println("Server gracefully stopped")
}
