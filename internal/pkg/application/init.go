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
