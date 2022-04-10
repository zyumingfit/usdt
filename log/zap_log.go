package log

import (
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

// Log print initialization, get *zap.Logger Info.
func InitLogger(logPath string, loglevel string) *zap.Logger {
	hook := lumberjack.Logger{
		Filename:   logPath, // log file path
		MaxSize:    128,     // megabytes
		MaxBackups: 30,      // max backup
		MaxAge:     7,       // days
		Compress:   true,    // is Compress, disabled by default
	}

	w := zapcore.AddSync(&hook)

	// set log print level
	var level zapcore.Level
	switch loglevel {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}
	encoderConfig := zap.NewProductionEncoderConfig()
	// time format
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), w), // this line enables log outputs to multiple destinations: log file/stdout
		level,
	)

	logger := zap.New(core, zap.AddStacktrace(zap.ErrorLevel))
	return logger
}

// Init zap log.
func init() {
	logger = InitLogger("usdt.log", "debug")
}

// Get zap log instance.
func Logger() *zap.Logger {
	return logger
}
