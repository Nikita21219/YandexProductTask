package pkg

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
	"main/pkg/utils"
	"time"
)

type DBClient interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

func NewClient(ctx context.Context, username, password, host, port, dbName string) (pool *pgxpool.Pool, err error) {
	queryConnection := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?%s", username, password, host, port, dbName)

	err = utils.DoWithTries(func() error {

		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		pool, err = pgxpool.Connect(ctx, queryConnection)
		if err != nil {
			return err
		}
		return nil

	}, 4, 5*time.Second)

	if err != nil {
		return nil, err
	}

	return pool, nil
}

//CREATE TABLE "order" (
//id SERIAL PRIMARY KEY,
//weight INTEGER,
//region INTEGER,
//delivery_time VARCHAR(13),
//price INTEGER
//);
