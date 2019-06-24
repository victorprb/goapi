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

type config struct {
	listenAddr string
	secret string
	dbHost string
	dbPort string
	dbUser string
	dbPass string
	dbName string
}

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	secretKey string
	users interface{
		Insert(*models.User) error
		Authenticate(*models.Credentials) (string, error)
		Get(string) (*models.User, error)
	}
}

func main() {
	// Server config
	cfg := new(config)
	flag.StringVar(&cfg.listenAddr, "addr", ":8080", "HTTP network address")
	flag.StringVar(&cfg.secret, "secret", "uQxaB8WAMqVW#a_+Rr-Nq=+Jkpvc?kSF", "Secret key")
	flag.StringVar(&cfg.dbHost, "dbhost", "localhost", "DB host address")
	flag.StringVar(&cfg.dbPort, "dbport", "5432", "DB port")
	flag.StringVar(&cfg.dbUser, "dbuser", "postgres", "DB username")
	flag.StringVar(&cfg.dbPass, "dbpass", "postgres", "DB user password")
	flag.StringVar(&cfg.dbName, "dbname", "goapi", "DB name")

	flag.Parse()

	// Custom loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Database info connection
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.dbHost,	cfg.dbPort,	cfg.dbUser,	cfg.dbPass,	cfg.dbName)

	// DB connection
	db, err := openDB(dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		secretKey: cfg.secret,
		users: &psql.UserModel{DB: db},
	}

	srv := &http.Server{
		Addr:         cfg.listenAddr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 10,
	}

	// Initiate the server
	infoLog.Printf("Starting server on %s", cfg.listenAddr)
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
