package storage

type Storage interface {
	CreateStudent(name string, email string, age int) (id int64, error error)
}
