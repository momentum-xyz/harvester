package log

import (
	"github.com/ory/x/errorsx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	DefaultLogger = NewLoggerX(zap.DebugLevel)
)

type Logger struct {
	unsugared *zap.Logger
	*zap.SugaredLogger
	skipCallerLogger *zap.SugaredLogger
	level            zapcore.Level
}

func NewLogger(level zapcore.Level, opts ...zap.Option) (*Logger, error) {
	l, err := zap.Config{
		DisableStacktrace: true,
		Level:             zap.NewAtomicLevelAt(level),
		Development:       true,
		Encoding:          "console",
		EncoderConfig:     zap.NewDevelopmentEncoderConfig(),
		OutputPaths:       []string{"stderr"},
		ErrorOutputPaths:  []string{"stderr"},
	}.Build(opts...)

	if err != nil {
		return nil, err
	}

	return &Logger{
		unsugared:        l,
		skipCallerLogger: l.WithOptions(zap.AddCallerSkip(1)).Sugar(),
		SugaredLogger:    l.Sugar(),
		level:            level,
	}, nil
}

func NewLoggerX(level zapcore.Level, opts ...zap.Option) *Logger {
	logger, err := NewLogger(level, opts...)
	if err != nil {
		panic(err)
	}
	return logger
}

func (l *Logger) Named(name string, opts ...zap.Option) *Logger {
	unsugared := l.unsugared.Named(name).WithOptions(opts...)

	return &Logger{
		unsugared:        unsugared,
		skipCallerLogger: unsugared.WithOptions(zap.AddCallerSkip(1)).Sugar(),
		SugaredLogger:    unsugared.Sugar(),
		level:            l.level,
	}
}

func (l *Logger) With(args ...interface{}) *Logger {
	sugared := l.SugaredLogger.With(args...)
	unsugared := sugared.Desugar()
	return &Logger{
		unsugared:        unsugared,
		skipCallerLogger: unsugared.WithOptions(zap.AddCallerSkip(1)).Sugar(),
		SugaredLogger:    sugared,
		level:            l.level,
	}
}

func (l *Logger) Error(err error) bool {
	if err == nil {
		return false
	}
	if l.level == zap.DebugLevel {
		err = errorsx.WithStack(err)
		// include the stack trace in the error message
		l.skipCallerLogger.Errorf("error: %+v\n", err)
	} else {
		l.skipCallerLogger.Errorf("error: %v\n", err)
	}
	return true
}

func Debugf(a string, args ...interface{}) {
	DefaultLogger.skipCallerLogger.Debugf(a, args...)
}

func Infof(a string, args ...interface{}) {
	DefaultLogger.skipCallerLogger.Infof(a, args...)
}

func Errorf(a string, args ...interface{}) {
	DefaultLogger.skipCallerLogger.Errorf(a, args...)
}

func Debug(args ...interface{}) {
	DefaultLogger.skipCallerLogger.Debug(args...)
}

func Info(args ...interface{}) {
	DefaultLogger.skipCallerLogger.Info(args...)
}

func Error(err error) bool {
	if err == nil {
		return false
	}
	if DefaultLogger.level == zap.DebugLevel {
		// include the stack trace in the error message
		DefaultLogger.skipCallerLogger.Errorf("error: %+v\n", err)
	} else {
		DefaultLogger.skipCallerLogger.Errorf("error: %v\n", err)
	}
	return true
}

func Fatal(args ...interface{}) {
	DefaultLogger.Fatal(args...)
}
