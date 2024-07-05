package cart

import (
	"backend/server/utils"
	"context"
	"database/sql"
	"fmt"
)

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type repository struct {
	db DBTX
}

func NewRepository(db DBTX) Repository {
	return &repository{db}
}

func (r *repository) AddNewCart(ctx context.Context, cart *Cart) (*Cart, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		tx.Rollback()
		return &Cart{}, err
	}

	result, err := tx.ExecContext(ctx, "insert into cart (UserID,ProductID,Status,Total,Price,ProductName,Created_At,Last_Updated) values(?,?,?,?,?,?,?,?)", cart.UserID, cart.ProductID, cart.Status, cart.Total, cart.Price, cart.ProductName, cart.CreatedAt, cart.LastUpdated)
	if err != nil {
		tx.Rollback()
		return &Cart{}, err
	}
	cartID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return &Cart{}, err
	}
	cart.ID = int(cartID)

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return &Cart{}, err
	}

	return cart, nil

}

func (r *repository) GetCartByProductID(ctx context.Context, userID, productID int) error {
	c := &Cart{}
	err := r.db.QueryRowContext(ctx, "select id,userID,productID,status,total,price,productname,created_at,last_updated from cart where userID = ? and productID=?", userID, productID).Scan(&c.ID, &c.UserID, &c.ProductID, &c.Status, &c.Total, &c.Price, &c.ProductName, &c.CreatedAt, &c.LastUpdated)
	if err != nil {
		if err == sql.ErrNoRows {
			return utils.ErrNotFound
		}
		return err
	}
	return nil
}

func (r *repository) UpdateCart(ctx context.Context, cart *Cart) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	result, err := tx.ExecContext(ctx, "update cart set status = ? ,total = ? ,price = ?,productname = ?,last_updated =? where id= ? and userid=?  ", cart.Status, cart.Total, cart.Price, cart.ProductName, cart.LastUpdated, cart.ID, cart.UserID)
	if err != nil {
		tx.Rollback()
		return err
	}

	aff, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("retrieving error:%v", err)
	}

	if aff <= 0 {
		tx.Rollback()
		return fmt.Errorf("affected rows '%v' ,error message:%v", aff, err)
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r *repository) GetCartByID(ctx context.Context, ID int) (*Cart, error) {
	c := &Cart{}
	err := r.db.QueryRowContext(ctx, "select id,userID,productID,status,total,price,productname,created_at,last_updated from cart where id = ? ", ID).Scan(&c.ID, &c.UserID, &c.ProductID, &c.Status, &c.Total, &c.Price, &c.ProductName, &c.CreatedAt, &c.LastUpdated)
	if err != nil {
		if err == sql.ErrNoRows {
			return &Cart{}, utils.ErrNotFound
		}
		return &Cart{}, err
	}
	return c, nil
}

func (r *repository) GetCartByUserID(ctx context.Context, userID int) ([]*Cart, error) {
	cartSlice := []*Cart{}

	rows, err := r.db.QueryContext(ctx, "select id,userid,productid,status,total,price,productname,created_at,last_updated from cart where userid=?", userID)
	if err != nil {
		return []*Cart{}, err
	}
	defer rows.Close()
	for rows.Next() {
		c := &Cart{}
		err := rows.Scan(&c.ID, &c.UserID, &c.ProductID, &c.Status, &c.Total, &c.Price, &c.ProductName, &c.CreatedAt, &c.LastUpdated)
		if err != nil {
			return []*Cart{}, err
		}
		cartSlice = append(cartSlice, c)
	}

	return cartSlice, nil
}
