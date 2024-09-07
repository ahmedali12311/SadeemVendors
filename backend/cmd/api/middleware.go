package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"project/internal/data"
	"project/utils"
	"project/utils/validator"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type contextKey string

const UserIDKey contextKey = "userID"
const UserRoleKey contextKey = "userRole"

func (app *application) AuthMiddleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			app.jwtErrorResponse(w, r, utils.ErrMissingToken)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			app.jwtErrorResponse(w, r, utils.ErrInvalidToken)
			return
		}

		tokenString := parts[1]
		token, err := utils.ValidateToken(tokenString)
		if err != nil {
			switch err.Error() {
			case "token contains an invalid number of segments":
				app.jwtErrorResponse(w, r, utils.ErrInvalidToken)
			default:
				app.jwtErrorResponse(w, r, utils.ErrInvalidClaims)
			}
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			app.jwtErrorResponse(w, r, utils.ErrInvalidClaims)
			return
		}

		if exp, ok := claims["exp"].(float64); ok {
			expTime := time.Unix(int64(exp), 0)
			if expTime.Before(time.Now()) {
				app.jwtErrorResponse(w, r, utils.ErrExpiredToken)
				return
			}
		} else {
			app.jwtErrorResponse(w, r, utils.ErrInvalidClaims)
			return
		}

		userID, okID := claims["id"].(string)
		userRole, okRole := claims["userRole"].(string)

		if !okID || !okRole {
			app.jwtErrorResponse(w, r, utils.ErrInvalidClaims)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		ctx = context.WithValue(ctx, UserRoleKey, userRole)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) requireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v := validator.New()
		if !app.isAdmin(v, r) {
			if len(v.Errors) > 0 {
				app.failedValidationResponse(w, r, v.Errors)
			} else {
				app.jwtErrorResponse(w, r, errors.New("you do not have permission to access this resource"))
			}
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Check if the user is an admin
func (app *application) isAdmin(v *validator.Validator, r *http.Request) bool {
	userIDStr, ok := r.Context().Value(UserIDKey).(string)
	if !ok {
		v.AddError("Token", "User ID is missing from context")
		return false
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		v.AddError("Token", "Invalid user ID format")
		return false
	}
	adminRoles, err := app.Model.User_roleDB.GetUserRole(userID)
	if err != nil {
		v.AddError("role", "Error retrieving user roles or insufficient permissions")
		return false
	}

	data.ValidatingUserRole(v, adminRoles)
	return v.Valid()
}

// Middleware to require vendor permissions
func (app *application) requireVendorPermission(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vendorIDStr := r.PathValue("id")

		if vendorIDStr == "" {
			app.badRequestResponse(w, r, errors.New("vendor ID is required"))
			return
		}

		vendorID, err := uuid.Parse(vendorIDStr)
		if err != nil {
			app.badRequestResponse(w, r, errors.New("invalid vendor ID format"))
			return
		}

		err = app.isVendorPermission(r, vendorID)
		if err != nil {
			if errors.Is(err, data.ErrRecordNotFound) {
				app.jwtErrorResponse(w, r, errors.New("you do not have permission to access this resource"))
			} else {
				app.jwtErrorResponse(w, r, errors.New("internal server error while checking vendor permissions"))
			}
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Check if the user has vendor permissions
func (app *application) isVendorPermission(r *http.Request, vendorID uuid.UUID) error {
	userIDStr, ok := r.Context().Value(UserIDKey).(string)
	if !ok {
		return errors.New("user ID is missing from context")
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return errors.New("invalid user ID format")
	}

	_, err = app.Model.VendorAdminDB.GetVendorAdmins(r.Context(), userID, vendorID)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			return errors.New("you do not have permission to access this resource")
		}
		return errors.New("error checking vendor permissions")
	}
	return nil
}

func (app *application) AuthorizeUserUpdate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(UserIDKey).(string)
		if !ok {
			app.errorResponse(w, r, http.StatusUnauthorized, "user ID is missing from context")
			return
		}

		userRole, ok := r.Context().Value(UserRoleKey).(string)
		if !ok {
			app.errorResponse(w, r, http.StatusUnauthorized, "user role is missing from context")
			return
		}

		// Check if the user is updating their own account
		if r.Method == http.MethodPut {
			if userID != r.Context().Value(UserIDKey) && userRole != "1" {
				app.errorResponse(w, r, http.StatusForbidden, "you do not have permission to update this user")
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		// CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000") // Allow only your frontend's origin
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method,
			r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
