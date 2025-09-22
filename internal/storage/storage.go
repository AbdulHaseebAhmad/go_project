package storage

import (
	"context"

	"github.com/AbdulHaseebAhmad/go_project/internal/types"
)

type Storage interface {
	CreateStudent(name string, email string, age int) (student types.Student, error error)
	GetStudentById(id int64) (student types.Student, error error)
	GetStudentList() ([]types.Student, error)
	DeleteStudent(id int64) (studentid types.Student, error error)
	RegisterStudent(context.Context, types.Credentials) (email string, error error)
	StudentSignin(context.Context, types.Login) (string, error)
}
