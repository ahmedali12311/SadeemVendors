package main

import (
<<<<<<< HEAD
	"encoding/json"
	"errors"
	"log"
=======
	"errors"
>>>>>>> d27b46be5e9dd1ccbadff4044dcca4c39a7d905c
	"net/http"
	"project/internal/data"
	"project/utils"
)

func (app *application) handleRetrievalError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, data.ErrRecordNotFound):
<<<<<<< HEAD
		app.errorResponse(w, r, http.StatusNotFound, "Invalid email or password!")
	case errors.Is(err, data.ErrDuplicatedKey):
		app.errorResponse(w, r, http.StatusConflict, utils.Envelope{"email": "Email already exists, try something else"})
	case errors.Is(err, data.ErrUserNotFound):
		app.errorResponse(w, r, http.StatusNotFound, data.ErrUserNotFound.Error())
=======
		app.errorResponse(w, r, http.StatusNotFound, "Resources could not be found")
	case errors.Is(err, data.ErrDuplicatedKey):
		app.errorResponse(w, r, http.StatusConflict, "Email already exists, try something else")
>>>>>>> d27b46be5e9dd1ccbadff4044dcca4c39a7d905c
	case errors.Is(err, data.ErrDuplicatedRole):
		app.errorResponse(w, r, http.StatusConflict, "User already have the role")
	case errors.Is(err, data.ErrHasRole):
		app.errorResponse(w, r, http.StatusConflict, "User already has a role")
	case errors.Is(err, data.ErrHasNoRoles):
		app.errorResponse(w, r, http.StatusConflict, data.ErrHasNoRoles.Error())
<<<<<<< HEAD
	case errors.Is(err, data.ErrForeignKeyViolation):
		app.errorResponse(w, r, http.StatusNotFound, "Vendor ID is incorrect")
	case errors.Is(err, data.ErrUserAlreadyhaveatable):
		app.errorResponse(w, r, http.StatusConflict, data.ErrUserAlreadyhaveatable.Error())
	case errors.Is(err, data.ErrUserHasNoTable):
		app.errorResponse(w, r, http.StatusNotFound, data.ErrUserHasNoTable.Error())
	case errors.Is(err, data.ErrItemAlreadyInserted):
		app.errorResponse(w, r, http.StatusConflict, data.ErrItemAlreadyInserted.Error())
	case errors.Is(err, data.ErrInvalidQuantity):
		app.errorResponse(w, r, http.StatusConflict, data.ErrInvalidQuantity.Error())
	case errors.Is(err, data.ErrRecordNotFoundOrders):
		app.errorResponse(w, r, http.StatusConflict, data.ErrRecordNotFoundOrders.Error())
=======
>>>>>>> d27b46be5e9dd1ccbadff4044dcca4c39a7d905c
	default:
		app.serverErrorResponse(w, r, err)
	}
}
func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	env := utils.Envelope{"error": message}
	err := utils.SendJSONResponse(w, status, env)
	if err != nil {
<<<<<<< HEAD
		app.logError(r, err)
=======
		app.logError(err)
>>>>>>> d27b46be5e9dd1ccbadff4044dcca4c39a7d905c
		w.WriteHeader(500)

	}
}

<<<<<<< HEAD
func (app *application) logError(r *http.Request, err error) {
	log.Printf("Error: %v, Method: %s, URL: %s", err, r.Method, r.URL.String())
}
func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)
	message := "the server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
=======
func (app *application) logError(err error) {
	app.log.Panic(err)

}
func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(err)
	message := "the server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)

>>>>>>> d27b46be5e9dd1ccbadff4044dcca4c39a7d905c
}
func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "resources not found"
	app.errorResponse(w, r, http.StatusNotFound, message)
<<<<<<< HEAD
=======

>>>>>>> d27b46be5e9dd1ccbadff4044dcca4c39a7d905c
}
func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}
func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
<<<<<<< HEAD
=======

>>>>>>> d27b46be5e9dd1ccbadff4044dcca4c39a7d905c
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}
func (app *application) jwtErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	var message string
	switch {
	case errors.Is(err, utils.ErrInvalidToken):
		message = "invalid token"
	case errors.Is(err, utils.ErrExpiredToken):
		message = "token has expired"
	case errors.Is(err, utils.ErrMissingToken):
		message = "missing authorization token"
	case errors.Is(err, utils.ErrInvalidClaims):
		message = "invalid token claims"
	default:
<<<<<<< HEAD
		app.errorResponse(w, r, http.StatusUnauthorized, "You don't have a premission")
=======
		app.serverErrorResponse(w, r, err)
>>>>>>> d27b46be5e9dd1ccbadff4044dcca4c39a7d905c
		return
	}
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}
<<<<<<< HEAD

type ErrorResponse struct {
	Error string `json:"error"`
}

func (app *application) ErrorHandlerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log the error
				app.log.Println("Recovered from panic:", err)

				// Send the error response
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)

				response := ErrorResponse{
					Error: "Internal Server Error",
				}
				json.NewEncoder(w).Encode(response)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

/*
func (app *application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	message := "rate limit exceeded"
	app.errorResponse(w, r, http.StatusTooManyRequests, message)
}
*/
=======
>>>>>>> d27b46be5e9dd1ccbadff4044dcca4c39a7d905c
