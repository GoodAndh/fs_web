package product

import (
	"backend/server/utils"
	"context"
	"database/sql"
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

	return &repository{db: db}
}

func (r *repository) CreateProduct(ctx context.Context, product *Product) (int, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	result, err := tx.ExecContext(ctx, "insert into product (userid,name,description,price,stock,created_at,last_updated) values(?,?,?,?,?,?,?)", product.UserID, product.Name, product.Description, product.Price, product.Stock, product.CreatedAt, product.LastUpdated)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return int(id), nil
}

func (r *repository) GetAllProduct(ctx context.Context) ([]*Product, error) {
	rows, err := r.db.QueryContext(ctx, "select id,userID,name,description,price,stock,created_at,last_updated from product")
	if err != nil {
		return []*Product{}, err
	}
	defer rows.Close()
	pr := []*Product{}
	for rows.Next() {
		product := &Product{}
		err := rows.Scan(&product.ID, &product.UserID, &product.Name, &product.Description, &product.Price, &product.Stock, &product.CreatedAt, &product.LastUpdated)
		if err != nil {
			return []*Product{}, err
		}
		pr = append(pr, product)
	}
	return pr, nil
}

func (r *repository) GetProductByID(ctx context.Context, id int) (*Product, error) {
	product := &Product{}

	err := r.db.QueryRowContext(ctx, "select id,userID,name,description,price,stock,created_at,last_updated from product where id = ?", id).Scan(&product.ID, &product.UserID, &product.Name, &product.Description, &product.Price, &product.Stock, &product.CreatedAt, &product.LastUpdated)
	if err != nil {
		if err == sql.ErrNoRows {
			return &Product{}, utils.ErrNotFound
		}
		return &Product{}, err
	}

	return product, nil

}

func (r *repository) GetProductByName(ctx context.Context, name string) (*Product, error) {
	product := &Product{}

	err := r.db.QueryRowContext(ctx, "select id,userID,name,description,price,stock,created_at,last_updated from product where name = ?", name).Scan(&product.ID, &product.UserID, &product.Name, &product.Description, &product.Price, &product.Stock, &product.CreatedAt, &product.LastUpdated)
	if err != nil {
		if err == sql.ErrNoRows {
			return &Product{}, utils.ErrNotFound
		}
		return &Product{}, err
	}

	return product, nil
}

func (r *repository) GetProductByUserID(ctx context.Context, userID int) ([]*Product, error) {
	rows, err := r.db.QueryContext(ctx, "select id,userID,name,description,price,stock,created_at,last_updated from product where UserID = ?", userID)
	if err != nil {
		return []*Product{}, err
	}
	defer rows.Close()
	pr := []*Product{}
	for rows.Next() {
		product := &Product{}
		err := rows.Scan(&product.ID, &product.UserID, &product.Name, &product.Description, &product.Price, &product.Stock, &product.CreatedAt, &product.LastUpdated)
		if err != nil {
			return []*Product{}, err
		}
		pr = append(pr, product)
	}
	return pr, nil
}

func (r *repository) CreateProductImage(ctx context.Context, img *ProductImage) error {

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {

		return err
	}

	_, err = tx.ExecContext(ctx, "insert into product_image(productID,Url,Captions) values(?,?,?)", img.ProductID, img.Url, img.Captions)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r *repository) GetProductImage(ctx context.Context, productID int) ([]*ProductImage, error) {
	pm := []*ProductImage{}
	rows, err := r.db.QueryContext(ctx, "select id,productid,url,captions from product_image where productid = ?", productID)
	if err != nil {
		return []*ProductImage{}, err
	}

	defer rows.Close()
	for rows.Next() {
		pr := &ProductImage{}
		err := rows.Scan(&pr.ID, &pr.Captions, &pr.Url, &pr.Captions)
		if err != nil {
			return []*ProductImage{}, err
		}
		pm = append(pm, pr)
	}

	return pm, nil
}


func (r *repository) CreateRoomChat(ctx context.Context, URC *UlasanRoomChatProduct) (int, string, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, "", err
	}

	result, err := tx.ExecContext(ctx, "insert into ulasan_room_chat_product(roomID,userID,ProductID,username) values(?,?,?,?)", URC.RoomID, URC.UserID, URC.ProductID, URC.Username)
	if err != nil {
		tx.Rollback()
		return 0, "", err
	}

	roomID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, "", err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return 0, "", err
	}

	return int(roomID), URC.RoomID, nil

}

func (r *repository) CreateRoomChatMessage(ctx context.Context, URC *UlasanRoomChatProductMessage) (int, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	result, err := tx.ExecContext(ctx, "insert into ulasan_room_chat_product_message(roomChatID,message,sendAt,isDeleted) values(?,?,?,?)", URC.RoomID, URC.Message, URC.SendAt, URC.IsDeleted)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return int(id), nil
}

func (r *repository) GetRoomByProductID(ctx context.Context, ProductID int) ([]*UlasanRoomChatProduct, error) {
	ulrSlice := []*UlasanRoomChatProduct{}
	rows, err := r.db.QueryContext(ctx, "select id,roomid,userID,productID,username from ulasan_room_chat_product where productID = ?", ProductID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		ulr := &UlasanRoomChatProduct{}
		err := rows.Scan(&ulr.ID, &ulr.RoomID, &ulr.UserID, &ulr.ProductID, &ulr.Username)
		if err != nil {
			return nil, err
		}
		ulrSlice = append(ulrSlice, ulr)
	}
	return ulrSlice, nil
}
