package storage

import (
	"context"

	"github.com/AbdulHaseebAhmad/go_project/internal/types"
)

type Storage interface {
	CreateStudent(name string, email string, age int) (student types.Student, error error)
	GetStudentById(id int64) (student types.Student, error error)
	GetStudentList() ([]types.Student, error)
	DeleteStudent(id int64) error
	RegisterStudent(context.Context, types.Credentials) (email string, error error)
	StudentSignin(context.Context, types.Login) (string, string, error)
	AuthorizeStudent(context context.Context, sessiontoken string, csrfHeader string) bool
}
