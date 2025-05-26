package order

import (
	"context"
	"database/sql"
	"log"

	"github.com/lib/pq"
)

type Repository interface {
	Close() error
	PutOrder(ctx context.Context, o Order) error
	GetOrderforAccount(ctx context.Context, accountId string) ([]Order, error)
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

func (r *postgresRepository) PutOrder(ctx context.Context, o Order) (err error) {
	txn, err := r.db.BeginTx(ctx, nil)
	defer func() {
		if err != nil {
			log.Println(err)
			txn.Rollback()
			return
		}
		err = txn.Commit()
	}()
	txn.ExecContext(ctx, "INSERT INTO orders(id,created_at,account_id,total_price) VALUES ($1,$2,$3,$4)", o.ID, o.CreatedAt, o.AccountId, o.TotalPrice)
	if err != nil {
		log.Println(err)
		return
	}

	stmt, _ := txn.PrepareContext(ctx, pq.CopyIn("order_products", "order_id", "product_id", "quantity"))
	for _, p := range o.Products {
		_, err = stmt.ExecContext(ctx, o.ID, p.ID, p.Quantity)
		if err != nil {
			log.Println(err)
			return
		}

	}
	_, err = stmt.ExecContext(ctx)
	if err != nil {
		log.Println(err)
		return
	}
	stmt.Close()
	return
}

func (r *postgresRepository) GetOrderforAccount(ctx context.Context, accountId string) ([]Order, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT o.id,o.created_at,o.account_id, o.total_price::money::numeric::float8,op.product_id,op.quantity FROM orders o JOIN order_products op ON(o.id=op.order_id) WHERE o.account_id=$1,ORDER BY o.id`, accountId,
	)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	orders := []Order{}
	order := &Order{}
	lastOrder := &Order{}
	orderedProduct := &OrderedProduct{}
	products := []OrderedProduct{}

	for rows.Next() {
		if err = rows.Scan(
			&order.ID,
			&order.CreatedAt,
			&order.AccountId,
			&order.TotalPrice,
			&orderedProduct.ID,
			&orderedProduct.Quantity,
		); err != nil {
			return nil, err
		}

		if lastOrder.ID != " " && lastOrder.ID != order.ID {
			newOrder := Order{
				ID:         lastOrder.ID,
				AccountId:  lastOrder.AccountId,
				CreatedAt:  lastOrder.CreatedAt,
				TotalPrice: lastOrder.TotalPrice,
				Products:   lastOrder.Products,
			}
			orders = append(orders, newOrder)
			products = []OrderedProduct{}

		}

		products = append(products, OrderedProduct{
			ID:       orderedProduct.ID,
			Quantity: orderedProduct.Quantity,
		})
		*lastOrder = *order
	}

	newOrder := Order{
		ID:         lastOrder.ID,
		AccountId:  lastOrder.AccountId,
		CreatedAt:  lastOrder.CreatedAt,
		TotalPrice: lastOrder.TotalPrice,
		Products:   lastOrder.Products,
	}
	orders = append(orders, newOrder)

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}
