package postgress

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	//

	passwordhash "github.com/AbdulHaseebAhmad/go_project/internal/Utils/hash"
	"github.com/AbdulHaseebAhmad/go_project/internal/Utils/tokens"
	"github.com/AbdulHaseebAhmad/go_project/internal/config"
	"github.com/AbdulHaseebAhmad/go_project/internal/types"
	"github.com/lib/pq"
	// _ "github.com/lib/pq"
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

func (s *Postgres) DeleteStudent(id int64) error {
	// delete the student,return the id, if deletion unsuccessfull return error
	stmt, err := s.DB.Prepare("DELETE FROM students WHERE id=$1")
	if err != nil {
		slog.Info("Error Prepping Query")
		return err
	}
	defer stmt.Close()
	result, qerr := stmt.Exec(id)
	if qerr != nil {
		return qerr
	}
	rowsAffected, rerr := result.RowsAffected()
	if rerr != nil {
		return rerr
	}

	if rowsAffected == 0 {
		slog.Info("No user found to delete")
		return sql.ErrNoRows
	}
	return nil
}

func (s *Postgres) RegisterStudent(ctx context.Context, credentials types.Credentials) (email string, err error) {
	hashedPassword, hasherror := passwordhash.Hashpassword(credentials.Password)
	if hasherror != nil {
		return "", fmt.Errorf("failed to insert user: %w", hasherror)
	}
	var newEmail string
	qerer := s.DB.QueryRowContext(ctx, "INSERT INTO users (email, hashed_password) VALUES ($1, $2) RETURNING email",
		credentials.Email, hashedPassword).Scan(&newEmail)

	if qerer != nil {
		pqErr, ok := qerer.(*pq.Error) // check if the qerer is of Pq.Error type, if it is it will return ok true and it will return qerer in pqErr
		if ok {
			if pqErr.Code == "23505" {
				return "", fmt.Errorf("user already exists with email %s", credentials.Email)
			}
		}
		return "", fmt.Errorf("failed to insert user: %w", qerer)
	}

	return fmt.Sprintf("User with email address %s Registered", newEmail), nil
}

func (s *Postgres) StudentSignin(ctx context.Context, credentials types.Login) (string, string, error) {
	var retrievedPassword string
	var role string
	qerer := s.DB.QueryRowContext(ctx, "SELECT hashed_password,role from users where email=$1", credentials.Email).Scan(&retrievedPassword, &role)
	if qerer != nil {
		if errors.Is(qerer, sql.ErrNoRows) {
			return "", "", fmt.Errorf("invalid email or password")
		}
		return "", "", qerer
	}
	_, passerror := passwordhash.Unhashpassword(credentials.Password, retrievedPassword)
	if passerror != nil {
		return "", "", fmt.Errorf("invalid email or password")
	}
	sessiontoken, terr := tokens.GenerateToken(32)
	if terr != nil {
		return "", "", terr
	}
	csrfToken, cterr := tokens.GenerateToken(32)
	if cterr != nil {
		return "", "", terr
	}
	_, tqerr := s.DB.ExecContext(
		ctx,
		`INSERT INTO sessions (session_token, csrf_token, email, role, created_at)
     VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP)
     ON CONFLICT (email)
     DO UPDATE SET
         session_token = EXCLUDED.session_token,
         csrf_token = EXCLUDED.csrf_token,
         role = EXCLUDED.role,
         created_at = CURRENT_TIMESTAMP`,
		sessiontoken, csrfToken, credentials.Email, role,
	)
	if tqerr != nil {
		return "", "", tqerr
	}
	return sessiontoken, csrfToken, nil
}

func (s *Postgres) AuthorizeStudent(ctx context.Context, sessiontoken string, csrfHeader string) bool {
	var csrfToken string
	qerr := s.DB.QueryRowContext(ctx, "SELECT csrf_token from sessions WHERE session_token = $1", sessiontoken).Scan(&csrfToken)
	if qerr != nil {
		if errors.Is(qerr, sql.ErrNoRows) {
			return false
		}
		return false
	}

	return csrfToken == csrfHeader
}
