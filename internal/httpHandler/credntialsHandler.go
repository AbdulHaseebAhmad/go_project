package credentials

import (
	"log/slog"
	"net/http"

	"github.com/AbdulHaseebAhmad/go_project/internal/storage"
)

// the first handler is for the route sign up
// we are returning the http.HandlerFunc because that is the type of the func we will use in the route in main
// we are passing storage as argument because we already have a connection pool ready in the main
// now in the main when we call this func we wil just pass it storage and then the request will be made to the db using that connection pool
func NewStudentRegister(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Request to Register Student")
		// check if the incoming fields are valid
		// if they are save them to db
		// return a success message
	}
}
