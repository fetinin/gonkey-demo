package app

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func NewDB(ctx context.Context, dsn string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	if err := conn.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %v", err)
	}

	return conn, nil
}
