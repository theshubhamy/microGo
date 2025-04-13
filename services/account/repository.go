package account

import (
	"context"
	"database/sql"
)

type Repository interface {
	Close() error
	PutAccount(ctx context.Context, acc Account) error
	GetAccountbyId(ctx context.Context, id string) (*Account, error)
	ListAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error)
}

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (Repository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &postgresRepository{db}, nil
}

func (r *postgresRepository) Close() error {
	return r.db.Close()
}

func (r *postgresRepository) Ping() error {
	return r.db.Ping()
}

func (r *postgresRepository) PutAccount(ctx context.Context, acc Account) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO accounts(id,name) VALUES($1,$2)", acc.ID, acc.Name)
	return err

}
func (r *postgresRepository) GetAccountbyId(ctx context.Context, id string) (*Account, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id,name FROM account WHERE id = $1", id)
	a := &Account{}
	err := row.Scan(&a.ID, &a.Name)
	if err != nil {
		return nil, err
	}
	return a, err
}
func (r *postgresRepository) ListAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id,name FROM accounts ORDER BY id DESC OFFSET $1 LIMIT $2", skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// Create a slice to hold the accounts

	accounts := []Account{}

	// Iterate through the rows and scan each row into an Account struct

	for rows.Next() {
		a := Account{}
		err = rows.Scan(&a.ID, &a.Name)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, a)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return accounts, err

}
