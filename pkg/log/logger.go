package log

import (
	"errors"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	ToolLogger                 *zap.Logger
	FileCopyLogger             *zap.Logger
	FileInfoLogger             *zap.Logger
	FileInfoErrorLogger 	   *zap.Logger
)

func Init(outputPath string) (err error) {
	ToolLogger, err = getLogger([]string{filepath.Join(outputPath, "tool.log"), filepath.Join(outputPath, "error.log")}, "Time", "Level")
	if err != nil {
		return err
	}
	FileCopyLogger, err = getLogger([]string{filepath.Join(outputPath, "file_copy.log"), ""}, "Time", "Level")
	if err != nil {
		return err
	}
	FileInfoLogger, err = getLogger([]string{filepath.Join(outputPath, "file_info.csv"), ""}, "", "")
	if err != nil {
		return err
	}
	FileInfoErrorLogger, err = getLogger([]string{filepath.Join(outputPath, "process_failed_files.log"), ""}, "Time", "Level")
	if err != nil {
		return err
	}
	return
}

func getLogger(logName []string, time string, level string) (*zap.Logger, error) {

	outputPath := []string{"stdout", logName[0]}
	errOutputPath := []string{"stderr"}
	if len(logName[1]) > 0{
		errOutputPath = append(errOutputPath, logName[1])
	}

	conf := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:       true,
		Encoding:          "console",
		DisableStacktrace: true,
		DisableCaller:     true,
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        time,
			LevelKey:       level,
			NameKey:        "Name",
			CallerKey:      "Caller",
			MessageKey:     "Msg",
			StacktraceKey:  "Stack",
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      outputPath,
		ErrorOutputPaths: errOutputPath,
	}

	logger, err := conf.Build()
	if err != nil {
		return nil, errors.New("[-] Failed to build zap logger")
	}

	return logger, nil
}
