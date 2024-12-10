package graphql

import (
	"context"
	"time"

	"github.com/asliddinberdiev/microservice_graphql/order"
	"github.com/pkg/errors"
)

var (
	ErrInvalidParameter = errors.New("invalid parameter")
)

type mutationResolver struct {
	server *Server
}

func (r *mutationResolver) CreateAccount(ctx context.Context, input AccountInput) (*Account, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	account, err := r.server.accountClient.CreateAccount(ctx, input.Name)
	if err != nil {
		return nil, err
	}

	return &Account{ID: account.ID, Name: account.Name}, nil
}

func (r *mutationResolver) CreateProduct(ctx context.Context, input ProductInput) (*Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	product, err := r.server.catalogClient.CreateProduct(ctx, input.Name, input.Description, uint64(input.Price), uint64(input.Quantity))
	if err != nil {
		return nil, err
	}

	return &Product{ID: product.ID, Name: product.Name, Description: product.Description, Price: int(product.Price), Quantity: int(product.Quantity)}, nil
}

func (r *mutationResolver) CreateOrder(ctx context.Context, input OrderInput) (*Order, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	products := make([]order.OrderedProduct, len(input.Products))
	for _, p := range input.Products {
		if p.Quantity <= 0 {
			return nil, ErrInvalidParameter
		}

		products = append(products, order.OrderedProduct{ID: p.ID, Quantity: uint64(p.Quantity)})
	}

	o, err := r.server.orderClient.CreateOrder(ctx, input.AccountID, products)
	if err != nil {
		return nil, err
	}

	return &Order{ID: o.ID, TotalPrice: int(o.TotalPrice), CreatedAt: o.CreatedAt}, nil	
}
