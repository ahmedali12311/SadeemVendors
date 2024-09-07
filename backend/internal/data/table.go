package data

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// Table represents a table in the database.
type Table struct {
	ID              uuid.UUID  `db:"id" json:"id"`
	Name            string     `db:"name" json:"name"`
	VendorID        uuid.UUID  `db:"vendor_id" json:"vendor_id"`
	CustomerID      *uuid.UUID `db:"customer_id,omitempty" json:"customer_id,omitempty"`
	IsAvailable     bool       `db:"is_available" json:"is_available"`
	IsNeedsServices bool       `db:"is_needs_service" json:"is_needs_service"`
}

// TableDB wraps a sqlx.DB connection pool.
type TableDB struct {
	DB *sqlx.DB
}

// GetTables retrieves all tables from the database.
func (db *TableDB) GetTables(ctx context.Context) (*[]Table, error) {
	var tables []Table
	query, args, err := QB.Select(strings.Join(tableColumns, ",")).From("tables").ToSql()
	if err != nil {
		return nil, err
	}
	err = db.DB.SelectContext(ctx, &tables, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &tables, nil
}

// GetTable retrieves a table by its ID.
func (db *TableDB) GetTable(ctx context.Context, id uuid.UUID) (*Table, error) {
	var table Table
	query, args, err := QB.Select(strings.Join(tableColumns, ",")).From("tables").Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, err
	}
	err = db.DB.QueryRowxContext(ctx, query, args...).StructScan(&table)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &table, nil
}

// Insert inserts a new table into the database.
func (db *TableDB) Insert(ctx context.Context, table *Table) error {

	query, args, err := QB.
		Insert("tables").
		Columns("name", "vendor_id", "customer_id", "is_available", "is_needs_service").
		Values(table.Name, table.VendorID, table.CustomerID, table.IsAvailable, table.IsNeedsServices).
		Suffix(fmt.Sprintf("RETURNING %s", strings.Join(tableColumns, ", "))).
		ToSql()
	if err != nil {
		return err
	}

	err = db.DB.QueryRowxContext(ctx, query, args...).StructScan(table)
	if err != nil {
		return err
	}
	return nil
}

// DeleteTable deletes a table by its ID.
func (db *TableDB) DeleteTable(ctx context.Context, id uuid.UUID) (*Table, error) {
	var table Table
	query, args, err := QB.Delete("tables").Where(squirrel.Eq{"id": id}).Suffix(fmt.Sprintf("RETURNING %s", strings.Join(tableColumns, ", "))).ToSql()
	if err != nil {
		return nil, err
	}
	err = db.DB.QueryRowxContext(ctx, query, args...).StructScan(&table)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &table, nil
}

// Update updates an existing table in the database.
func (db *TableDB) Update(ctx context.Context, table *Table) error {
	query, args, err := QB.Update("tables").
		Set("name", table.Name).
		Set("is_available", table.IsAvailable).
		Set("is_needs_service", table.IsNeedsServices).
		Where(squirrel.Eq{"id": table.ID}).
		Suffix(fmt.Sprintf("RETURNING %s", strings.Join(tableColumns, ", "))).
		ToSql()
	if err != nil {
		return err
	}

	result, err := db.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("table not found")
	}
	return nil
}
