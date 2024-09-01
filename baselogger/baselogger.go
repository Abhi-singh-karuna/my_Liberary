package baselogger

import (
	"errors"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// BaseLogger
type BaseLogger struct {
	sugarLogger *zap.SugaredLogger
	level       zap.AtomicLevel
}

// BaseLogger constructor
func NewBaseLogger() (bl *BaseLogger) {
	bl = &BaseLogger{}
	bl.init()
	return
}

// For mapping config logger to app logger levels
var levels = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

// Init BaseLogger
func (bl *BaseLogger) init() {
	encoderCfg := zap.NewDevelopmentEncoderConfig()
	encoderCfg.LevelKey = "LEVEL"
	encoderCfg.CallerKey = "CALLER"
	encoderCfg.TimeKey = "TIME"
	encoderCfg.NameKey = "NAME"
	encoderCfg.MessageKey = "MESSAGE"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	encoder := zapcore.NewJSONEncoder(encoderCfg)
	logWriter := zapcore.AddSync(os.Stderr)
	bl.level = zap.NewAtomicLevel()

	core := zapcore.NewCore(encoder, logWriter, bl.level)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	bl.sugarLogger = logger.Sugar()
}

func (bl *BaseLogger) SetLevel(level string) error {
	zLevel, found := levels[level]
	if found {
		bl.level.SetLevel(zLevel)
	} else {
		return errors.New("Wrong logger level")
	}
	return nil
}

// BaseLogger methods
func (bl *BaseLogger) Sync() error {
	return bl.sugarLogger.Sync()
}

func (bl *BaseLogger) Debug(args ...interface{}) {
	bl.sugarLogger.Debug(args...)
}

func (bl *BaseLogger) Debugf(template string, args ...interface{}) {
	bl.sugarLogger.Debugf(template, args...)
}

func (bl *BaseLogger) Info(args ...interface{}) {
	bl.sugarLogger.Info(args...)
}

func (bl *BaseLogger) Infof(template string, args ...interface{}) {
	bl.sugarLogger.Infof(template, args...)
}

func (bl *BaseLogger) Warn(args ...interface{}) {
	bl.sugarLogger.Warn(args...)
}

func (bl *BaseLogger) Warnf(template string, args ...interface{}) {
	bl.sugarLogger.Warnf(template, args...)
}

func (bl *BaseLogger) Error(args ...interface{}) {
	bl.sugarLogger.Error(args...)
}

func (bl *BaseLogger) Errorf(template string, args ...interface{}) {
	bl.sugarLogger.Errorf(template, args...)
}

func (bl *BaseLogger) DPanic(args ...interface{}) {
	bl.sugarLogger.DPanic(args...)
}

func (bl *BaseLogger) DPanicf(template string, args ...interface{}) {
	bl.sugarLogger.DPanicf(template, args...)
}

func (bl *BaseLogger) Panic(args ...interface{}) {
	bl.sugarLogger.Panic(args...)
}

func (bl *BaseLogger) Panicf(template string, args ...interface{}) {
	bl.sugarLogger.Panicf(template, args...)
}

func (bl *BaseLogger) Fatal(args ...interface{}) {
	bl.sugarLogger.Fatal(args...)
}

func (bl *BaseLogger) Fatalf(template string, args ...interface{}) {
	bl.sugarLogger.Fatalf(template, args...)
}
