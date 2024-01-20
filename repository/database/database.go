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

	InsertInvoiceToSaleFactorMaxIdQuery = "INSERT INTO invoice_to_sale_factor (i_id, created_at) VALUES(?, ?);"
	GetInvoiceToSaleFactorMaxIdQuery    = "SELECT Max(i_id) FROM invoice_to_sale_factor;"

	InsertInvoiceToGoodsMaxIdQuery = "INSERT INTO invoice_to_goods (i_id, created_at) VALUES(?, ?);"
	GetInvoiceToGoodsMaxIdQuery    = "SELECT Max(i_id) FROM invoice_to_goods;"

	InsertInvoiceToSaleOrderMaxIdQuery = "INSERT INTO invoice_to_sale_order (i_id, created_at) VALUES(?, ?);"
	GetInvoiceToSaleOrderMaxIdQuery    = "SELECT Max(i_id) FROM invoice_to_sale_order;"

	InsertInvoiceToSalePaymentMaxIdQuery = "INSERT INTO invoice_to_sale_payment (i_id, created_at) VALUES(?, ?);"
	GetInvoiceToSalePaymentMaxIdQuery    = "SELECT Max(i_id) FROM invoice_to_sale_payment;"

	InsertInvoiceToSalerSelectMaxIdQuery = "INSERT INTO invoice_to_saler_select (i_id, created_at) VALUES(?, ?);"
	GetInvoiceToSalerSelectMaxIdQuery    = "SELECT Max(i_id) FROM invoice_to_saler_select;"

	InsertInvoiceToSaleProformaMaxIdQuery = "INSERT INTO invoice_to_sale_proforma (i_id, created_at) VALUES(?, ?);"
	GetInvoiceToSaleProformaMaxIdQuery    = "SELECT Max(i_id) FROM invoice_to_sale_proforma;"

	InsertTreasuriesMaxIdQuery = "INSERT INTO treasuries (t_id, created_at) VALUES(?, ?);"
	GetTreasuriesMaxIdQuery    = "SELECT Max(t_id) FROM treasuries;"

	InsertRevertedMaxIdQuery = "INSERT INTO reverted (r_id, created_at) VALUES(?, ?);"
	GetRevertedMaxIdQuery    = "SELECT Max(r_id) FROM reverted;"
)

type DatabaseInterface interface {
	InsertProduct(p_id int) error
	GetProduct() (int, error)

	InsertCustomer(c_id int) error
	GetCustomer() (int, error)

	InsertInvoiceToSaleFactor(i_id int) error
	GetInvoiceToSaleFactor() (int, error)

	InsertInvoiceToGoods(i_id int) error
	GetInvoiceToGoods() (int, error)

	InsertInvoiceToSaleOrder(i_id int) error
	GetInvoiceToSaleOrder() (int, error)

	InsertInvoiceToSalePayment(i_id int) error
	GetInvoiceToSalePayment() (int, error)

	InsertInvoiceToSalerSelect(i_id int) error
	GetInvoiceToSalerSelect() (int, error)

	InsertInvoiceToSaleProforma(i_id int) error
	GetInvoiceToSaleProforma() (int, error)

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

func (d *Database) InsertInvoiceToSaleFactor(i_id int) error {
	_, err := d.sdb.Exec(InsertInvoiceToSaleFactorMaxIdQuery, i_id, time.Now())
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) GetInvoiceToSaleFactor() (int, error) {
	var maxId int
	err := d.sdb.QueryRow(GetInvoiceToSaleFactorMaxIdQuery).Scan(&maxId)
	if err != nil {
		return 0, err
	}

	return maxId, err
}

func (d *Database) InsertInvoiceToGoods(i_id int) error {
	_, err := d.sdb.Exec(InsertInvoiceToGoodsMaxIdQuery, i_id, time.Now())
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) GetInvoiceToGoods() (int, error) {
	var maxId int
	err := d.sdb.QueryRow(GetInvoiceToGoodsMaxIdQuery).Scan(&maxId)
	if err != nil {
		return 0, err
	}

	return maxId, err
}

func (d *Database) InsertInvoiceToSaleOrder(i_id int) error {
	_, err := d.sdb.Exec(InsertInvoiceToSaleOrderMaxIdQuery, i_id, time.Now())
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) GetInvoiceToSaleOrder() (int, error) {
	return 0, nil
}

func (d *Database) InsertInvoiceToSalePayment(i_id int) error {
	return nil
}

func (d *Database) GetInvoiceToSalePayment() (int, error) {
	return 0, nil
}

func (d *Database) InsertInvoiceToSalerSelect(i_id int) error {
	return nil
}

func (d *Database) GetInvoiceToSalerSelect() (int, error) {
	return 0, nil
}

func (d *Database) InsertInvoiceToSaleProforma(i_id int) error {
	return nil
}

func (d *Database) GetInvoiceToSaleProforma() (int, error) {
	return 0, nil
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
