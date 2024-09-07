package data

import (
	"errors"
	"fmt"
	"os"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrDuplicatedKey  = errors.New("user already have the value")
	ErrDuplicatedRole = errors.New("user Already have the role")
	ErrHasRole        = errors.New("user Already has a role")
	ErrHasNoRoles     = errors.New("user Has no roles")

	ErrUserNotFound = errors.New("User Not Found")

	QB     = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	Domain = os.Getenv("DOMAIN")

	user_columns = []string{
		"id",
		"name",
		"email",
		"password",
		"phone",
		"created_at",
		"updated_at",
		fmt.Sprintf("CASE WHEN NULLIF(img, '') IS NOT NULL THEN FORMAT('%s/%%s', img) ELSE NULL END AS img", Domain),
	}
	vendors_columns = []string{
		"id",
		"name",
		"description",
		"created_at",
		"updated_at",
		fmt.Sprintf("CASE WHEN NULLIF(img, '') IS NOT NULL THEN FORMAT('%s/%%s', img) ELSE NULL END AS img", Domain),
	}
	user_roles = []string{
		"user_id",
		"role_id",
	}
	tableColumns = []string{"id", "name", "vendor_id", "customer_id", "is_available", "is_needs_service"}
)

type Model struct {
	UserDB        userDB
	TableDB       TableDB
	VendorDB      VendorDB
	User_roleDB   user_roleDB
	VendorAdminDB VendorAdminDB
}

func NewModels(db *sqlx.DB) Model {
	return Model{
		UserDB:        userDB{db},
		TableDB:       TableDB{db},
		VendorDB:      VendorDB{db},
		User_roleDB:   user_roleDB{db},
		VendorAdminDB: VendorAdminDB{db},
	}
}
