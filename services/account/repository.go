package account

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Account struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `josn:"phone"`
	Password string `json:"password"`
}

type Repository interface {
	Close() error
	PutAccount(ctx context.Context, acc Account) error
	GetAccount(ctx context.Context, key, value string) (*Account, error)
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
	_, err := r.db.ExecContext(ctx, "INSERT INTO accounts(id,name,email,phone,password) VALUES($1,$2,$3,$4,$5)", acc.ID, acc.Name, acc.Email, acc.Phone, acc.Password)
	return err
}

func (r *postgresRepository) GetAccount(ctx context.Context, key, value string) (*Account, error) {
	if !isAllowedKey(key) {
		return nil, fmt.Errorf("invalid column key: %s", key)
	}
	// construct the query safely since key is validated
	query := fmt.Sprintf(`SELECT id, name, email, phone,password FROM accounts WHERE %s = $1`, key)

	row := r.db.QueryRowContext(ctx, query, value)
	a := &Account{}
	err := row.Scan(&a.ID, &a.Name, &a.Email, &a.Phone, &a.Password)
	if err != nil {
		return nil, err
	}
	return a, err
}

func (r *postgresRepository) ListAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id,name,email,phone FROM accounts ORDER BY id DESC OFFSET $1 LIMIT $2", skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// Create a slice to hold the accounts

	accounts := []Account{}

	// Iterate through the rows and scan each row into an Account struct

	for rows.Next() {
		a := Account{}
		err = rows.Scan(&a.ID, &a.Name, &a.Email, &a.Phone)
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
