package main

import (
	"errors"
	"fmt"
	"net/http"
	"project/internal/data"
	"project/utils"
	"strconv"

	"github.com/google/uuid"
)

// IndexUserRoles handles the listing of user roles
func (app *application) IndexUserRoles(w http.ResponseWriter, r *http.Request) {
<<<<<<< HEAD
	userRoles, err := app.Model.UserRoleDB.GetUserRoles()
=======
	userRoles, err := app.Model.User_roleDB.GetUserRoles()
>>>>>>> d27b46be5e9dd1ccbadff4044dcca4c39a7d905c
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	utils.SendJSONResponse(w, http.StatusOK, utils.Envelope{"user_roles": userRoles})
}

// ShowUserRoleHandler handles the retrieval of a user role by ID
func (app *application) ShowUserRoleHandler(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid ID")
		return
	}

<<<<<<< HEAD
	userRoles, err := app.Model.UserRoleDB.GetUserRole(id)
=======
	userRoles, err := app.Model.User_roleDB.GetUserRole(id)
>>>>>>> d27b46be5e9dd1ccbadff4044dcca4c39a7d905c
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.errorResponse(w, r, http.StatusNotFound, fmt.Sprintf("User role with ID %v was not found", id))
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}
	utils.SendJSONResponse(w, http.StatusOK, utils.Envelope{"user_roles": userRoles})
}
<<<<<<< HEAD

// UpdateUserRoleHandler handles the updating of a user's role
func (app *application) GrantRole(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
=======
func (app *application) GrantRoleHandler(w http.ResponseWriter, r *http.Request) {
	var vendoradmin data.VendorAdmin
	id, err := uuid.Parse(r.FormValue("id"))
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid ID")
		return
	}

	_, err = app.Model.UserDB.GetUser(id)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.errorResponse(w, r, http.StatusNotFound, fmt.Sprintf("User with ID %v was not found", id))
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	roleID, err := strconv.Atoi(r.FormValue("user_role"))
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid role ID")
		return
	}

	if roleID == 2 {
		if r.FormValue("vendorID") == "" {
			app.errorResponse(w, r, http.StatusBadRequest, "Must enter the vendor ID")
			return
		}
		vendorID, err := uuid.Parse(r.FormValue("vendorID"))
		if err != nil {
			app.errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		vendoradmin.UserID = id
		vendoradmin.VendorID = vendorID
		if _, err := app.Model.VendorAdminDB.InsertVendorAdmin(r.Context(), vendoradmin); err != nil {
			app.handleRetrievalError(w, r, err)
			return
		}
	}

	user, err := app.Model.User_roleDB.GrantRole(id, roleID)
	if err != nil {
		if roleID == 2 {
			// Clean up vendor admin record if role assignment fails
			app.Model.VendorAdminDB.DeleteVendorAdmin(r.Context(), vendoradmin.UserID, vendoradmin.VendorID)
		}
		if errors.Is(err, data.ErrDuplicatedKey) {
			utils.SendJSONResponse(w, http.StatusConflict, utils.Envelope{"message": fmt.Sprintf("Role %d already granted to user %v", roleID, id)})
			return
		}
		app.handleRetrievalError(w, r, err)
		return
	}

	utils.SendJSONResponse(w, http.StatusCreated, utils.Envelope{"added user": user})
}

// UpdateUserRoleHandler handles the updating of a user's role
func (app *application) UpdateUserRoleHandler(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.FormValue("id"))
>>>>>>> d27b46be5e9dd1ccbadff4044dcca4c39a7d905c
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid ID")
		return
	}

	newRole, err := strconv.Atoi(r.FormValue("role"))
	if err != nil {
<<<<<<< HEAD
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid role ID")
=======
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid new role ID")
>>>>>>> d27b46be5e9dd1ccbadff4044dcca4c39a7d905c
		return
	}

	if newRole == 2 {
		if r.FormValue("vendorID") == "" {
			app.errorResponse(w, r, http.StatusBadRequest, "Must enter the vendor ID")
			return
		}
		vendorID, err := uuid.Parse(r.FormValue("vendorID"))
		if err != nil {
			app.errorResponse(w, r, http.StatusBadRequest, err.Error())
			return
		}

		vendoradmin := data.VendorAdmin{
			UserID:   id,
			VendorID: vendorID,
		}
		_, err = app.Model.VendorAdminDB.InsertVendorAdmin(r.Context(), vendoradmin)
		if err != nil {
			app.handleRetrievalError(w, r, err)
			return
		}
	}

<<<<<<< HEAD
	user, err := app.Model.UserRoleDB.UpdateRole(id, newRole)
=======
	user, err := app.Model.User_roleDB.UpdateRole(id, newRole)
>>>>>>> d27b46be5e9dd1ccbadff4044dcca4c39a7d905c
	if err != nil {
		if newRole == 2 {
			vendor, err := uuid.Parse(r.FormValue("vendorID"))
			if err != nil {
				app.errorResponse(w, r, http.StatusBadRequest, err)
				return
			}
			err = app.Model.VendorAdminDB.DeleteVendorAdmin(r.Context(), id, vendor)
			if err != nil {
				app.errorResponse(w, r, http.StatusBadRequest, err)
				return
			}
<<<<<<< HEAD

		}

=======
		}
>>>>>>> d27b46be5e9dd1ccbadff4044dcca4c39a7d905c
		app.handleRetrievalError(w, r, err)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, utils.Envelope{"Updated user role": user})
}

func (app *application) RevokeRoleHandler(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.FormValue("id"))
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid ID")
		return
	}

	role, err := strconv.Atoi(r.FormValue("user_role"))
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid role ID")
		return
	}

<<<<<<< HEAD
	err = app.Model.UserRoleDB.RevokeRole(id, role)
=======
	err = app.Model.User_roleDB.RevokeRole(id, role)
>>>>>>> d27b46be5e9dd1ccbadff4044dcca4c39a7d905c
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.errorResponse(w, r, http.StatusNotFound, fmt.Sprintf("user with ID %v's role already deleted", id))
			return
		}

		app.handleRetrievalError(w, r, err)
		return
	}
<<<<<<< HEAD
	if role == 2 {
		_, err = app.Model.VendorDB.GetUserVendors(r.Context(), id)
		if err != nil {
			if err == data.ErrRecordNotFound {
				user, err := app.Model.UserRoleDB.UpdateRole(id, 3)
				fmt.Print(user)
				if err != nil {
					app.handleRetrievalError(w, r, err)
					return
				}

			} else {
				app.handleRetrievalError(w, r, err)
				return
			}
		}
	}
	utils.SendJSONResponse(w, http.StatusOK, utils.Envelope{fmt.Sprintf("Deleted user %v 's role ", id): role})

=======

	utils.SendJSONResponse(w, http.StatusOK, utils.Envelope{fmt.Sprintf("Deleted user %v 's role ", id): role})
>>>>>>> d27b46be5e9dd1ccbadff4044dcca4c39a7d905c
}
