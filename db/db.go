package db

import (
	"context"
	"fmt"

	"github.com/spf13/viper"
	"github.com/surrealdb/surrealdb.go"
	"go.opentelemetry.io/otel/codes" // Import for setting span status
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type DatabaseParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Logger    *zap.Logger
	Config    *viper.Viper
	Tracer    trace.Tracer // Injected Tracer from OpenTelemetry
}

type Database struct {
	DB     *surrealdb.DB
	Logger *zap.Logger
	Tracer trace.Tracer
}

// Module exports the database module.
var Module = fx.Options(
	fx.Provide(NewDatabase),
)

// NewDatabase initializes the database connection using the provided configuration and logger.
func NewDatabase(params DatabaseParams) (*Database, error) {
	// Start a new span for the database initialization
	_, span := params.Tracer.Start(context.Background(), "Database Initialization")
	defer span.End()

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
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to connect to DB")
		return nil, err
	}

	// Sign in to the database
	if _, err = db.Signin(map[string]interface{}{
		"user": dbUsername,
		"pass": dbPasswd,
	}); err != nil {
		params.Logger.Error("Couldn't sign in to DB!",
			zap.Error(err),
			zap.String("username", dbUsername),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to sign in to DB")
		return nil, err
	}

	// Use the specified namespace and database
	if _, err = db.Use("PackageLock", "db1.0"); err != nil {
		params.Logger.Panic("Couldn't use 'PackageLock' Namespace and 'db1.0' Database.",
			zap.Error(err),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to use Namespace and Database")
		return nil, err
	}

	params.Logger.Info("Successfully connected to DB.",
		zap.String("address", dbAddress),
		zap.String("port", dbPort),
	)
	span.AddEvent("Successfully connected and authenticated to DB")

	database := &Database{
		DB:     db,
		Logger: params.Logger,
		Tracer: params.Tracer,
	}

	// Use Lifecycle to manage the database connection
	params.Lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			params.Logger.Info("Closing database connection.")
			span.AddEvent("Closing database connection")
			db.Close()
			return nil
		},
	})

	return database, nil
}
