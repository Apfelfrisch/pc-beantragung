package database

import (
	"database/sql"
	_ "embed"
	"log"
	"os"
	"pc-beantragung/internal/signon"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed blank-sqlite.sqlite
var sqlite []byte

// Service represents a service that interacts with a database.
type Service interface {
	DbInstance() *sql.DB
	SignonRepo() signon.SignOnRepo
	Close() error
}

type service struct {
	db *sql.DB
}

var (
	dburl      = "./signons.sqlite"
	dbInstance *service
)

func New() Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}

	if _, err := os.Stat(dburl); err != nil {
		file, _ := os.Create("./signons.sqlite")
		defer file.Close()
		_, err = file.Write(sqlite)
		file.Sync()
	}

	db, err := sql.Open("sqlite3", dburl)

	if err != nil {
		log.Fatal(err)
	}

	dbInstance = &service{
		db: db,
	}
	return dbInstance
}

func (s *service) DbInstance() *sql.DB {
	return s.db
}

func (s *service) SignonRepo() signon.SignOnRepo {
	return signon.SignonRepo(signon.New(s.db), s.db)
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *service) Close() error {
	log.Printf("Disconnected from database: %s", dburl)
	return s.db.Close()
}
