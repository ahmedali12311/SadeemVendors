package data

import (
	"database/sql"
	"fmt"
	"os"
	"project/utils"
	"project/utils/validator"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Vendor struct {
	ID          uuid.UUID `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Img         *string   `db:"img" json:"img"`
	Description string    `db:"description" json:"description"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

type VendorDB struct {
	db *sqlx.DB
}

func ValidatingVendor(v *validator.Validator, vendor *Vendor) {
	if vendor.Name != "" {
		v.Check(vendor.Name != "", "name", "Name can not be empty")
		v.Check(len(vendor.Name) >= 3, "name", "Name can't be less than 3 letters")
		v.Check(len(vendor.Name) <= 20, "name", "Name can't be larger than 20 letters ")
	}
	if vendor.Description != "" {
		v.Check(vendor.Description != "", "description", "description can't be empty ")
		v.Check(len(vendor.Description) >= 5, "description", "description can't be less than 5 letters ")
		v.Check(len(vendor.Description) <= 60, "description", "description can't be larger than 60 letters ")

	}

}
func (v *VendorDB) InsertVendor(vendor *Vendor) (*Vendor, error) {
	query, args, err := QB.Insert("vendors").Columns("name", "description", "img").
		Values(vendor.Name, vendor.Description, vendor.Img).
		Suffix(fmt.Sprintf("RETURNING %s", fmt.Sprint(strings.Join(vendors_columns, ",")))).ToSql()
	if err != nil {
		return nil, err
	}
	err = v.db.QueryRowx(query, args...).StructScan(vendor)
	if err != nil {
		return nil, fmt.Errorf("error while inserting vendor : %v", err)
	}
	return vendor, nil
}

func (v *VendorDB) DeleteVendor(id uuid.UUID) (*Vendor, error) {
	var vendor Vendor
	query, args, err := QB.Delete("vendors").Where(squirrel.Eq{"id": id}).Suffix(fmt.Sprintf("RETURNING %s", strings.Join(vendors_columns, ","))).ToSql()
	if err != nil {
		return nil, err
	}
	err = v.db.QueryRowx(query, args...).StructScan(&vendor)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	if vendor.Img != nil {
		imgfile := strings.TrimPrefix(*vendor.Img, Domain+"/")
		// Check if the file exists
		if _, err := os.Stat(imgfile); err == nil {
			// File exists, attempt to delete it
			err = utils.DeleteImageFile(imgfile)
			if err != nil {
				return nil, fmt.Errorf("failed to delete file %s: %w", imgfile, err)
			}
		} else if os.IsNotExist(err) {
			// File does not exist, log a warning but do not treat it as a fatal error
			fmt.Printf("Warning: image file %s does not exist\n", imgfile)
		} else {
			// Handle other potential errors from os.Stat
			return nil, fmt.Errorf("failed to check file %s: %w", imgfile, err)
		}
	}
	return &vendor, nil
}

func (v *VendorDB) UpdateVendor(vendor *Vendor) error {
	query, args, err := QB.Update("vendors").
		Set("name", vendor.Name).Set("img", vendor.Img).Set("description", vendor.Description).Set("updated_at", time.Now()).Where(squirrel.Eq{"id": vendor.ID}).Suffix(fmt.Sprintf("RETURNING %s", strings.Join(vendors_columns, ","))).ToSql()
	if err != nil {
		return err
	}
	err = v.db.QueryRowx(query, args...).StructScan(vendor)
	if err != nil {
		return fmt.Errorf("UpdateVendor: %v", err)
	}
	return nil
}
func (v *VendorDB) GetVendor(id uuid.UUID) (*Vendor, error) {
	query, args, err := QB.Select(strings.Join(vendors_columns, ",")).From("vendors").Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, err
	}
	var vendor Vendor
	err = v.db.QueryRowx(query, args...).StructScan(&vendor)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrRecordNotFound // or handle as not found
		}
		return &vendor, fmt.Errorf("GetVendor: %v", err)
	}
	return &vendor, nil
}
func (v *VendorDB) GetVendors(filters utils.Filters) (*[]Vendor, int, error) {
	var vendors []Vendor
	offset := (filters.Page - 1) * filters.PageSize

	queryBuilder := QB.Select(strings.Join(vendors_columns, ",")).From("vendors")

	// Apply search filter
	if filters.Search != "" {
		searchTerm := "%" + filters.Search + "%"
		queryBuilder = queryBuilder.Where("name ILIKE ?", searchTerm)
	}

	// Apply sorting
	switch filters.Sort {
	case "latest":
		queryBuilder = queryBuilder.OrderBy("created_at DESC")
	case "name_asc":
		queryBuilder = queryBuilder.OrderBy("name ASC")
	case "name_desc":
		queryBuilder = queryBuilder.OrderBy("name DESC")
	default:
		queryBuilder = queryBuilder.OrderBy("created_at DESC")

	}

	// Apply pagination
	query, args, err := queryBuilder.Limit(uint64(filters.PageSize)).Offset(uint64(offset)).ToSql()
	if err != nil {
		return nil, 0, err
	}

	err = v.db.Select(&vendors, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, 0, ErrRecordNotFound
		}
		return nil, 0, err
	}

	// Get the total number of vendors
	var totalVendorsCount int
	countQuery, _, err := QB.Select("COUNT(*)").From("vendors").ToSql()
	if err != nil {
		return nil, 0, err
	}
	err = v.db.Get(&totalVendorsCount, countQuery)
	if err != nil {
		return nil, 0, err
	}

	return &vendors, totalVendorsCount, nil
}

/*
	func (v *VendorDB) GetVendors() (*[]Vendor, error) {
		var vendors []Vendor
		query, args, err := QB.Select(strings.Join(vendors_columns, ",")).From("vendors").ToSql()
		if err != nil {
			return nil, err
		}
		err = v.db.Select(&vendors, query, args...)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, ErrRecordNotFound
			}

			return nil, err
		}

		return &vendors, nil
	}
*/
