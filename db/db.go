package db

import (
	"context"
	"fmt"

	"github.com/spf13/viper"
	"github.com/surrealdb/surrealdb.go"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type DatabaseParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Logger    *zap.Logger
	Config    *viper.Viper
}

type Database struct {
	DB *surrealdb.DB
}

// Module exports the database module.
var Module = fx.Options(
	fx.Provide(NewDatabase),
)

// NewDatabase initializes the database connection using the provided configuration and logger.
func NewDatabase(params DatabaseParams) (*Database, error) {
	dbAddress := params.Config.GetString("database.address")
	dbPort := params.Config.GetString("database.port")
	dbUsername := params.Config.GetString("database.username")
	dbPasswd := params.Config.GetString("database.password")

	connString := fmt.Sprintf("ws://%s:%s/rpc", dbAddress, dbPort)
	db, err := surrealdb.New(connString)
	if err != nil {
		params.Logger.Error("Couldn't connect to DB!",
			zap.Error(err),
			zap.String("address", dbAddress),
			zap.String("port", dbPort),
			zap.String("connString", connString),
		)
		return nil, err
	}

	if _, err = db.Signin(map[string]interface{}{
		"user": dbUsername,
		"pass": dbPasswd,
	}); err != nil {
		params.Logger.Error("Couldn't sign in to DB!",
			zap.Error(err),
			zap.String("username", dbUsername),
		)
		return nil, err
	}

	if _, err = db.Use("PackageLock", "db1.0"); err != nil {
		params.Logger.Panic("Couldn't use 'PackageLock' Namespace and 'db1.0' Database.",
			zap.Error(err),
		)
		return nil, err
	}

	params.Logger.Info("Successfully connected to DB.",
		zap.String("address", dbAddress),
		zap.String("port", dbPort),
	)

	database := &Database{
		DB: db,
	}

	// Use Lifecycle to manage the database connection
	params.Lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			params.Logger.Info("Closing database connection.")
			db.Close()
			return nil
		},
	})

	return database, nil
}
