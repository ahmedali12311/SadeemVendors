package main

import (
	"net/http"

	"github.com/go-michi/michi"
)

func (app *application) Router() *michi.Router {
	r := michi.NewRouter()

	r.Route("/", func(sub *michi.Router) {
		r.Use(app.logRequest)
		r.Use(app.recoverPanic)
		r.Use(secureHeaders)
		r.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

		// User routes
		sub.HandleFunc("GET users", app.AuthMiddleware(app.requireAdmin(http.HandlerFunc(app.IndexUserHandler))))
		sub.HandleFunc("GET users/{id}", http.HandlerFunc(app.ShowUserHandler))
		sub.HandleFunc("PUT users/{id}", app.AuthMiddleware(app.AuthorizeUserUpdate(http.HandlerFunc(app.UpdateUserHandler))))
		sub.HandleFunc("DELETE users/{id}", app.AuthMiddleware(app.AuthorizeUserUpdate(http.HandlerFunc(app.DeleteUserHandler))))

		// Table routes (authentication required)
		/*

			sub.HandleFunc("POST tables", app.AuthMiddleware(http.HandlerFunc(app.StoreTableHandler)))
			sub.HandleFunc("PUT tables/{id}", app.AuthMiddleware(http.HandlerFunc(app.UpdateTableHandler)))
			sub.HandleFunc("DELETE tables/{id}", app.AuthMiddleware(http.HandlerFunc(app.DeleteTableHandler)))
			sub.HandleFunc("GET tables", app.)
		*/

		// Vendor routes
		sub.HandleFunc("GET vendors", app.IndexVendorHandler) // May be public
		sub.HandleFunc("GET vendors/{id}", app.ShowVendorHandler)
		sub.HandleFunc("POST vendors", app.AuthMiddleware(app.requireAdmin(http.HandlerFunc(app.CreateVendor))))
		sub.HandleFunc("PUT vendors/{id}", app.AuthMiddleware(app.requireVendorPermission(http.HandlerFunc(app.UpdateVendorHandler))))
		sub.HandleFunc("DELETE vendors/{id}", app.AuthMiddleware(app.requireAdmin(http.HandlerFunc(app.DeleteVendorHandler))))

		// Vendor Admin routes
		sub.HandleFunc("GET vendors/{id}/admins", app.AuthMiddleware(app.requireVendorPermission(http.HandlerFunc(app.GetVendorAdminsHandler))))
		sub.HandleFunc("POST vendors/{id}/admins", app.AuthMiddleware(app.requireVendorPermission(http.HandlerFunc(app.CreateVendorAdminHandler))))
		sub.HandleFunc("GET vendors/{id}/admins/{adminId}", app.AuthMiddleware(app.requireVendorPermission(http.HandlerFunc(app.GetVendorAdminHandler))))
		sub.HandleFunc("PUT vendors/{id}/admins/{adminId}", app.AuthMiddleware(app.requireVendorPermission(http.HandlerFunc(app.UpdateVendorAdminHandler))))
		sub.HandleFunc("DELETE vendors/{id}/admins/{adminId}", app.AuthMiddleware(app.requireVendorPermission(http.HandlerFunc(app.DeleteVendorAdminHandler))))
		sub.HandleFunc("GET uservendors/{id}", app.AuthMiddleware(app.requireVendorPermission(http.HandlerFunc(app.GetUserVendor))))

		// User role management (requires admin)
		sub.HandleFunc("POST userrole", app.AuthMiddleware(app.requireAdmin(http.HandlerFunc(app.GrantRoleHandler))))
		sub.HandleFunc("PUT userrole", app.AuthMiddleware(app.requireAdmin(http.HandlerFunc(app.UpdateUserRoleHandler))))
		sub.HandleFunc("DELETE userrole", app.AuthMiddleware(app.requireAdmin(http.HandlerFunc(app.RevokeRoleHandler))))
		sub.HandleFunc("GET userroles", app.AuthMiddleware(app.requireAdmin(http.HandlerFunc(app.IndexUserRoles))))
		sub.HandleFunc("GET userroles/{id}", app.AuthMiddleware(app.requireAdmin(http.HandlerFunc(app.ShowUserRoleHandler))))

		// Auth routes (public)
		sub.HandleFunc("POST signin", http.HandlerFunc(app.LoginHandler))
		sub.HandleFunc("POST signup", http.HandlerFunc(app.SignupHandler))
		sub.HandleFunc("GET me", app.AuthMiddleware(http.HandlerFunc(app.MeHandler)))

		// Custom route for getting a user's vendors
		sub.HandleFunc("GET users/{id}/vendors", app.AuthMiddleware(http.HandlerFunc(app.GetUserVendor)))
	})

	return r
}
