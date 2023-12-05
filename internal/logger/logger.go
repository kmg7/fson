package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Options struct {
	Development bool
	Output      *os.File
	Silent      bool
}

var logger *zap.SugaredLogger

func InitLogger(opt Options) {
	level := zap.NewAtomicLevel()
	lcfg := zap.NewProductionEncoderConfig()
	if opt.Development {
		lcfg = zap.NewDevelopmentEncoderConfig()

	}
	lcfg.LevelKey = "level"
	lcfg.TimeKey = "timestamp"
	lcfg.CallerKey = "caller"
	lcfg.MessageKey = "message"

	lcfg.EncodeLevel = zapcore.CapitalLevelEncoder
	lcfg.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(lcfg)

	lcfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	lcfg.EncodeTime = func(t time.Time, pae zapcore.PrimitiveArrayEncoder) {
		pae.AppendString(t.Format("15:04:05"))
	}
	consoleEncoder := zapcore.NewConsoleEncoder(lcfg)
	cores := []zapcore.Core{
		zapcore.NewCore(fileEncoder, opt.Output, level),
	}
	if !opt.Silent {
		stdout := zapcore.AddSync(os.Stdout)
		cores = append(cores, zapcore.NewCore(consoleEncoder, stdout, level))
	}
	core := zapcore.NewTee(cores...)
	l := zap.New(core)
	logger = l.Sugar()
	Info = logger.Info
	Warn = logger.Warn
	Error = logger.Error
	Fatal = logger.Fatal
	Debug = logger.Debug
	Info("Logger Initialized Successfully")
}

var Info func(args ...interface{})
var Warn func(args ...interface{})
var Error func(args ...interface{})
var Fatal func(args ...interface{})
var Debug func(args ...interface{})
