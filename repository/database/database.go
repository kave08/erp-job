package database

import (
	"context"
	"database/sql"
)

const (
	//TODO: basedata and details

	InsertProductsMaxIdQuery = "INSERT INTO products (p_id, created_at) VALUES(?, ?);"
	GetProductsMaxIdQuery    = "SELECT Max(p_id) FROM products LIMIT 1;"

	InsertCustomersMaxIdQuery = "INSERT INTO customers (c_id, created_at) VALUES(?, ?);"
	GetCustomersMaxIdQuery    = "SELECT Max(c_id) FROM customers LIMIT 1;"

	//TODO: invoice id is unique?
	InsertInvoicesMaxIdQuery = "INSERT INTO invoices (i_id, created_at) VALUES(?, ?);"
	GetInvoicesMaxIdQuery    = "SELECT Max(i_id) FROM invoices LIMIT 1;"

	//TODO: invoice id is unique?
	InsertTreasuriesMaxIdQuery = "INSERT INTO treasuries (t_id, created_at) VALUES(?, ?);"
	GetTreasuriesMaxIdQuery    = "SELECT Max(t_id) FROM treasuries LIMIT 1;"

	//TODO: invoice id is unique?
	InsertRevertedMaxIdQuery = "INSERT INTO reverted (r_id, created_at) VALUES(?, ?);"
	GetRevertedMaxIdQuery    = "SELECT Max(r_id) FROM reverted LIMIT 1;"
)

type DatabaseInterface interface {
	InsertProduct(ctx context.Context, p_id int) error
	GetProduct(ctx context.Context, p_id int) (int, error)

	InsertCustomer(ctx context.Context, c_id int) error
	GetCustomer(ctx context.Context, c_id int) (int, error)

	InsertInvoice(ctx context.Context, i_id int) error
	GetInvoice(ctx context.Context, i_id int) (int, error)

	InsertTreasuries(ctx context.Context, t_id int) error
	GetTreasuries(ctx context.Context, t_id int) (int, error)

	InsertReverted(ctx context.Context, r_id int) error
	GetReverted(ctx context.Context, r_id int) (int, error)
}

type Database struct {
	sdb *sql.DB
}


func NewDatabase(sdb *sql.DB) DatabaseInterface {
	return &Database{
		sdb: sdb,
	}
}


// GetCustomer implements DatabaseInterface
func (d *Database) GetCustomer(ctx context.Context, c_id int) (int, error) {
	panic("unimplemented")
}

// GetInvoice implements DatabaseInterface
func (d *Database) GetInvoice(ctx context.Context, i_id int) (int, error) {
	panic("unimplemented")
}

// GetProduct implements DatabaseInterface
func (d *Database) GetProduct(ctx context.Context, p_id int) (int, error) {
	panic("unimplemented")
}

// GetReverted implements DatabaseInterface
func (d *Database) GetReverted(ctx context.Context, r_id int) (int, error) {
	panic("unimplemented")
}

// GetTreasuries implements DatabaseInterface
func (d *Database) GetTreasuries(ctx context.Context, t_id int) (int, error) {
	panic("unimplemented")
}

// InsertCustomer implements DatabaseInterface
func (d *Database) InsertCustomer(ctx context.Context, c_id int) error {
	panic("unimplemented")
}

// InsertInvoice implements DatabaseInterface
func (d *Database) InsertInvoice(ctx context.Context, i_id int) error {
	panic("unimplemented")
}

// InsertProduct implements DatabaseInterface
func (d *Database) InsertProduct(ctx context.Context, p_id int) error {
	panic("unimplemented")
}

// InsertReverted implements DatabaseInterface
func (d *Database) InsertReverted(ctx context.Context, r_id int) error {
	panic("unimplemented")
}

// InsertTreasuries implements DatabaseInterface
func (d *Database) InsertTreasuries(ctx context.Context, t_id int) error {
	panic("unimplemented")
}