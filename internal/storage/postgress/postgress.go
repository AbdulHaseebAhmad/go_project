package postgress

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	//
	"github.com/AbdulHaseebAhmad/go_project/internal/config"
	"github.com/AbdulHaseebAhmad/go_project/internal/types"
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

func (s *Postgres) CreateStudent(name string, email string, age int) (student types.Student, error error) {
	var newStudent types.Student

	qerr := s.DB.QueryRow(
		"INSERT INTO students (name,email,age) VALUES ($1, $2, $3) RETURNING id, name, email, age",
		name, email, age,
	).Scan(
		&newStudent.Id,
		&newStudent.Name,
		&newStudent.Email,
		&newStudent.Age,
	)
	if qerr != nil {
		return types.Student{}, nil
	}

	return newStudent, nil
}

func (s *Postgres) GetStudentById(id int64) (student types.Student, eror error) {
	stmt, err := s.DB.Prepare("SELECT * FROM students WHERE id=$1 LIMIT 1")
	if err != nil {
		slog.Info("Error Fetching Student: " + err.Error())
		return types.Student{}, err
	}

	defer stmt.Close()
	var returnedStudent types.Student
	qerr := stmt.QueryRow(id).Scan(&returnedStudent.Id, &returnedStudent.Name, &returnedStudent.Email, &returnedStudent.Age) // execute the query and pass it to struct with Scan
	if qerr != nil {
		if qerr == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("no students found with id: %s", fmt.Sprint(id))

		}
		return types.Student{}, fmt.Errorf("query error: %w", qerr)
	}
	return returnedStudent, nil
}
func (s *Postgres) GetStudentList() (stuentsList []types.Student, eror error) {
	stmt, err := s.DB.Prepare("SELECT * FROM students")
	if err != nil {
		slog.Info("Error Fetching Students: " + err.Error())
		return nil, err
	}

	defer stmt.Close()
	rows, qerr := stmt.Query()
	if qerr != nil {
		return nil, fmt.Errorf("query error: %w", qerr)
	}

	defer rows.Close()
	var students []types.Student

	for rows.Next() {
		var returnedStudent types.Student
		err := rows.Scan(&returnedStudent.Id, &returnedStudent.Name, &returnedStudent.Email, &returnedStudent.Age)
		if err != nil {
			return nil, err
		}
		students = append(students, returnedStudent)
	}
	return students, nil
}

func (s *Postgres) DeleteStudent(id int64) (studentId types.Student, error error) {
	// delete the student,return the id, if deletion unsuccessfull return error
	stmt, err := s.DB.Prepare("DELETE FROM students WHERE id=$1")
	if err != nil {
		slog.Info("Error Prepping Query")
		return types.Student{}, err
	}
	defer stmt.Close()
	var returnedStudent types.Student
	qerr := stmt.QueryRow(id).Scan(&returnedStudent.Id, &returnedStudent.Name, &returnedStudent.Email, &returnedStudent.Age)
	if qerr != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Info("No User Found")
			return types.Student{}, sql.ErrNoRows
		}
	}

	return returnedStudent, nil
}
