package logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

// InitLogger initializes the logger.
// env examples: "dev", "prod" (case-insensitive). Unknown values will be used as-is in file names.
func InitLogger(env string) {
	// Normalize and default
	env = strings.ToLower(strings.TrimSpace(env))
	if env == "" {
		env = "prod"
	}

	// Ensure logs directory exists
	_ = os.MkdirAll("logs", 0o755)

	// Asia/Kolkata timestamp
	loc, _ := time.LoadLocation("Asia/Kolkata")
	timeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.In(loc).Format("02 Jan 2006 15:04:05"))
	}

	// Common encoder config for JSON files
	fileEncCfg := zapcore.EncoderConfig{
		TimeKey:       "timestamps",
		LevelKey:      "level",
		MessageKey:    "message",
		CallerKey:     "caller",
		StacktraceKey: "trace",
		EncodeLevel:   zapcore.CapitalLevelEncoder, // INFO/WARN/ERROR
		EncodeTime:    timeEncoder,
		EncodeCaller:  zapcore.ShortCallerEncoder,
	}

	// Console encoder config (colored)
	consoleEncCfg := fileEncCfg
	consoleEncCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	jsonEncoder := zapcore.NewJSONEncoder(fileEncCfg)
	consoleEncoder := zapcore.NewConsoleEncoder(consoleEncCfg)

	// Build filenames based on env
	infoFilePath := fmt.Sprintf("logs/info.%s.log", env)
	errorFilePath := fmt.Sprintf("logs/error.%s.log", env)

	// Open files (fallback to stdout/stderr if file open fails)
	infoFile, err1 := os.OpenFile(infoFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err1 != nil {
		infoFile = os.Stdout
	}
	errorFile, err2 := os.OpenFile(errorFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err2 != nil {
		errorFile = os.Stderr
	}

	infoWS := zapcore.AddSync(infoFile)
	errorWS := zapcore.AddSync(errorFile)

	// Level filters
	infoLevel := zap.LevelEnablerFunc(func(l zapcore.Level) bool { return l >= zapcore.InfoLevel })
	errorLevel := zap.LevelEnablerFunc(func(l zapcore.Level) bool { return l >= zapcore.ErrorLevel })

	// Build cores: info file, error file, and dev console
	cores := []zapcore.Core{
		zapcore.NewCore(jsonEncoder, infoWS, infoLevel),   // writes INFO/WARN/ERROR to info.<env>.log
		zapcore.NewCore(jsonEncoder, errorWS, errorLevel), // writes ERROR+ to error.<env>.log
	}
	if env == "dev" {
		cores = append(cores, zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), infoLevel))
	}

	log = zap.New(
		zapcore.NewTee(cores...),
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel), // include stack trace for Error+
	)
}

// Wrapper functions for direct usage (Winston-like)
func Info(msg string, fields ...zap.Field)  { log.Info(msg, fields...) }
func Error(msg string, fields ...zap.Field) { log.Error(msg, fields...) }
func Debug(msg string, fields ...zap.Field) { log.Debug(msg, fields...) }
func Warn(msg string, fields ...zap.Field)  { log.Warn(msg, fields...) }

// Sync flushes any buffered log entries (call at shutdown).
func Sync() { _ = log.Sync() }
