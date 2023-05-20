package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
)

var Db *sql.DB

func NewDatabase() {
	connectionString := fmt.Sprintf(
		"root:%s@tcp(localhost)/%s",
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_DB"),
	)
	fmt.Println("connection string:", connectionString)
	db, err := sql.Open("mysql", "root:dbpass@tcp(db)/hackernews")
	if err != nil {
		log.Panic(err)
	}

	if err = db.Ping(); err != nil {
		log.Panic(err)
	}

	Db = db
}

func CloseDatabase() error {
	return Db.Close()
}

func Migrate() {
	if err := Db.Ping(); err != nil {
		log.Fatal(err)
	}

	driver, _ := mysql.WithInstance(Db, &mysql.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/pkg/db/migrations/mysql",
		"mysql",
		driver,
	)

	if err != nil {
		log.Panic(err)
	}

	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			log.Panic(err)
		}
	}

	fmt.Println("successfully migrated the database")
}
