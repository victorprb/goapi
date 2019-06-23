package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/victorprb/goapi/pkg/models"
	"github.com/victorprb/goapi/pkg/models/psql"

	_ "github.com/lib/pq"
)

// Config for application
type Config struct {
	ListenAddr string
	DbHost string
	DbPort string
	DbUser string
	DbPass string
	DbName string
}

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	users interface{
		Insert(*models.User) error
	}
}

func main() {
	cfg := new(Config)
	flag.StringVar(&cfg.ListenAddr, "addr", ":8080", "HTTP network address")
	flag.StringVar(&cfg.DbHost, "dbhost", "localhost", "DB host address")
	flag.StringVar(&cfg.DbPort, "dbport", "5432", "DB port")
	flag.StringVar(&cfg.DbUser, "dbuser", "postgres", "DB username")
	flag.StringVar(&cfg.DbPass, "dbpass", "postgres", "DB user password")
	flag.StringVar(&cfg.DbName, "dbname", "goapi", "DB name")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DbHost,
		cfg.DbPort,
		cfg.DbUser,
		cfg.DbPass,
		cfg.DbName,
	)
	db, err := openDB(dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		users: &psql.UserModel{DB: db},
	}

	srv := &http.Server{
		Addr:         cfg.ListenAddr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 10,
	}

	infoLog.Printf("Starting server on %s", cfg.ListenAddr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

// The openDB() function wraps sql.Open() and returns a sql.DB connection pool
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
