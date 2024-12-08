package account

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	_ "github.com/lib/pq"
)

type Repository interface {
	Close()
	Ping() error
	CreateAccount(ctx context.Context, name string) (*Account, error)
	GetAccountByID(ctx context.Context, id string) (*Account, error)
	GetAccounts(ctx context.Context, skip, take uint64) ([]*Account, error)
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

func (r *PostgresRepository) Ping() error {
	return r.db.Ping()
}

func (r *PostgresRepository) CreateAccount(ctx context.Context, name string) (*Account, error) {
	query := `
		INSERT INTO accounts (name) 
		VALUES ($1)
		RETURNING id
	`
	var account Account
	account.Name = name
	if err := r.db.QueryRowContext(ctx, query, name).Scan(&account.ID); err != nil {
		return nil, errors.Wrap(err, "failed to scan account id")
	}

	return &account, nil
}

func (r *PostgresRepository) GetAccountByID(ctx context.Context, id string) (*Account, error) {
	query := `
		SELECT id, name 
		FROM accounts 
		WHERE id = $1
	`

	var account Account
	if err := r.db.QueryRowContext(ctx, query, id).Scan(&account.ID, &account.Name); err != nil {
		return nil, errors.Wrap(err, "failed to scan account")
	}

	if account.ID == "" {
		return nil, nil
	}

	return &account, nil
}

func (r *PostgresRepository) GetAccounts(ctx context.Context, skip, take uint64) ([]*Account, error) {
	query := `
		SELECT id, name 
		FROM accounts ORDER BY id DESC LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, take, skip)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query accounts")
	}
	defer rows.Close()

	accounts := make([]*Account, 0, take)
	for rows.Next() {
		var account Account
		if err := rows.Scan(&account.ID, &account.Name); err != nil {
			return nil, errors.Wrap(err, "failed to scan account")
		}

		accounts = append(accounts, &account)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "failed to get accounts")
	}

	return accounts, nil
}
