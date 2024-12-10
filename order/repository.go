package order

import (
	"context"
	"database/sql"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

type Repository interface {
	Close()
	UpdateOrder(ctx context.Context, o Order) error
	GetOrdersForAccount(ctx context.Context, accountID string) ([]*Order, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open database")
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(20)

	if err := db.Ping(); err != nil {
		return nil, errors.Wrap(err, "failed to ping database")
	}

	return &PostgresRepository{db: db}, nil
}

func (r *PostgresRepository) Close() {
	r.db.Close()
}

func (r *PostgresRepository) UpdateOrder(ctx context.Context, o Order) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	query := `
		INSERT INTO orders(id, account_id, total_price, created_at)
		VALUES ($1, $2, $3, $4)
	`
	if _, err := tx.ExecContext(ctx, query, o.ID, o.AccountID, o.TotalPrice, o.CreatedAt); err != nil {
		return errors.Wrap(err, "failed to insert order")
	}

	stmt, _ := tx.PrepareContext(ctx, pq.CopyIn("order_products", "order_id", "product_id", "quantity"))
	for _, p := range o.Products {
		if _, err := stmt.ExecContext(ctx, o.ID, p.ID, p.Quantity); err != nil {
			return errors.Wrap(err, "failed to insert order")
		}
	}
	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx); err != nil {
		return errors.Wrap(err, "failed to insert order")
	}

	return nil
}

func (r *PostgresRepository) GetOrdersForAccount(ctx context.Context, accountID string) ([]*Order, error) {
	query := `
		SELECT o.id, o.account_id, o.total_price::money::numeric, o.created_at, op.product_id, op.quantity
		FROM orders o JOIN order_products op ON o.id = op.order_id
		WHERE o.account_id = $1
		ORDER BY o.created_at DESC 
	`
	rows, err := r.db.QueryContext(ctx, query, accountID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query orders for account")
	}
	defer rows.Close()

	orderMap := make(map[string]*Order)

	for rows.Next() {
		var (
			id, accountID, productID string
			totalPrice, quantity     uint64
			createdAt                sql.NullTime
		)

		if err := rows.Scan(&id, &accountID, &totalPrice, &quantity, &productID, &createdAt); err != nil {
			return nil, errors.Wrap(err, "failed to scan order")
		}

		order, exists := orderMap[id]
		if !exists {
			order = &Order{ID: id, AccountID: accountID, TotalPrice: totalPrice, CreatedAt: createdAt.Time, Products: make([]OrderedProduct, 0)}
			orderMap[id] = order
		}

		order.Products = append(order.Products, OrderedProduct{ID: productID, Quantity: quantity})
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "failed to get rows")
	}

	orders := make([]*Order, 0, len(orderMap))
	for _, order := range orderMap {
		orders = append(orders, order)
	}

	return orders, nil
}
