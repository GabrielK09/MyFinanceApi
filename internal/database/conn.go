package database

import (
	"context"
	"fmt"
	l "my_finance/internal/logger"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var Pool *pgxpool.Pool

func Init() (*pgxpool.Pool, error) {
	l.InfoLogger.Println("Conectando ao db ...")

	_ = godotenv.Load()

	if Pool != nil {
		return Pool, nil
	}

	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		return nil, fmt.Errorf("DB_URL não definida")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 35*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(dsn)

	if err != nil {
		l.ErrorLogger.Println("Erro ao carregar as configurações:", err)
		return nil, err
	}

	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		l.ErrorLogger.Println("Erro ao criar pool de conexões:", err)
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		l.ErrorLogger.Println("Erro ao conectar ao banco:", err)
		return nil, err
	}

	Pool = pool

	return Pool, nil
}
