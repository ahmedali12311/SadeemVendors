<<<<<<< HEAD
package main

import (
	"errors"
=======
/*
package main

import (

	"errors"
	"fmt"
>>>>>>> d27b46be5e9dd1ccbadff4044dcca4c39a7d905c
	"net/http"
	"project/internal/data"
	"project/utils"

	"github.com/google/uuid"
<<<<<<< HEAD
)

// GetTablesHandler retrieves all tables from the database.
func (app *application) GetALLTablesHandler(w http.ResponseWriter, r *http.Request) {
	tables, err := app.Model.TableDB.GetTables(r.Context())
	if err != nil {
		app.handleRetrievalError(w, r, err)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, utils.Envelope{"tables": tables})
}

// GetTablesHandler retrieves all tables for a specific vendor from the database.
func (app *application) GetTablesHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve vendor ID from URL path
	vendorIDStr := r.PathValue("id")
	if vendorIDStr == "" {
		app.notFoundResponse(w, r)
		return
	}

	vendorID, err := uuid.Parse(vendorIDStr)
	if err != nil {
		app.badRequestResponse(w, r, errors.New("invalid vendor ID"))
		return
	}

	// Get tables for the specific vendor
	tables, err := app.Model.TableDB.GetVendorTables(r.Context(), vendorID)
	if err != nil {
		if err == data.ErrRecordNotFound {
			http.Error(w, "No tables found for the specified vendor", http.StatusNotFound)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, utils.Envelope{"tables": tables})
}

// GetTableHandler retrieves a single table by its ID and ensures it belongs to the vendor specified in the URL.
func (app *application) GetTableHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve table ID from URL path
	tableIDStr := r.PathValue("table_id")
	if tableIDStr == "" {
		app.notFoundResponse(w, r)
		return
	}

	tableID, err := uuid.Parse(tableIDStr)
	if err != nil {
		app.badRequestResponse(w, r, errors.New("invalid table ID"))
		return
	}

	// Get the table from the database
	table, err := app.Model.TableDB.GetTable(r.Context(), tableID)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, utils.Envelope{"table": table})
}

func (app *application) CreateTableHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve vendor ID from URL path
	vendorIDStr := r.PathValue("id")
	if vendorIDStr == "" {
		app.notFoundResponse(w, r)
		return
	}

	vendorID, err := uuid.Parse(vendorIDStr)
	if err != nil {
		app.badRequestResponse(w, r, errors.New("invalid vendor ID"))
		return
	}

	// Parse the form values
	err = r.ParseForm()
	if err != nil {
		app.badRequestResponse(w, r, errors.New("failed to parse form"))
		return
	}

	name := r.FormValue("name")
	isAvailableStr := r.FormValue("is_available")
	isNeedsServiceStr := r.FormValue("is_needs_service")

	// Parse boolean values
	isAvailable, err := utils.ParseBoolOrDefault(isAvailableStr, true)
	if err != nil {
		app.badRequestResponse(w, r, errors.New("invalid is_available value"))
		return
	}

	isNeedsService, err := utils.ParseBoolOrDefault(isNeedsServiceStr, false)
	if err != nil {
		app.badRequestResponse(w, r, errors.New("invalid is_needs_service value"))
		return
	}

	// Create a new table entity
	table := &data.Table{
		ID:              uuid.New(),
		Name:            name,
		VendorID:        vendorID,
		IsAvailable:     isAvailable,
		IsNeedsServices: isNeedsService,
	}
	tableCount, err := app.Model.TableDB.CountTablesForVendor(r.Context(), vendorID)
	if err != nil {

		app.serverErrorResponse(w, r, err)
		return
	}

	if tableCount >= 12 {
		app.badRequestResponse(w, r, errors.New("vendor has reached the maximum number of tables (12)"))
		return
	}

	// Insert the table into the database
	err = app.Model.TableDB.Insert(r.Context(), table)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	utils.SendJSONResponse(w, http.StatusCreated, utils.Envelope{"table": table})
}

// DeleteTableHandler handles the deletion of a table by its ID.
func (app *application) DeleteTableHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve vendor ID from URL path
	tableIDStr := r.PathValue("table_id")
	if tableIDStr == "" {
		app.notFoundResponse(w, r)
		return
	}

	tableID, err := uuid.Parse(tableIDStr)
	if err != nil {
		app.badRequestResponse(w, r, errors.New("invalid Table ID"))
		return
	}

	// Get the table from the database
	_, err = app.Model.TableDB.GetTable(r.Context(), tableID)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	table, err := app.Model.TableDB.DeleteTable(r.Context(), tableID)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, utils.Envelope{"message": "table deleted successfully", "table": table})
}

func (app *application) UpdateTableHandler(w http.ResponseWriter, r *http.Request) {
	tableIDStr := r.PathValue("table_id")
	if tableIDStr == "" {
		app.notFoundResponse(w, r)
		return
	}

	tableID, err := uuid.Parse(tableIDStr)
	if err != nil {
		app.badRequestResponse(w, r, errors.New("invalid vendor ID"))
		return
	}

	// Get the table from the database
	table, err := app.Model.TableDB.GetTable(r.Context(), tableID)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	if r.FormValue("name") != "" {
		table.Name = r.FormValue("name")
	}

	if err := app.Model.TableDB.Update(r.Context(), table); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if err := app.Model.OrderDB.DeleteUserOrders(*table.CustomerID); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	utils.SendJSONResponse(w, http.StatusOK, utils.Envelope{"table": table})
}
func (app *application) UpdateTableNeedsServiceHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve table ID and customer ID from URL path
	tableIDStr := r.PathValue("table_id")
	customerIDStr := r.Context().Value(UserIDKey).(string)
	if tableIDStr == "" || customerIDStr == "" {
		app.notFoundResponse(w, r)
		return
	}

	tableID, err := uuid.Parse(tableIDStr)
	if err != nil {
		app.badRequestResponse(w, r, errors.New("invalid table ID"))
		return
	}

	customerID, err := uuid.Parse(customerIDStr)
	if err != nil {
		app.badRequestResponse(w, r, errors.New("invalid customer ID"))
		return
	}

	// Get the table from the database
	table, err := app.Model.TableDB.GetTable(r.Context(), tableID)
	if err != nil {

		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
		} else {

			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if table.CustomerID != nil && *table.CustomerID != uuid.Nil && *table.CustomerID != customerID {
		app.errorResponse(w, r, http.StatusConflict, "The table is not available! Try again later")
		return
	}

	isNeedsServiceStr := r.FormValue("is_needs_service")
	isNeedsService, err := utils.ParseBoolOrDefault(isNeedsServiceStr, false)
	if err != nil {
		app.badRequestResponse(w, r, errors.New("invalid is_needs_service value"))
		return
	}

	table.IsNeedsServices = isNeedsService
	err = app.Model.TableDB.AssignCustomer(r.Context(), tableID, customerID)
	if err != nil {
		app.handleRetrievalError(w, r, err)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, utils.Envelope{"table": table})
}
func (app *application) FreeTableHandler(w http.ResponseWriter, r *http.Request) {
	tableIDStr := r.PathValue("table_id")
	customerIDStr := r.Context().Value(UserIDKey).(string) // Extract customer ID from context
	if tableIDStr == "" || customerIDStr == "" {
		app.notFoundResponse(w, r)
		return
	}

	// Parse the table ID and customer ID as UUIDs
	tableID, err := uuid.Parse(tableIDStr)
	if err != nil {
		app.badRequestResponse(w, r, errors.New("invalid table ID"))
		return
	}

	customerID, err := uuid.Parse(customerIDStr)
	if err != nil {
		app.badRequestResponse(w, r, errors.New("invalid customer ID"))
		return
	}

	// Get the table from the database
	table, err := app.Model.TableDB.GetTable(r.Context(), tableID)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	zeroUUID := uuid.UUID{}
	if table.CustomerID == nil || (table.CustomerID != nil && *table.CustomerID == zeroUUID) || *table.CustomerID != customerID {
		app.errorResponse(w, r, http.StatusConflict, "The table is not available! Try again later.")
		return
	}
	// Free the table
	if err := app.Model.TableDB.FreeTable(r.Context(), tableID); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, utils.Envelope{"message": "Table freed and user orders deleted successfully"})
}

func (app *application) FreeCustomerTableHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve table ID from URL path
	tableIDStr := r.PathValue("table_id")

	tableID, err := uuid.Parse(tableIDStr)
	if err != nil {
		app.badRequestResponse(w, r, errors.New("invalid table ID"))
		return
	}

	// Get the table from the database
	table, err := app.Model.TableDB.GetTable(r.Context(), tableID)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err := app.Model.TableDB.FreeTableVendor(r.Context(), tableID, table.VendorID); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, utils.Envelope{"message": "Table freed and user orders deleted successfully"})
}

func (app *application) GetCustomertable(w http.ResponseWriter, r *http.Request) {

	customerIDStr := r.Context().Value(UserIDKey).(string) // Extract customer ID from context
	if customerIDStr == "" {
		app.notFoundResponse(w, r)
		return
	}

	customerID, err := uuid.Parse(customerIDStr)
	if err != nil {
		app.badRequestResponse(w, r, errors.New("invalid customer ID"))
		return
	}

	table, err := app.Model.TableDB.GetCustomertable(r.Context(), customerID)
	if err != nil {
		if err != data.ErrUserHasNoTable {
			app.handleRetrievalError(w, r, err)
			return
		} else {
			utils.SendJSONResponse(w, http.StatusNotFound, utils.Envelope{"tables": "User has no tables"})
			return
		}

	}

	utils.SendJSONResponse(w, http.StatusOK, utils.Envelope{"tables": table})
}
=======

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
>>>>>>> d27b46be5e9dd1ccbadff4044dcca4c39a7d905c
