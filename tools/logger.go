package tools

//######################################################################################################################
//
//    Next raspiBackup version written in go
//
//    Copyright (C) 2018 framp at linux-tips-and-tricks dot de
//
//#######################################################################################################################

import (
	"fmt"
	"os"

	"go.uber.org/zap"
)

var (
	Logger *zap.SugaredLogger
)

// NewLogger -
func NewLogger(debug bool) {

	cfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	_, envDefind := os.LookupEnv("DEBUG")

	if !envDefind && debug {
		cfg.OutputPaths = []string{"./raspiBackup.log"}
		cfg.ErrorOutputPaths = []string{"./raspiBackup.log"}
	} else if envDefind && debug {
		cfg.OutputPaths = append(cfg.OutputPaths, "./raspiBackup.log")
		cfg.ErrorOutputPaths = append(cfg.ErrorOutputPaths, "./raspiBackup.log")
	} else {
		cfg.Level = zap.NewAtomicLevelAt(zap.FatalLevel)
	}
	logger, err := cfg.Build()
	if err != nil {
		fmt.Printf("Unable to create logger. Root cause: %s", err.Error())
		os.Exit(42)
	}

	zap.ReplaceGlobals(logger)
	Logger = zap.S()

}

// HandleError -
func HandleError(err error) {
	if err != nil {
		panic(err)
	}
}
