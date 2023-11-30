package database

import (
	"database/sql"
	"time"
)

const (
	//TODO: basedata and details

	InsertProductMaxIdQuery = "INSERT INTO products (p_id, created_at) VALUES(?, ?);"
	GetProductMaxIdQuery    = "SELECT Max(p_id) FROM products;"

	InsertCustomerMaxIdQuery = "INSERT INTO customers (c_id, created_at) VALUES(?, ?);"
	GetCustomerMaxIdQuery    = "SELECT Max(c_id) FROM customers;"

	//TODO: invoice id is unique?
	InsertInvoiceMaxIdQuery = "INSERT INTO invoices (i_id, created_at) VALUES(?, ?);"
	GetInvoiceMaxIdQuery    = "SELECT Max(i_id) FROM invoices;"

	//TODO: invoice id is unique?
	InsertTreasuriesMaxIdQuery = "INSERT INTO treasuries (t_id, created_at) VALUES(?, ?);"
	GetTreasuriesMaxIdQuery    = "SELECT Max(t_id) FROM treasuries;"

	//TODO: invoice id is unique?
	InsertRevertedMaxIdQuery = "INSERT INTO reverted (r_id, created_at) VALUES(?, ?);"
	GetRevertedMaxIdQuery    = "SELECT Max(r_id) FROM reverted;"
)

type DatabaseInterface interface {
	InsertProduct(p_id int) error
	GetProduct() (int, error)

	InsertCustomer(c_id int) error
	GetCustomer() (int, error)

	InsertInvoice(i_id int) error
	GetInvoice() (int, error)

	InsertTreasuries(t_id int) error
	GetTreasuries() (int, error)

	InsertReverted(r_id int) error
	GetReverted() (int, error)
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
func (d *Database) GetCustomer() (int, error) {
	var maxId int
	err := d.sdb.QueryRow(GetCustomerMaxIdQuery).Scan(&maxId)
	if err != nil {
		return 0, err
	}

	return maxId, err
}

// GetInvoice implements DatabaseInterface
func (d *Database) GetInvoice() (i_id int, err error) {
	var maxId int
	err = d.sdb.QueryRow(GetInvoiceMaxIdQuery).Scan(&maxId)
	if err != nil {
		return 0, err
	}

	return maxId, err
}

// GetProduct implements DatabaseInterface
func (d *Database) GetProduct() (int, error) {
	var maxId int
	err := d.sdb.QueryRow(GetProductMaxIdQuery).Scan(&maxId)
	if err != nil {
		return 0, err
	}

	return maxId, err
}

// GetReverted implements DatabaseInterface
func (d *Database) GetReverted() (int, error) {
	var maxId int
	err := d.sdb.QueryRow(GetRevertedMaxIdQuery).Scan(&maxId)
	if err != nil {
		return 0, err
	}

	return maxId, err
}

// GetTreasuries implements DatabaseInterface
func (d *Database) GetTreasuries() (int, error) {
	var maxId int
	err := d.sdb.QueryRow(GetTreasuriesMaxIdQuery).Scan(&maxId)
	if err != nil {
		return 0, err
	}

	return maxId, err
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
	_, err := d.sdb.Exec(InsertInvoiceMaxIdQuery, i_id)
	if err != nil {
		return err
	}

	return nil
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
	_, err := d.sdb.Exec(InsertRevertedMaxIdQuery, r_id)
	if err != nil {
		return err
	}

	return nil
}

// InsertTreasuries implements DatabaseInterface
func (d *Database) InsertTreasuries(t_id int) error {
	_, err := d.sdb.Exec(InsertTreasuriesMaxIdQuery, t_id)
	if err != nil {
		return err
	}

	return nil
}
