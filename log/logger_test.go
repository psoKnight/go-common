package log

import (
	"testing"
)

func TestLoggerDefault(t *testing.T) {
	cfg := &logConfig{
		LogFileName: "",
		LogLevel:    "",
		Log:         nil,
	}

	// 初始化logger
	NewLogger(cfg)

	Info("This is a Info.")
	Infof("This is a Infof, str: %v.", "infof")

	Debug("This is a Debug.")
	Debugf("This is a Debugf, str: %v.", "debugf")

	Error("This is a Error.")
	Errorf("This is a Errorf, str: %v.", "errorf")

	//Fatal("This is a Fatal.")
	//Fatalf("This is a Fatalf, str: %v.", "fatalf")

	panic("This is a panic.")
	Panicf("This is a Panicf, str: %v.", "panicf")
}

func TestLogger(t *testing.T) {
	cfg := &logConfig{
		LogFileName: "go-common",
		LogLevel:    "debug",
		Log: &Log{
			MaxSize:    100,
			MaxAge:     7,
			MaxBackups: 7,
			LocalTime:  true,
			Compress:   true,
		},
	}

	// 初始化logger
	NewLogger(cfg)

	Info("This is a Info.")
	Infof("This is a Infof, str: %v.", "infof")

	Debug("This is a Debug.")
	Debugf("This is a Debugf, str: %v.", "debugf")

	Error("This is a Error.")
	Errorf("This is a Errorf, str: %v.", "errorf")

	//Fatal("This is a Fatal.")
	//Fatalf("This is a Fatalf, str: %v.", "fatalf")

	panic("This is a panic.")
	Panicf("This is a Panicf, str: %v.", "panicf")
}
