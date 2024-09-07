package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"project/internal/data"
	"project/utils"
	"project/utils/validator"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func (app *application) IndexVendorHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve query parameters for pagination
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")
	sort := r.URL.Query().Get("sort")
	search := r.URL.Query().Get("search")

	// Set default values for page and pageSize
	page := 1
	pageSize := 10

	// Parse page and pageSize query parameters
	if pageStr != "" {
		parsedPage, err := strconv.Atoi(pageStr)
		if err != nil || parsedPage < 1 {
			page = 1
		} else {
			page = parsedPage
		}
	}

	if pageSizeStr != "" {
		parsedPageSize, err := strconv.Atoi(pageSizeStr)
		if err != nil || parsedPageSize < 1 {
			pageSize = 10
		} else {
			pageSize = parsedPageSize
		}
	}

	// Set default sort to "latest" if it's not provided or invalid
	if sort == "" || !validator.In(sort, "latest", "name_asc", "name_desc") {
		sort = "latest"
	}

	filters := utils.Filters{
		Page:         page,
		PageSize:     pageSize,
		Sort:         sort,
		SortSafelist: []string{"latest", "name_asc", "name_desc"},
		Search:       search,
	}

	v := validator.New()
	utils.ValidateFilters(v, filters)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Fetch vendors with pagination
	vendors, totalCount, err := app.Model.VendorDB.GetVendors(filters)
	if err != nil {
		app.handleRetrievalError(w, r, err)
		return
	}

	// Prepare response
	response := utils.Envelope{
		"Vendors":    vendors,
		"TotalCount": totalCount,
		"Page":       page,
		"PageSize":   pageSize,
	}

	utils.SendJSONResponse(w, http.StatusOK, response)
}

func (app *application) ShowVendorHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	vendor, err := app.Model.VendorDB.GetVendor(id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			app.errorResponse(w, r, http.StatusNotFound, fmt.Sprintf("Vendor with ID %d was not found", id))
			return
		default:
			app.serverErrorResponse(w, r, err)
		}
	}
	utils.SendJSONResponse(w, http.StatusOK, utils.Envelope{"vendor": vendor})
}

func (app *application) CreateVendor(w http.ResponseWriter, r *http.Request) {
	var vendor = data.Vendor{
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
	}
	var newImage *string
	file, fileHeader, err := r.FormFile("img")
	if err != nil && err != http.ErrMissingFile {
		app.badRequestResponse(w, r, errors.New("invalid file"))
		return
	} else if err == nil {
		defer file.Close()
		imageName, err := utils.SaveImageFile(file, "vendors", fileHeader.Filename)
		if err != nil {
			app.errorResponse(w, r, http.StatusInternalServerError, "Error saving image")
			return

		}
		vendor.Img = &imageName
		newImage = &imageName
	}

	v := validator.New()
	data.ValidatingVendor(v, &vendor)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	ven, err := app.Model.VendorDB.InsertVendor(&vendor)
	if err != nil {
		utils.DeleteImageFile(*newImage)
		app.serverErrorResponse(w, r, err)
		return
	}

	utils.SendJSONResponse(w, http.StatusCreated, utils.Envelope{"vendor": ven})
}

func (app *application) UpdateVendorHandler(w http.ResponseWriter, r *http.Request) {
	var oldImg *string
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	vendor, err := app.Model.VendorDB.GetVendor(id)
	if err != nil {
		app.handleRetrievalError(w, r, err)
		return
	}

	if r.FormValue("name") != "" {
		vendor.Name = r.FormValue("name")
	}
	if r.FormValue("description") != "" {
		vendor.Description = r.FormValue("description")
	}
	if vendor.Img != nil {
		*vendor.Img = strings.TrimPrefix(*vendor.Img, data.Domain+"/")
		oldImg = vendor.Img

	}
	file, fileheader, err := r.FormFile("img")
	if vendor.Img != nil {
		*vendor.Img = strings.TrimPrefix(*vendor.Img, data.Domain+"/")
		oldImg = vendor.Img

	}
	if err != nil && err != http.ErrMissingFile {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid file")
		return
	} else if err == nil {
		file.Close()
		newimg, err := utils.SaveImageFile(file, "vendors", fileheader.Filename)
		if err != nil {
			app.errorResponse(w, r, http.StatusBadRequest, "error while saving image")
			return

		}
		vendor.Img = &newimg

	}
	v := validator.New()
	data.ValidatingVendor(v, vendor)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.Model.VendorDB.UpdateVendor(vendor)
	if err != nil {
		if vendor.Img != nil && oldImg != nil {
			utils.DeleteImageFile(*vendor.Img)
		}
		app.serverErrorResponse(w, r, err)
		return
	}
	if oldImg != nil && vendor.Img != nil && *oldImg != *vendor.Img {
		utils.DeleteImageFile(*oldImg)
	}
	utils.SendJSONResponse(w, http.StatusOK, utils.Envelope{"vendor": vendor})
}
func (app *application) DeleteVendorHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	iduu, err := uuid.Parse(id)
	if err != nil {
		app.notFoundResponse(w, r)
		return // Ensure to return after handling the error
	}
	vendor, err := app.Model.VendorDB.DeleteVendor(iduu)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.errorResponse(w, r, http.StatusNotFound, fmt.Sprintf("vendor with ID %s was not found", id))
			return
		default:
			app.serverErrorResponse(w, r, err)
			return
		}
	}
	utils.SendJSONResponse(w, http.StatusOK, utils.Envelope{"deleted vendor": vendor})
}
