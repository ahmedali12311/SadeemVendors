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
	ErrRecordNotFound        = errors.New("record not found")
	ErrDuplicatedKey         = errors.New("user already have the value")
	ErrDuplicatedRole        = errors.New("user Already have the role")
	ErrHasRole               = errors.New("user Already has a role")
	ErrHasNoRoles            = errors.New("user Has no roles")
	ErrForeignKeyViolation   = errors.New("foreign key constraint violation")
	ErrUserNotFound          = errors.New("user Not Found")
	ErrUserAlreadyhaveatable = errors.New("user already have a table")
	ErrUserHasNoTable        = errors.New("user has no table")
	ErrItemAlreadyInserted   = errors.New("item already inserted! ")
	ErrInvalidQuantity       = errors.New("requested quantity is not available")
	ErrRecordNotFoundOrders  = errors.New("no orders available! ")
	ErrDescriptionMissing    = errors.New("description is required") // New error

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
		"subscription_end",
		"subscription_days",
		"is_visible",
		"created_at",
		"updated_at",
		fmt.Sprintf("CASE WHEN NULLIF(img, '') IS NOT NULL THEN FORMAT('%s/%%s', img) ELSE NULL END AS img", Domain),
	}
	user_roles = []string{
		"user_id",
		"role_id",
	}
	tableColumns     = []string{"id", "name", "vendor_id", "customer_id", "is_available", "is_needs_service"}
	cartItemsColumns = []string{
		"cart_id", "item_id", "quantity",
	}

	cartsColumns = []string{
		"id", "total_price", "quantity", "vendor_id", "description", "created_at", "updated_at",
	}

	orderItemsColumns = []string{
		"id", "order_id", "item_id", "quantity", "price",
	}

	ordersColumns = []string{
		"id", "total_order_cost", "customer_id", "vendor_id", "status", "description", "created_at", "updated_at",
	}

	itemsColumns = []string{
		"id",
		"vendor_id",
		"name",
		"price",
		"quantity",
		"discount",
		"discount_expiry",
		"created_at",
		"updated_at",
		fmt.Sprintf("CASE WHEN NULLIF(img, '') IS NOT NULL THEN FORMAT('%s/%%s', img) ELSE NULL END AS img", Domain),
	}
	orders_ColumnsJOIN = []string{
		"o.id",
		"o.total_order_cost",
		"o.customer_id",
		"o.vendor_id",
		"o.status",
		"o.description",
		"o.created_at",
		"o.updated_at",
	}

	vendorColumnsJOIN = []string{
		"v.name AS vendor_name",
	}

	userColumnsJOIN = []string{
		"c.name AS user_name",
		"c.phone AS user_phone", // add this
	}

	table_ColumnsJOIN = []string{
		"t.id AS table_id",
		"t.name AS table_name",
	}

	itemColumnsJOIN = []string{
		"ARRAY_AGG(i.name) AS item_names",
		"ARRAY_AGG(oi.price) AS item_prices",
		"ARRAY_AGG(oi.quantity) AS item_quantities",
	}

	allColumns = append(append(append(append(orders_ColumnsJOIN, vendorColumnsJOIN...), userColumnsJOIN...), table_ColumnsJOIN...), itemColumnsJOIN...)
)

type Model struct {
	UserDB        UserDB
	TableDB       TableDB
	VendorDB      VendorDB
	UserRoleDB    UserRoleDB
	VendorAdminDB VendorAdminDB
	CartItemDB    CartItemDB
	CartDB        CartDB
	OrderItemDB   OrderItemDB
	OrderDB       OrderDB
	ItemDB        ItemDB
	TransactionDB Transaction
}

func NewModels(db *sqlx.DB) Model {
	tx, err := db.Beginx()
	if err != nil {
		return Model{}
	}
	return Model{
		UserDB:        UserDB{db},
		TableDB:       TableDB{db},
		VendorDB:      VendorDB{db},
		UserRoleDB:    UserRoleDB{db},
		VendorAdminDB: VendorAdminDB{db},
		CartItemDB:    CartItemDB{db},
		CartDB:        CartDB{db},
		OrderItemDB:   OrderItemDB{db},
		OrderDB:       OrderDB{db},
		ItemDB:        ItemDB{db},
		TransactionDB: Transaction{tx},
	}
}
