package logger

import (
	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

func Init() {
	log, _ := zap.NewProduction(zap.AddCallerSkip(1))
	logger = log.Sugar()
}

func LocalInit() {
	log, _ := zap.NewDevelopment(zap.AddCallerSkip(1))
	logger = log.Sugar()
}

func Sync() {
	if logger != nil {
		logger.Sync()
	}
}

// Debug logs a message at level Debug on the standard logger.
func Debug(msg string, args ...interface{}) {
	if logger == nil {
		return
	}
	logger.Debug(args...)
}

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
	if logger == nil {
		return
	}
	logger.Info(args...)
}

// Warn logs a message at level Warn on the standard logger.
func Warn(args ...interface{}) {
	if logger == nil {
		return
	}
	logger.Warn(args...)
}

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) {
	if logger == nil {
		return
	}
	logger.Error(args...)
}

// Panic logs a message at level Panic on the standard logger.
func Panic(args ...interface{}) {
	if logger == nil {
		return
	}
	logger.Panic(args...)
}

// Fatal logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatal(args ...interface{}) {
	if logger == nil {
		return
	}
	logger.Fatal(args...)
}

// Debugf logs a message at level Debug on the standard logger.
func Debugf(format string, args ...interface{}) {
	if logger == nil {
		return
	}
	logger.Debugf(format, args...)
}

// Infof logs a message at level Info on the standard logger.
func Infof(format string, args ...interface{}) {
	if logger == nil {
		return
	}
	logger.Infof(format, args...)
}

// Warnf logs a message at level Warn on the standard logger.
func Warnf(format string, args ...interface{}) {
	if logger == nil {
		return
	}
	logger.Warnf(format, args...)
}

// Errorf logs a message at level Error on the standard logger.
func Errorf(format string, args ...interface{}) {
	if logger == nil {
		return
	}
	logger.Errorf(format, args...)
}

// Panicf logs a message at level Panic on the standard logger.
func Panicf(format string, args ...interface{}) {
	if logger == nil {
		return
	}
	logger.Panicf(format, args...)
}

// Fatalf logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatalf(format string, args ...interface{}) {
	if logger == nil {
		return
	}
	logger.Fatalf(format, args...)
}
