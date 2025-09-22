package credentialsss

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/AbdulHaseebAhmad/go_project/internal/Utils/response"
	"github.com/AbdulHaseebAhmad/go_project/internal/storage"
	"github.com/AbdulHaseebAhmad/go_project/internal/types"
	"github.com/go-playground/validator/v10"
)

func NewStudentRegister(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Request to Register Student")
		if r.Header.Get("Content-Type") != "application/json" {
			response.WriteJson(w, http.StatusUnsupportedMediaType, map[string]string{
				"error": "Content-Type must be application/json",
			})
			return
		}
		// cast incoming data to struct
		var credentials types.Credentials
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&credentials)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		// check if no missing field is present
		validationErr := validator.New().Struct(credentials)
		if validationErr != nil {
			validateErrors := validationErr.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(validateErrors))
			return
		}

		email, registrationErr := storage.RegisterStudent(r.Context(), credentials)
		if registrationErr != nil {
			response.WriteJson(w, http.StatusConflict, response.GeneralError(registrationErr))
			return
		}
		response.WriteJson(w, http.StatusCreated, email)

	}
}

func StudentSignin(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Request to Login Student")
		if r.Header.Get("Content-Type") != "application/json" {
			response.WriteJson(w, http.StatusUnsupportedMediaType, map[string]string{
				"error": "Content-Type must be application/json",
			})
			return
		}
		var loginData types.Login
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&loginData)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		sessiontoken, loginerr := storage.StudentSignin(r.Context(), loginData)
		if loginerr != nil {
			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(loginerr))
			return
		}
		//if token is generated then we attach it to a cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "Session_token",
			Value:    sessiontoken,
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
			Secure:   true, // ✅
			SameSite: http.SameSiteStrictMode,
		})
		response.WriteJson(w, http.StatusAccepted, sessiontoken)
	}
}
