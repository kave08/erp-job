package database

import (
	"database/sql"
	"time"
)

const (
	//TODO: basedata and details

	InsertProductMaxIdQuery = "INSERT INTO products (p_id, created_at) VALUES(?, ?);"
	GetProductMaxIdQuery    = "SELECT Max(p_id) FROM products LIMIT 1;"

	InsertCustomerMaxIdQuery = "INSERT INTO customers (c_id, created_at) VALUES(?, ?);"
	GetCustomerMaxIdQuery    = "SELECT Max(c_id) FROM customers LIMIT 1;"

	//TODO: invoice id is unique?
	InsertInvoiceMaxIdQuery = "INSERT INTO invoices (i_id, created_at) VALUES(?, ?);"
	GetInvoiceMaxIdQuery    = "SELECT Max(i_id) FROM invoices LIMIT 1;"

	//TODO: invoice id is unique?
	InsertTreasuriesMaxIdQuery = "INSERT INTO treasuries (t_id, created_at) VALUES(?, ?);"
	GetTreasuriesMaxIdQuery    = "SELECT Max(t_id) FROM treasuries LIMIT 1;"

	//TODO: invoice id is unique?
	InsertRevertedMaxIdQuery = "INSERT INTO reverted (r_id, created_at) VALUES(?, ?);"
	GetRevertedMaxIdQuery    = "SELECT Max(r_id) FROM reverted LIMIT 1;"
)

type DatabaseInterface interface {
	InsertProduct(p_id int) error
	GetProduct() (p_id int, err error)

	InsertCustomer(c_id int) error
	GetCustomer() (c_id int, err error)

	InsertInvoice(i_id int) error
	GetInvoice() (i_id int, err error)

	InsertTreasuries(t_id int) error
	GetTreasuries() (t_id int, err error)

	InsertReverted(r_id int) error
	GetReverted() (r_id int, err error)
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
func (d *Database) GetCustomer() (c_id int, err error) {
	err = d.sdb.QueryRow(GetCustomerMaxIdQuery).Err()
	if err != nil {
		return 0, err
	}
	return c_id, err
}

// GetInvoice implements DatabaseInterface
func (d *Database) GetInvoice() (i_id int, err error) {
	panic("unimplemented")
}

// GetProduct implements DatabaseInterface
func (d *Database) GetProduct() (p_id int, err error) {
	err = d.sdb.QueryRow(GetProductMaxIdQuery).Err()
	if err != nil {
		return 0, err
	}
	return p_id, err
}

// GetReverted implements DatabaseInterface
func (d *Database) GetReverted() (r_id int, err error) {
	panic("unimplemented")
}

// GetTreasuries implements DatabaseInterface
func (d *Database) GetTreasuries() (t_id int, err error) {
	panic("unimplemented")
}

// InsertCustomer implements DatabaseInterface
func (d *Database) InsertCustomer(c_id int) error {
	_, err := d.sdb.Exec(InsertCustomerMaxIdQuery, c_id, time.Now())
	if err != nil {
		return err
	}

	return nil
}

// InsertInvoice implements DatabaseInterface
func (d *Database) InsertInvoice(i_id int) error {
	panic("unimplemented")
}

// InsertProduct implements DatabaseInterface
func (d *Database) InsertProduct(p_id int) error {
	_, err := d.sdb.Exec(InsertProductMaxIdQuery, p_id)
	if err != nil {
		return err
	}

	return nil
}

// InsertReverted implements DatabaseInterface
func (d *Database) InsertReverted(r_id int) error {
	panic("unimplemented")
}

// InsertTreasuries implements DatabaseInterface
func (d *Database) InsertTreasuries(t_id int) error {
	panic("unimplemented")
}
