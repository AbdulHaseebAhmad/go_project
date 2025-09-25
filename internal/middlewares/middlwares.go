package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/AbdulHaseebAhmad/go_project/internal/Utils/response"
	"github.com/AbdulHaseebAhmad/go_project/internal/storage"
)

func Authorize(storage storage.Storage, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check if session token exist in db, if it does then user is authenticated
		sessionnId, err := r.Cookie("session_token")
		if err != nil {
			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(errors.New("missing session token")))
			return
		}
		fmt.Println("Got session_token cookie:", sessionnId.Value)

		csrfHeader := r.Header.Get("X-CSRF-Token")
		fmt.Println("Got CSRF header:", csrfHeader)
		sessionCookie := sessionnId.Value
		// fmt.Println(csrfHeader, sessionCookie)
		authorized := storage.AuthorizeStudent(r.Context(), sessionCookie, csrfHeader)
		fmt.Println(authorized)
		if !authorized {
			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(errors.New("unauthorized Access")))
			return
		}
		next.ServeHTTP(w, r)
	})

}
