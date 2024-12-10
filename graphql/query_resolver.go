package graphql

import (
	"context"
	"time"
)

type queryResolver struct {
	server *Server
}

func (r *queryResolver) Accounts(ctx context.Context, pagination *PaginationInput, id *string) ([]*Account, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if id != nil {
		r, err := r.server.accountClient.GetAccountByID(ctx, *id)
		if err != nil {
			return nil, err
		}

		return []*Account{{ID: r.ID, Name: r.Name}}, nil
	}

	skip, take := int(0), int(0)
	if pagination != nil {
		skip, take = pagination.bounds()
	}

	accountList, err := r.server.accountClient.GetAccounts(ctx, uint64(skip), uint64(take))
	if err != nil {
		return nil, err
	}

	accounts := make([]*Account, 0, len(accountList))
	for _, a := range accountList {
		accounts = append(accounts, &Account{ID: a.ID, Name: a.Name})
	}
	return accounts, nil
}

func (r *queryResolver) Products(ctx context.Context, pagination *PaginationInput, query *string, id *string) ([]*Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if id != nil {
		r, err := r.server.catalogClient.GetProductByID(ctx, *id)
		if err != nil {
			return nil, err
		}

		return []*Product{{ID: r.ID, Name: r.Name, Description: r.Description, Price: int(r.Price), Quantity: int(r.Quantity)}}, nil
	}

	skip, take := int(0), int(0)
	if pagination != nil {
		skip, take = pagination.bounds()
	}

	q := ""
	if query != nil {
		q = *query
	}

	productList, err := r.server.catalogClient.GetProducts(ctx, uint64(skip), uint64(take), nil, q)
	if err != nil {
		return nil, err
	}

	products := make([]*Product, 0, len(productList))
	for _, p := range productList {
		products = append(products, &Product{ID: p.ID, Name: p.Name, Description: p.Description, Price: int(p.Price), Quantity: int(p.Quantity)})
	}
	return products, nil
}

func (p PaginationInput) bounds() (int, int) {
	skipValue := int(0)
	takeValue := int(0)

	if p.Skip != nil {
		skipValue = *p.Skip
	}

	if p.Take != nil {
		takeValue = *p.Take
	}

	return skipValue, takeValue
}
