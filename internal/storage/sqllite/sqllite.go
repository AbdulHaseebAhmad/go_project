package sqllite

import (
	"database/sql"

	//
	"github.com/AbdulHaseebAhmad/go_project/internal/config"
	_ "github.com/lib/pq"
)

type Postgres struct {
	DB *sql.DB
}

func New(cfg *config.Config) (*Postgres, error) {

	// open connection
	db, err := sql.Open("postgres", cfg.Storage_path) //sql is the standard library for sql connection, 1st arg db name 2nd arg path
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &Postgres{DB: db}, nil
}

func (s *Postgres) CreateStudent(name string, email string, age int) (id int64, error error) {
	stmt, err := s.DB.Prepare("INSERT INTO students (name,email,age) VALUES ($1, $2, $3)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastId, nil
}
