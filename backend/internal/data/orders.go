package data

import (
	"database/sql"
	"fmt"
	"project/utils/validator"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type OrderDetails struct {
	ID             uuid.UUID       `db:"id" json:"id"`
	TotalOrderCost float64         `db:"total_order_cost" json:"total_order_cost"`
	VendorName     string          `db:"vendor_name" json:"vendor_name"`
	VendorID       uuid.UUID       `db:"vendor_id" json:"-"`
	UserName       string          `db:"user_name" json:"user_name"`
	CustomerID     string          `db:"customer_id" json:"customer_id"`
	ItemNames      pq.StringArray  `db:"item_names" json:"item_names"`
	ItemPrices     pq.Float64Array `db:"item_prices" json:"item_prices"`
	ItemQuantities pq.Int64Array   `db:"item_quantities" json:"item_quantities"`
	Status         string          `db:"status" json:"status"`
	TableID        uuid.UUID       `db:"table_id" json:"-"`
	TableName      string          `db:"table_name" json:"table_name"`
	Description    string          `db:"description" json:"description"` // New field
	CustomerPhone  string          `db:"user_phone" json:"CustomerPhone"`
	CreatedAt      time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time       `db:"updated_at" json:"updated_at"`
}
type Order struct {
	ID             uuid.UUID `db:"id" json:"id"`
	TotalOrderCost float64   `db:"total_order_cost" json:"total_order_cost"`
	CustomerID     uuid.UUID `db:"customer_id" json:"customer_id"`
	VendorID       uuid.UUID `db:"vendor_id" json:"vendor_id"`
	Status         string    `db:"status" json:"status"`
	Description    string    `db:"description" json:"description"` // New field
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}

type OrderDB struct {
	db *sqlx.DB
}

func ValidatingOrder(v *validator.Validator, order *Order, fields ...string) {
	for _, field := range fields {
		switch field {
		case "customer_id":
			v.Check(order.CustomerID != uuid.Nil, "customer_id", "Customer ID is required")
		case "vendor_id":
			v.Check(order.VendorID != uuid.Nil, "vendor_id", "Vendor ID is required")
		case "status":
			v.Check(order.Status != "", "status", "Order status is required")
			v.Check(order.Status == "completed" || order.Status == "preparing", "status", "Order must be of 2 statuses")
		case "total_order_cost":
			v.Check(order.TotalOrderCost >= 0, "total_order_cost", "Total order cost must be a non-negative number")

		}
	}
}
func (o *OrderDB) GetOrders(customerID uuid.UUID, tableID uuid.UUID) ([]OrderDetails, error) {
	// Prepare the SQL query with correct column names matching the struct tags
	query, args, err := QB.Select(
		allColumns...,
	).
		From("orders o").
		Join("vendors v ON o.vendor_id = v.id").
		Join("users c ON o.customer_id = c.id").
		Join("order_items oi ON o.id = oi.order_id").
		Join("items i ON oi.item_id = i.id").
		Join("tables t ON o.customer_id = t.customer_id").
		Where(squirrel.Eq{"o.customer_id": customerID, "t.id": tableID}).
		GroupBy("o.id, v.name, o.vendor_id, c.name, o.status, t.id, t.name, c.phone").
		ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := o.db.Queryx(query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrRecordNotFoundOrders
		}
		return nil, err
	}
	defer rows.Close()

	var orders []OrderDetails
	for rows.Next() {
		var order OrderDetails
		err := rows.StructScan(&order)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (o *OrderDB) InsertOrder(order *Order) error {
	query, args, err := QB.Insert("orders").
		Columns(strings.Join(ordersColumns, ",")).
		Values(order.ID, order.TotalOrderCost, order.CustomerID, order.Description, order.VendorID, order.Status).
		ToSql()
	if err != nil {
		return err
	}
	_, err = o.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("error while inserting order: %v", err)
	}
	return nil
}

func (o *OrderDB) DeleteOrder(orderID uuid.UUID) error {
	query, args, err := QB.Delete("orders").
		Where(squirrel.Eq{"id": orderID}).
		ToSql()
	if err != nil {
		return err
	}
	_, err = o.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("error while deleting order: %v", err)
	}
	return nil
}
func (o *OrderDB) GetVendorOrders(vendorID uuid.UUID) ([]OrderDetails, error) {
	query, args, err := QB.Select(
		allColumns...,
	).
		From("orders o").
		Join("vendors v ON o.vendor_id = v.id").
		Join("users c ON o.customer_id = c.id").
		Join("order_items oi ON o.id = oi.order_id").
		Join("items i ON oi.item_id = i.id").
		Join("tables t ON o.customer_id = t.customer_id").
		GroupBy("o.id, v.name, o.vendor_id, c.name, o.status, t.id, t.name, c.phone").
		ToSql()

	if err != nil {

		return nil, err
	}

	rows, err := o.db.Queryx(query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrRecordNotFoundOrders
		}
		return nil, err
	}
	defer rows.Close()

	var orders []OrderDetails
	for rows.Next() {
		var order OrderDetails
		err := rows.StructScan(&order)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

// UpdateOrder updates the order status to "completed".
func (o *OrderDB) UpdateOrder(orderID uuid.UUID, status string) error {
	if status != "completed" {
		return fmt.Errorf("only 'completed' status is allowed")
	}

	query, args, err := QB.Update("orders").
		Set("status", status).
		Where(squirrel.Eq{"id": orderID}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = o.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("error while updating order: %v", err)
	}
	return nil
}

// DeleteCompletedOrder deletes the order if its status is "completed".
func (o *OrderDB) DeleteCompletedOrder(orderID uuid.UUID) error {
	// Get the order by ID
	order, err := o.GetOrder(orderID)
	if err != nil {
		return fmt.Errorf("error while getting order: %v", err)
	}

	// Check if the order is completed
	if order.Status != "completed" {
		return fmt.Errorf("order is not completed")
	}

	// Delete the order
	query, args, err := QB.Delete("orders").
		Where(squirrel.Eq{"id": orderID}).
		ToSql()
	if err != nil {
		return err
	}
	_, err = o.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("error while deleting order: %v", err)
	}
	return nil
}

func (o *OrderDB) GetOrder(orderID uuid.UUID) (*Order, error) {
	query, args, err := QB.Select(strings.Join(ordersColumns, ",")).
		From("orders").
		Where(squirrel.Eq{"id": orderID}).
		ToSql()
	if err != nil {
		return nil, err
	}
	row, err := o.db.Queryx(query, args...)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	if !row.Next() {
		return nil, fmt.Errorf("order not found")
	}

	var order Order
	err = row.StructScan(&order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}
func (o *OrderDB) DeleteUserOrders(customerID uuid.UUID) error {
	query, args, err := QB.Delete("orders").
		Where(squirrel.Eq{"customer_id": customerID}).
		ToSql()
	if err != nil {
		return err
	}
	_, err = o.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("error while deleting user orders: %v", err)
	}
	return nil
}
