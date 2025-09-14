package student

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/AbdulHaseebAhmad/go_project/internal/Utils/response"
	"github.com/AbdulHaseebAhmad/go_project/internal/storage"
	"github.com/AbdulHaseebAhmad/go_project/internal/types"
	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Creating Student")
		var student types.Student                       // get the type from out internal Type
		err := json.NewDecoder(r.Body).Decode(&student) // decode body of request and put in the student Instance
		if errors.Is(err, io.EOF) {                     // check if the body is not empty
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		// now since the body is not empty  we check if there is any other error

		if err != nil { // this will now tell what other error there is if any
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// Now we have to validate the request
		// for validation we will use go play ground validator
		// we will use filled validation on the struct to use the validator
		validationErr := validator.New().Struct(student)
		if validationErr != nil {
			validateErrors := validationErr.(validator.ValidationErrors) // this is type casting
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrors))
			return
		}

		myStudent, err := storage.CreateStudent(student.Name, student.Email, student.Age)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}
		response.WriteJson(w, http.StatusCreated, myStudent)
	}
}

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id") // dynamic query parameter
		slog.Info("Getting a Student", slog.String("id: ", id))
		intId, conversionerr := strconv.ParseInt(id, 10, 64)
		if conversionerr != nil {
			// response.WriteJson(w, http.StatusBadRequest, response.GeneralError(conversionerr)) // my general error is not suitable to call it
			response.WriteJson(w, http.StatusBadRequest, conversionerr)
		}
		student, err := storage.GetStudentById(intId)
		if err != nil {
			// response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(conversionerr)) // my general error is not suitable to call it
			response.WriteJson(w, http.StatusInternalServerError, conversionerr)

		}
		response.WriteJson(w, http.StatusOK, student)
	}

}

func GetStudentList(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Getting a Student")
		stuents, err := storage.GetStudentList()
		if err != nil {
			// response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(conversionerr)) // my general error is not suitable to call it
			response.WriteJson(w, http.StatusInternalServerError, err)

		}
		response.WriteJson(w, http.StatusOK, stuents)
	}

}

func DeleteStudent(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Deleting Student")
		id := r.PathValue("id") // dynamic query parameter
		intId, conversionerr := strconv.ParseInt(id, 10, 64)
		if conversionerr != nil {
			response.WriteJson(w, http.StatusBadRequest, conversionerr)
		}
		student, err := storage.DeleteStudent(intId)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, conversionerr)
		}
		slog.Info("Deleting Student Success")
		response.WriteJson(w, http.StatusOK, student)
	}
}
