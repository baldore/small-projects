package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "asdfasdf"
	dbname   = "react_golang_jwt"
)

func main() {
	// database connection
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("error creating connection to db: %v", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatalf("error trying to ping db: %v", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("error trying to create db driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://./server/migrations", "postgres", driver)
	if err != nil {
		log.Fatalf("error trying to create migrate client: %v", err)
	}

	err = m.Up()
	if err != nil {
		log.Fatalf("error running migrate.up: %v", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		SetJSONHeader(w)
		w.Write([]byte(`{ "message": "hola mundo genial" }`))
	})

	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Panicf("error starting the server: %v", err)
	}
}

func SetJSONHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}
