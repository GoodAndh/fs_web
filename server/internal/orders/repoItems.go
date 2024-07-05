package orders

import "context"

func (r *repository) GetOrderItems(ctx context.Context, userID, productID int) (*OrderItems, error) {
	ord := &OrderItems{}
	err := r.db.QueryRowContext(ctx, "select id,orderID,userID,productID,total,price from orders_items where userid =? and productID=?", userID, productID).Scan(&ord.ID, &ord.OrderID, &ord.UserID, &ord.ProductID, &ord.Total, &ord.Price)
	if err != nil {
		return &OrderItems{}, err
	}

	return ord, nil

}

func (r *repository) GetItemsByUserID(ctx context.Context, userID int) ([]*OrderItems, error) {
	ord := []*OrderItems{}
	rows, err := r.db.QueryContext(ctx, "select id,orderID,userID,productID,total,price from orders_items where userid =? ", userID)
	if err != nil {
		return []*OrderItems{}, err
	}
	defer rows.Close()
	for rows.Next() {
		o := &OrderItems{}
		err := rows.Scan(&o.ID, &o.OrderID, &o.UserID, &o.ProductID, &o.Total, &o.Price)
		if err != nil {
			return []*OrderItems{}, err
		}
		ord = append(ord, o)
	}

	return ord, nil
}
