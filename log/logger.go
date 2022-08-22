package log

import (
	"fmt"
	"io"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var zapDefault, _ = zap.NewProduction()
var l Logger = zapDefault.Sugar()

var logFileHook *lumberjack.Logger

type GinRecovery struct {
	Writers []io.Writer
}

type Logger interface {
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Debugw(format string, keysAndValues ...interface{})

	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Infow(format string, keysAndValues ...interface{})

	Warn(v ...interface{})
	Warnf(format string, v ...interface{})
	Warnw(format string, keysAndValues ...interface{})

	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Errorw(format string, keysAndValues ...interface{})

	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Fatalw(format string, keysAndValues ...interface{})

	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
	Panicw(format string, keysAndValues ...interface{})

	Sync() error
	With(args ...interface{}) *zap.SugaredLogger
}

func Debug(v ...interface{}) {
	l.Debug(v...)
}
func Debugf(format string, v ...interface{}) {
	l.Debugf(format, v...)
}
func Debugw(format string, keysAndValues ...interface{}) {
	l.Debugw(format, keysAndValues...)
}

func Info(v ...interface{}) {
	l.Info(v...)
}
func Infof(format string, v ...interface{}) {
	l.Infof(format, v...)
}
func Infow(format string, keysAndValues ...interface{}) {
	l.Infow(format, keysAndValues...)
}

func Warn(v ...interface{}) {
	l.Warn(v...)
}
func Warnf(format string, v ...interface{}) {
	l.Warnf(format, v...)
}
func Warnw(format string, keysAndValues ...interface{}) {
	l.Warnw(format, keysAndValues...)
}

func Error(v ...interface{}) {
	l.Error(v...)
}
func Errorf(format string, v ...interface{}) {
	l.Errorf(format, v...)
}
func Errorw(format string, keysAndValues ...interface{}) {
	l.Errorw(format, keysAndValues...)
}

func Fatal(v ...interface{}) {
	l.Fatal(v...)
}
func Fatalf(format string, v ...interface{}) {
	l.Fatalf(format, v...)
}
func Fatalw(format string, keysAndValues ...interface{}) {
	l.Fatalw(format, keysAndValues...)
}

func Panic(v ...interface{}) {
	l.Panic(v...)
}
func Panicf(format string, v ...interface{}) {
	l.Panicf(format, v...)
}
func Panicw(format string, keysAndValues ...interface{}) {
	l.Panicw(format, keysAndValues...)
}

func Sync() error {
	return l.Sync()
}

func With(args ...interface{}) *zap.SugaredLogger {
	return l.With(args)
}

func GetLogFileHook() *lumberjack.Logger {
	return logFileHook
}

func SetLogger(logger Logger) {
	l = logger
}

func GetLogger() Logger {
	return l
}

type logConfig struct {
	LogFileName string `json:"log_file_name"` // 日志文件名称
	LogLevel    string `json:"log_level"`     // 日志等级
	Log         *Log   `json:"log"`
}

type Log struct {
	MaxSize    int  `json:"max_size"`    // MaxSize 是日志文件在轮换之前的最大大小（以 MB 为单位）。默认为 100MB
	MaxAge     int  `json:"max_age"`     // MaxAge 是根据文件名中编码的时间戳保留旧日志文件的最大天数，单位：天（也就是24h）
	MaxBackups int  `json:"max_backups"` // MaxBackups 是要保留的旧日志文件的最大数量。默认是保留所有旧的日志文件（尽管 MaxAge 可能仍会导致它们被删除）
	LocalTime  bool `json:"local_time"`  // LocalTime 确定用于格式化备份文件中时间戳的时间是否是计算机的本地时间。默认是使用UTC 时间
	Compress   bool `json:"compress"`    // Compress 确定是否应使用gzip 压缩旋转的日志文件。默认是不执行压缩
}

// NewLogger 初始化日志空间, logFileName 用于输出的日志文件名, logLevel 表示日志级别
func NewLogger(cfg *logConfig) {

	// Log config 的前置参数检查
	cfg.preCheck()

	// 初始化日志级别
	lLevel := StringLevel(cfg.LogLevel)

	logFileHook = &lumberjack.Logger{
		Filename:   fmt.Sprintf("logs/%s.log", cfg.LogFileName),
		MaxBackups: cfg.Log.MaxBackups,
		MaxSize:    cfg.Log.MaxSize,
		MaxAge:     cfg.Log.MaxAge,
		LocalTime:  cfg.Log.LocalTime,
		Compress:   cfg.Log.Compress,
	}

	consoleColoredEncoderConfig := zap.NewProductionEncoderConfig()
	consoleColoredEncoderConfig.TimeKey = "time"
	consoleColoredEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleColoredEncoderConfig.EncodeTime = func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}

	fileEncoderConfig := zap.NewProductionEncoderConfig()
	fileEncoderConfig.TimeKey = "time"
	fileEncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	fileEncoderConfig.EncodeTime = func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}

	filePriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= lLevel
	})
	errPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= lLevel && lvl >= zapcore.ErrorLevel
	})
	stdoutPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= lLevel && lvl < zapcore.ErrorLevel
	})

	consoleEncoder := zapcore.NewConsoleEncoder(consoleColoredEncoderConfig)
	fileEncoder := zapcore.NewConsoleEncoder(fileEncoderConfig)

	cores := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(logFileHook), filePriority),
		zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), stdoutPriority),
		zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stderr), errPriority),
	)

	logger := zap.New(cores, zap.AddStacktrace(zap.NewAtomicLevelAt(zap.ErrorLevel)))

	l = logger.Sugar()
}

func StringLevel(levelString string) zapcore.Level {
	level := zap.NewAtomicLevelAt(zapcore.DebugLevel)
	level.UnmarshalText([]byte(levelString))
	return level.Level()
}

func (mo *GinRecovery) Write(p []byte) (n int, err error) {
	for _, writer := range mo.Writers {
		n, err = writer.Write(p)
	}
	return
}

func (cfg *logConfig) preCheck() {
	if cfg.LogFileName == "" {
		cfg.LogFileName = "default"
	}
	if cfg.LogLevel == "" {
		cfg.LogLevel = "debug"
	}

	if cfg.Log == nil {
		cfg.Log = &Log{
			MaxSize:    100,   // 100MB
			MaxAge:     7,     // 7天
			MaxBackups: 10,    // 10个备份文件
			LocalTime:  true,  // 启用本地时间
			Compress:   false, // 不采用gzip 压缩
		}
	} else {
		if cfg.Log.MaxSize == 0 {
			cfg.Log.MaxSize = 100 // 100MB
		}
		if cfg.Log.MaxAge == 0 {
			cfg.Log.MaxAge = 7 // 7天
		}
		if cfg.Log.MaxBackups == 0 {
			cfg.Log.MaxBackups = 10 // 10个备份文件
		}
	}
}
