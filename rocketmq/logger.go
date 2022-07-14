package rocketmq

import (
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

type loggerWrap struct {
	logger *logrus.Logger
}

func (l *loggerWrap) Debug(msg string, fields map[string]interface{}) {
	if msg == "" && len(fields) == 0 {
		return
	}
	l.logger.WithField("component", "[ROCKETMQ]").WithFields(fields).Debug(msg)
}

func (l *loggerWrap) Info(msg string, fields map[string]interface{}) {
	if msg == "" && len(fields) == 0 {
		return
	}
	l.logger.WithField("component", "[ROCKETMQ]").WithFields(fields).Info(msg)
}

func (l *loggerWrap) Warning(msg string, fields map[string]interface{}) {
	if msg == "" && len(fields) == 0 {
		return
	}
	l.logger.WithField("component", "[ROCKETMQ]").WithFields(fields).Warning(msg)
}

func (l *loggerWrap) Error(msg string, fields map[string]interface{}) {
	if msg == "" && len(fields) == 0 {
		return
	}
	l.logger.WithField("component", "[ROCKETMQ]").WithFields(fields).WithFields(fields).Error(msg)
}

func (l *loggerWrap) Fatal(msg string, fields map[string]interface{}) {
	if msg == "" && len(fields) == 0 {
		return
	}
	l.logger.WithField("component", "[ROCKETMQ]").WithFields(fields).Fatal(msg)
}

func (l *loggerWrap) Level(level string) {
	switch strings.ToLower(level) {
	case "debug":
		l.logger.SetLevel(logrus.DebugLevel)
	case "warn":
		l.logger.SetLevel(logrus.WarnLevel)
	case "error":
		l.logger.SetLevel(logrus.ErrorLevel)
	default:
		l.logger.SetLevel(logrus.InfoLevel)
	}
}

func (l *loggerWrap) OutputPath(path string) (err error) {
	var file *os.File
	file, err = os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return
	}

	l.logger.Out = file
	return
}
