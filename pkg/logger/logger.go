// pkg/logger/logger.go
package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
	With(fields ...zap.Field) *zap.Logger
}

type Config struct {
	Level      string `mapstructure:"level"`       // debug, info, warn, error
	Encoding   string `mapstructure:"encoding"`    // json, console
	OutputPath string `mapstructure:"output_path"` // stdout, file path
}

func New(cfg Config) (*zap.Logger, error) {
	level := zap.NewAtomicLevel()
	if err := level.UnmarshalText([]byte(cfg.Level)); err != nil {
		return nil, err
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var outputPaths []string
	if cfg.OutputPath != "" {
		outputPaths = append(outputPaths, cfg.OutputPath)
	}
	outputPaths = append(outputPaths, "stdout")

	config := zap.Config{
		Level:            level,
		Development:      false,
		Encoding:         cfg.Encoding,
		EncoderConfig:    encoderConfig,
		OutputPaths:      outputPaths,
		ErrorOutputPaths: []string{"stderr"},
	}

	return config.Build()
}

func SugarLogger(logger *zap.Logger) *zap.SugaredLogger {
	return logger.Sugar()
}
