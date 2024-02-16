package logger

import (
	"io"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type AppLogger interface {
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	// Debug(args ...interface{}) //TODO
}

type logger struct {
	zap *zap.SugaredLogger
}

// Logger options.
// There will be no log when Output file nil and Silent is true
type Options struct {
	Files  []io.Writer //It should implement Sync() error method.
	Stdout io.Writer   //It must be os.Stdout
	// Development bool
	// FatalFunc   *func()     //If it left nil logger call os.Exit after logging //TODO
}

// Info level log.
func (l *logger) Info(args ...interface{}) {
	l.zap.Info(args...)
}

// Warn level log.
func (l *logger) Warn(args ...interface{}) {
	l.zap.Warn(args...)
}

// Error level log.
func (l *logger) Error(args ...interface{}) {
	l.zap.Error(args...)
}

// Fatal level log. Calls os.Exit
func (l *logger) Fatal(args ...interface{}) {
	l.zap.Fatal(args...)
}

// // Debug level log.
// func (l *logger) Debug(args ...interface{}) {
// 	l.zap.Debug(args...)
// }

// Returns a new logger.
// Caller is not responsible for closing files.
func New(opt Options) *logger {
	level := zap.NewAtomicLevel()
	fileEncoder, consoleEncoder := getEncoders(true)
	cores := []zapcore.Core{}

	for _, writer := range opt.Files {
		out := zapcore.AddSync(writer) //it is intelligent about sync stuff.
		cores = append(cores, zapcore.NewCore(fileEncoder, out, level))
	}
	if opt.Stdout != nil {
		out := zapcore.AddSync(opt.Stdout) // os.Stdout satisfies Sync() already but just in case
		cores = append(cores, zapcore.NewCore(consoleEncoder, out, level))
	}
	core := zapcore.NewTee(cores...)
	l := zap.New(core)

	sugaredLogger := l.Sugar()
	return &logger{
		zap: sugaredLogger,
	}
}

// Standart zap config for this app.
func getEncoders(development bool) (zapcore.Encoder, zapcore.Encoder) {

	cfg := zap.NewProductionEncoderConfig()
	if development {
		cfg = zap.NewDevelopmentEncoderConfig()
		cfg.EncodeCaller = zapcore.ShortCallerEncoder
		cfg.EncodeDuration = zapcore.MillisDurationEncoder
	}

	cfg.LevelKey = "l"
	cfg.TimeKey = "t"
	cfg.CallerKey = "c"
	cfg.MessageKey = "m"

	cfg.EncodeLevel = zapcore.LowercaseLevelEncoder
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	file := zapcore.NewJSONEncoder(cfg)

	cfg.LevelKey = "level"
	cfg.TimeKey = "timestamp"
	cfg.CallerKey = "caller"
	cfg.MessageKey = "message"

	cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	cfg.EncodeTime = func(t time.Time, pae zapcore.PrimitiveArrayEncoder) {
		pae.AppendString(t.Format("15:04:05"))
	}
	console := zapcore.NewConsoleEncoder(cfg)

	return file, console
}
