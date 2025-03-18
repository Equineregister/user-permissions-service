package application

import (
	"log"
	"log/slog"
	"os"

	"github.com/Equineregister/slogger"
)

func InitLogger() {
	loggerEnv, ok := slogger.EnvLevels[os.Getenv("LOG_LEVEL")]
	if !ok {
		slog.Warn("unsupported env supplied for logging options", "LOG_LEVEL", os.Getenv("LOG_LEVEL"))

		loggerEnv = slogger.EnvProd
	}

	// Setup logger.
	slogHandler, err := slogger.NewHandler(loggerEnv)
	if err != nil {
		log.Printf("failed to initialize slog handler: %v\n", err)
	}
	logger := slog.New(slogHandler)
	slog.SetDefault(logger) // Updates slogs default instance of slog with our own handler.
}

/*
func InitPostgres(ctx context.Context, cfg *config.Config) *pgxpool.Pool {
	db := cfg.Database()
	postgresConn, _, err := postgresutil.Connect(ctx, db.Port, db.Host, db.Username, db.Password, db.DatabaseName)
	if err != nil {
		slog.Error("Error connecting to database", "error", err.Error())

		os.Exit(1)
	}

	err = migrations.Migrate(ctx, postgresConn, postgres.Migrations)
	if err != nil {
		slog.Error("failed to migrate", "error", err.Error())

		os.Exit(1)
	}

	return postgresConn
}
*/
