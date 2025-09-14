package storage

import "github.com/AbdulHaseebAhmad/go_project/internal/types"

type Storage interface {
	CreateStudent(name string, email string, age int) (id int64, error error)
	GetStudentById(id int64) (student types.Student, error error)
	GetStudentList() ([]types.Student, error)
}
