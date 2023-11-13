package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	Orders
}

func NewPostgres(pool *pgxpool.Pool, ctx context.Context) *Postgres {
	return &Postgres{
		Orders: &ordersDatabase{
			pool: pool,
			ctx:  ctx,
		},
	}
}
