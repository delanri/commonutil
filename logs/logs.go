package logs

import (
	"fmt"

	"go.uber.org/zap"
)

type (
	Logger interface {
		Info(...interface{})
		Infof(string, ...interface{})
		Debug(...interface{})
		Debugf(string, ...interface{})
		Error(...interface{})
		Errorf(string, ...interface{})
		Warning(...interface{})
		Warningf(string, ...interface{})
		Fatal(...interface{})
		Fatalf(string, ...interface{})
		Print(...interface{})
		Printf(string, ...interface{})
		Println(...interface{})
		Instance() interface{}
		AddDefault(string, interface{}) *logger
	}

	Level     string
	Formatter string

	Option struct {
		Level       Level
		LogFilePath string
		Formatter   Formatter
	}

	logger struct {
		instance *zap.SugaredLogger
	}
)

const (
	Info  Level = "INFO"
	Debug Level = "DEBUG"
	Error Level = "ERROR"

	JSONFormatter Formatter = "JSON"
	TextFormatter Formatter = "TEXT"
)

func (l *logger) Info(args ...interface{}) {
	l.instance.Info(args...)
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.instance.Infof(format, args...)
}

func (l *logger) Debug(args ...interface{}) {
	l.instance.Debug(args...)
}

func (l *logger) Debugf(format string, args ...interface{}) {
	l.instance.Debugf(format, args...)
}

func (l *logger) Error(args ...interface{}) {
	l.instance.Error(args...)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.instance.Errorf(format, args...)
}

func (l *logger) Warning(args ...interface{}) {
	l.instance.Warn(args...)
}

func (l *logger) Warningf(format string, args ...interface{}) {
	l.instance.Warnf(format, args...)
}

func (l *logger) Fatal(args ...interface{}) {
	l.instance.Fatal(args...)
}

func (l *logger) Fatalf(format string, args ...interface{}) {
	l.instance.Fatalf(format, args...)
}

func (l *logger) Print(args ...interface{}) {
	l.instance.Info(args...)
}

func (l *logger) Printf(format string, args ...interface{}) {
	l.instance.Infof(format, args...)
}

func (l *logger) Println(args ...interface{}) {
	l.instance.Infoln(args...)
}

func (l *logger) Instance() interface{} {
	return l.instance
}
func New(option *Option) (Logger, error) {
	instance, err := zap.NewProduction()
	if err != nil {
		instance.Error("logger initialization error", zap.Any("error", err))
		panic(fmt.Sprintf("logger initialization failed %v", err))
	}

	return &logger{instance.Sugar()}, nil
}

func DefaultLog() (Logger, error) {
	return New(&Option{
		Level:     Info,
		Formatter: TextFormatter,
	})
}

func (l *logger) AddDefault(key string, value interface{}) *logger {
	l.instance.With(key, value)
	return l
}
