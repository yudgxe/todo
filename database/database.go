package database

import (
	"context"
	"medods/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

var database *pgxpool.Pool

func MustInitDatabase(ctx context.Context, connStr string) {
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		utils.Panicf("error on parse config - %v", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		utils.Panicf("error on connect db - %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		utils.Panicf("error on Ping db - %v", err)
	}

	database = pool
}

func GetDatabase() *pgxpool.Pool {
	if database != nil {
		return database
	}

	panic("error on GetDatabase - need setup db")
}

func SetDatabase(db *pgxpool.Pool) { // for tests
	database = db
}
