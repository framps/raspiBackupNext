package tools

import (
	"fmt"
	"os"

	"go.uber.org/zap"
)

var (
	// Log -
	Log *zap.SugaredLogger
)

// NewLogger -
func NewLogger(debug bool) *zap.SugaredLogger {

	cfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	_, envDefind := os.LookupEnv("DEBUG")

	if !envDefind && !debug {
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
		cfg.OutputPaths = append(cfg.OutputPaths, "./raspiBackup.log")
		cfg.ErrorOutputPaths = append(cfg.ErrorOutputPaths, "./raspiBackup.log")
	}
	logger, err := cfg.Build()
	if err != nil {
		fmt.Printf("Unable to create logger. Root cause: %s", err.Error())
		os.Exit(42)
	}
	Log = logger.Sugar()
	return Log
}

// Sync -
func Sync() {
	Log.Sync()
}

func HandleError(err error) {
	if err != nil {
		panic(err)
	}
}
