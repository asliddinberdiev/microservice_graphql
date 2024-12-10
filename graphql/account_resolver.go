package graphql

import (
	"context"
	"time"
)

type accountResolver struct {
	server *Server
}

func (r *accountResolver) Orders(ctx context.Context, account *Account) ([]*Order, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	orderList, err := r.server.orderClient.GetOrdersForAccount(ctx, account.ID)
	if err != nil {
		return nil, err
	}

	orders := make([]*Order, 0, len(orderList))
	for _, o := range orderList {
		products := make([]*OrderedProduct, len(o.Products))
		for _, p := range o.Products {
			products = append(products, &OrderedProduct{ID: p.ID, Name: p.Name, Description: p.Description, Price: int(p.Price), Quantity: int(p.Quantity)})
		}

		orders = append(orders, &Order{ID: o.ID, TotalPrice: int(o.TotalPrice), CreatedAt: o.CreatedAt, Products: products})
	}

	return orders, nil
}
