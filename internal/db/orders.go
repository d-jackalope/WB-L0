package db

import (
	"context"

	"github.com/d-jackalope/L0/pkg/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Orders interface {
	Create() error
	Insert(uid string, data []byte) error
	GetAllData() (map[string]models.Order, error)
	Exist(uid string) (bool, error)
}

type ordersDatabase struct {
	pool *pgxpool.Pool
	ctx  context.Context
}

func (db *ordersDatabase) Create() error {
	conn, err := db.pool.Acquire(db.ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	query := `
		CREATE TABLE IF NOT EXISTS orders (
			order_uid VARCHAR PRIMARY KEY,
			data JSONB
		);
	`

	_, err = conn.Exec(context.Background(), query)
	if err != nil {
		return err
	}
	return nil
}

func (db *ordersDatabase) Insert(uid string, data []byte) error {
	conn, err := db.pool.Acquire(db.ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	query := `
		INSERT INTO orders (order_uid, data)
		VALUES ($1, $2)
	`
	_, err = conn.Exec(context.Background(), query, uid, data)
	if err != nil {
		return err
	}
	return nil
}

func (db *ordersDatabase) GetAllData() (map[string]models.Order, error) {
	conn, err := db.pool.Acquire(db.ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	query := `
		SELECT data FROM orders;
	`
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	var order models.Order
	orders := make(map[string]models.Order)
	for rows.Next() {
		if err := rows.Scan(&order); err != nil {
			return nil, err
		}
		orders[order.OrderUID] = order
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}

func (db *ordersDatabase) Exist(uid string) (bool, error) {
	conn, err := db.pool.Acquire(db.ctx)
	if err != nil {
		return false, err
	}
	defer conn.Release()

	query := `
		SELECT EXISTS (SELECT 1 FROM orders WHERE order_uid = $1);
	`
	var exist bool
	err = conn.QueryRow(context.Background(), query, uid).Scan(&exist)
	if err != nil {
		return false, err
	}

	return exist, nil
}
