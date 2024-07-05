package orders

import "context"

type repository struct {
	db DBTX
}

func NewRepository(db DBTX) Repository {
	return &repository{db: db}
}

func (r *repository) GetOrderStatus(ctx context.Context, userID, ProductID int) (*OrderStatus, error) {
	orS := &OrderStatus{}

	err := r.db.QueryRowContext(ctx, "select id,productID,userID,status,totalprice from orders_status where userid=? and productid=?", userID, ProductID).Scan(&orS.ID, &orS.ProductID, &orS.UserID, &orS.Status, &orS.TotalPrice)
	if err != nil {
		return &OrderStatus{}, err
	}

	return orS, nil
}

func (r *repository) CreateOrders(ctx context.Context, ord *OrderStatus, ordItems *OrderItems) (int, int, error) {

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, 0, err
	}

	result, err := tx.ExecContext(ctx, "insert into orders_status (productID,userID,status,totalprice) values (?,?,?,?)", ord.ProductID, ord.UserID, ord.Status, ord.TotalPrice)
	if err != nil {
		tx.Rollback()
		return 0, 0, err
	}

	idStatus, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, 0, err
	}

	res, err := tx.ExecContext(ctx, "insert into orders_items(orderID,userID,productID,total,price) values(?,?,?,?,?)", idStatus, ordItems.UserID, ordItems.ProductID, ordItems.Total, ordItems.Price)
	if err != nil {
		tx.Rollback()
		return 0, 0, err
	}

	idItems, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, 0, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return 0, 0, err
	}

	return int(idStatus), int(idItems), nil
}

func (r *repository) GetStatusByUserID(ctx context.Context, userID int) ([]*OrderStatus, error) {
	ord := []*OrderStatus{}
	rows, err := r.db.QueryContext(ctx, "select id,productID,userID,status,totalprice from orders_status where userid=? ", userID)
	if err != nil {
		return []*OrderStatus{}, err
	}
	defer rows.Close()
	for rows.Next() {
		o := &OrderStatus{}
		err := rows.Scan(&o.ID, &o.ProductID, &o.UserID, &o.Status, &o.TotalPrice)
		if err != nil {
			return []*OrderStatus{}, err
		}
		ord = append(ord, o)
	}

	return ord, nil
}
