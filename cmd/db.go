package main

import (
	"context"
	"log"
	"strings"
	"time"

	"corrigan.io/go_api_seed/internal/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
)

type myQueryTracer struct {
	logger *zerolog.Logger
}

func (tracer *myQueryTracer) TraceQueryStart(
	ctx context.Context,
	_ *pgx.Conn,
	data pgx.TraceQueryStartData) context.Context {

	start := time.Now()
	cleanedSQL := strings.ReplaceAll(data.SQL, "\n", " ")
	operationName := strings.Split(cleanedSQL, ":")

	if len(operationName) > 2 {
		context.AfterFunc(ctx, func() {
			tracer.logger.Trace().
				Ctx(ctx).
				Str("operation", operationName[1]).
				Interface("args", data.Args).
				Dur("time", time.Since(start)).Msg("DATABASE")
		})
	}

	return ctx
}

func (tracer *myQueryTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
}

func configDB(cfg config.Config, logger *zerolog.Logger) *pgxpool.Config {
	const defaultMaxConns = int32(4)
	const defaultMinConns = int32(0)
	const defaultMaxConnLifetime = time.Hour
	const defaultMaxConnIdleTime = time.Minute * 30
	const defaultHealthCheckPeriod = time.Minute
	const defaultConnectTimeout = time.Second * 5

	dbConfig, err := pgxpool.ParseConfig(cfg.DB.DSN)
	if err != nil {
		log.Fatal("Failed to create a config, error: ", err)
	}

	dbConfig.MaxConns = defaultMaxConns
	dbConfig.MinConns = defaultMinConns
	dbConfig.MaxConnLifetime = defaultMaxConnLifetime
	dbConfig.MaxConnIdleTime = defaultMaxConnIdleTime
	dbConfig.HealthCheckPeriod = defaultHealthCheckPeriod
	dbConfig.ConnConfig.ConnectTimeout = defaultConnectTimeout
	// if cfg.DB.Logging {
	// 	dbConfig.ConnConfig.Tracer = &myQueryTracer{logger}
	// }
	dbConfig.ConnConfig.Tracer = &myQueryTracer{logger}

	return dbConfig
}

func openDBPool(cfg config.Config, logger *zerolog.Logger) (*pgxpool.Pool, error) {
	dbConfigurationOptions := configDB(cfg, logger)

	dbpool, err := pgxpool.NewWithConfig(context.Background(), dbConfigurationOptions)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	connection, err := dbpool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer connection.Release()

	err = dbpool.Ping(ctx)

	if err != nil {
		dbpool.Close()
		return nil, err
	}

	return dbpool, nil
}
