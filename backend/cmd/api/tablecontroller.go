/*
package main

import (

	"errors"
	"fmt"
	"net/http"
	"project/internal/data"
	"project/utils"

	"github.com/google/uuid"

)

	func (app *application) IndexTableHandler(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tables, err := app.Model.TableDB.GetTables(ctx)
		if err != nil {
			app.notFoundResponse(w, r)
			return
		}

		utils.SendJSONResponse(w, http.StatusOK, utils.Envelope{"tables": tables})
	}

// ShowTableHandler handles showing a specific table by its ID.

	func (app *application) ShowTableHandler(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := r.PathValue("id")
		tableID, err := uuid.Parse(id)
		if err != nil {
			app.notFoundResponse(w, r)
			return
		}
		table, err := app.Model.TableDB.GetTable(ctx, tableID)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.errorResponse(w, r, http.StatusNotFound, fmt.Sprintf("Table with %v was not found", tableID))
				return
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}
		utils.SendJSONResponse(w, http.StatusOK, utils.Envelope{"table": table})
	}

// StoreTableHandler handles creating a new table.

	func (app *application) StoreTableHandler(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, err := uuid.Parse(ctx.Value(UserIDKey).(string))
		if err != nil {
			app.badRequestResponse(w, r, errors.New("invalid user ID"))
			return
		}
		if r.FormValue("vendor_id") == "" {
			app.badRequestResponse(w, r, errors.New(" must enter a Vendor ID"))
			return
		}
		table := &data.Table{
			Name:            r.FormValue("name"),
			VendorID:        uuid.Must(uuid.Parse(r.FormValue("vendor_id"))),
			IsAvailable:     r.FormValue("is_available") == "true",
			IsNeedsServices: r.FormValue("is_needs_service") == "true",
		}

		if customerID := r.FormValue("customer_id"); customerID != "" {
			custUUID, err := uuid.Parse(customerID)
			if err != nil {
				app.badRequestResponse(w, r, errors.New("invalid customer_id"))
				return
			}
			table.CustomerID = &custUUID
		}

		if VendorID := r.FormValue("vendor_id"); VendorID != "" {
			vendorUUID, err := uuid.Parse(VendorID)
			if err != nil {
				app.badRequestResponse(w, r, errors.New("invalid vendor id"))
				return
			}
			table.VendorID = vendorUUID
		} else {
			app.badRequestResponse(w, r, errors.New("vendor cannot be empty"))
			return
		}

		if r.FormValue("customer_id") != "" {
			_, err := app.Model.UserDB.GetUser(*table.CustomerID)
			if err != nil {
				if errors.Is(err, data.ErrRecordNotFound) {
					app.errorResponse(w, r, http.StatusNotFound, fmt.Sprintf("User with ID %v was not found", table.CustomerID))
					return
				}
				app.serverErrorResponse(w, r, err)
				return
			}
		}

		if err != nil {
			if errors.Is(err, data.ErrRecordNotFound) {
				app.errorResponse(w, r, http.StatusNotFound, fmt.Sprintf("Vendor with ID %v was not found", table.VendorID))
				return
			}
			app.serverErrorResponse(w, r, err)
			return
		}

		_, err = app.Model.VendorAdminDB.GetVendorAdmin(ctx, userID)
		if err != nil {
			if errors.Is(err, data.ErrRecordNotFound) {
				app.errorResponse(w, r, http.StatusForbidden, "You do not have permission to insert a table for this vendor")
				return
			}
			app.serverErrorResponse(w, r, err)
			return
		}

		err = app.Model.TableDB.Insert(ctx, table)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		utils.SendJSONResponse(w, http.StatusCreated, utils.Envelope{"table": table})
	}

// UpdateTableStatusHandler handles updating the status of a table.

	func (app *application) Update(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		userID, err := uuid.Parse(ctx.Value(UserIDKey).(string))

		if err != nil {
			app.notFoundResponse(w, r)
			return
		}
		id := r.PathValue("id")
		tableID, err := uuid.Parse(id)
		if err != nil {
			app.notFoundResponse(w, r)
			return
		}

		table, err := app.Model.TableDB.GetTable(ctx, tableID)
		if err != nil {
			if errors.Is(err, data.ErrRecordNotFound) {
				app.errorResponse(w, r, http.StatusNotFound, "Table not found")
				return
			}
			app.serverErrorResponse(w, r, err)
			return
		}
		vendor_id := table.VendorID
		vendor, err := app.Model.VendorDB.GetVendor(vendor_id)
		if err != nil {
			if errors.Is(err, data.ErrRecordNotFound) {
				app.errorResponse(w, r, http.StatusNotFound, "Table not found")
				return
			}
			app.serverErrorResponse(w, r, err)
			return

		}
		if table.VendorID != vendor.ID {
			app.errorResponse(w, r, http.StatusNotFound, "wrong vendor to be updated")
			return
		}

		// Update status fields based on request parameters
		if r.FormValue("is_available") != "" {
			table.IsAvailable = r.FormValue("is_available") == "true"
		}

		if r.FormValue("is_needs_service") != "" {
			table.IsNeedsServices = r.FormValue("is_needs_service") == "true"
		}

		_, err = app.Model.VendorAdminDB.GetVendorAdmin(ctx, userID)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.errorResponse(w, r, http.StatusNotFound, fmt.Sprintf("you dont have premission to update table  %s", tableID))
				return
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}
		err = app.Model.TableDB.Update(ctx, table)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		utils.SendJSONResponse(w, http.StatusOK, utils.Envelope{"table": table})
	}

	func (app *application) DeleteTableHandler(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, err := uuid.Parse(ctx.Value(UserIDKey).(string))
		if err != nil {
			app.notFoundResponse(w, r)
			return
		}

		id := r.PathValue("id")
		tableID, err := uuid.Parse(id)
		if err != nil {
			app.notFoundResponse(w, r)
			return
		}

		// Fetch the table to check the vendor
		table, err := app.Model.TableDB.GetTable(ctx, tableID)
		if err != nil {
			if errors.Is(err, data.ErrRecordNotFound) {
				app.errorResponse(w, r, http.StatusNotFound, fmt.Sprintf("Table with %v was not found", tableID))
				return
			}
			app.serverErrorResponse(w, r, err)
			return
		}

		_, err = app.Model.VendorAdminDB.GetVendorAdmin(ctx, userID)
		if err != nil {
			if errors.Is(err, data.ErrRecordNotFound) {
				app.errorResponse(w, r, http.StatusForbidden, "You do not have permission to delete this table")
				return
			}
			app.serverErrorResponse(w, r, err)
			return
		}

		_, err = app.Model.TableDB.DeleteTable(ctx, tableID)
		if err != nil {
			if errors.Is(err, data.ErrRecordNotFound) {
				app.errorResponse(w, r, http.StatusNotFound, fmt.Sprintf("Table with %v was not found", tableID))
				return
			}
			app.serverErrorResponse(w, r, err)
			return
		}

		utils.SendJSONResponse(w, http.StatusOK, utils.Envelope{"deleted_table": table})
	}
*/
package main
