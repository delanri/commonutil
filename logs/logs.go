package logs

import (
	"os"

	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
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
		Instance() interface{}
		WithFields(keyValues Fields) Logger
	}
	Fields map[string]interface{}

	Level     string
	Formatter string

	Option struct {
		Level       Level
		LogFilePath string
		Formatter   Formatter
	}

	logger struct {
		instance *logrus.Logger
	}

	logrusLogEntry struct {
		entry *logrus.Entry
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
	l.instance.Warning(args...)
}

func (l *logger) Warningf(format string, args ...interface{}) {
	l.instance.Warningf(format, args...)
}

func (l *logger) Fatal(args ...interface{}) {
	l.instance.Fatal(args...)
}

func (l *logger) Fatalf(format string, args ...interface{}) {
	l.instance.Fatalf(format, args...)
}

func (l *logger) Print(args ...interface{}) {
	l.instance.Print(args...)
}

func (l *logger) Printf(format string, args ...interface{}) {
	l.instance.Printf(format, args...)
}

func (l *logger) Instance() interface{} {
	return l.instance
}

func New(option *Option) (Logger, error) {
	instance := logrus.New()

	if option.Level == Info {
		instance.Level = logrus.InfoLevel
	}

	if option.Level == Debug {
		instance.Level = logrus.DebugLevel
	}

	if option.Level == Error {
		instance.Level = logrus.ErrorLevel
	}

	var formatter logrus.Formatter

	if option.Formatter == JSONFormatter {
		formatter = &logrus.JSONFormatter{}
	} else {
		formatter = &logrus.TextFormatter{}
	}

	instance.Formatter = formatter

	// - check if log file path does exists
	if option.LogFilePath != "" {
		if _, err := os.Stat(option.LogFilePath); os.IsNotExist(err) {
			if _, err = os.Create(option.LogFilePath); err != nil {
				return nil, errors.Wrapf(err, "failed to create log file %s", option.LogFilePath)
			}
		}
		maps := lfshook.PathMap{
			logrus.InfoLevel:  option.LogFilePath,
			logrus.DebugLevel: option.LogFilePath,
			logrus.ErrorLevel: option.LogFilePath,
		}
		instance.Hooks.Add(lfshook.NewHook(maps, formatter))
	}

	return &logger{instance}, nil
}

func DefaultLog() (Logger, error) {
	return New(&Option{
		Level:     Info,
		Formatter: TextFormatter,
	})
}

func (l *logger) AddDefault(key string, value interface{}) *logger {
	l.instance.WithField(key, value)
	return l
}

func (l *logger) WithFields(fields Fields) Logger {
	return &logrusLogEntry{
		entry: l.instance.WithFields(convertToLogrusFields(fields)),
	}
}

func (l *logrusLogEntry) Panicf(format string, args ...interface{}) {
	l.entry.Fatalf(format, args...)
}

func (l *logrusLogEntry) WithFields(fields Fields) Logger {
	return &logrusLogEntry{
		entry: l.entry.WithFields(convertToLogrusFields(fields)),
	}
}

func (l *logrusLogEntry) Info(args ...interface{}) {
	l.entry.Info(args...)
}

func (l *logrusLogEntry) Infof(format string, args ...interface{}) {
	l.entry.Infof(format, args...)
}

func (l *logrusLogEntry) Debug(i ...interface{}) {
	l.entry.Debug(i...)
}

func (l *logrusLogEntry) Debugf(format string, args ...interface{}) {
	l.entry.Debugf(format, args...)
}

func (l *logrusLogEntry) Error(i ...interface{}) {
	l.entry.Error(i...)
}

func (l *logrusLogEntry) Errorf(format string, args ...interface{}) {
	l.entry.Errorf(format, args...)
}

func (l *logrusLogEntry) Warning(args ...interface{}) {
	l.entry.Warn(args...)
}

func (l *logrusLogEntry) Warningf(format string, args ...interface{}) {
	l.entry.Warnf(format, args...)
}

func (l *logrusLogEntry) Fatal(args ...interface{}) {
	l.entry.Fatal(args...)
}

func (l *logrusLogEntry) Fatalf(format string, args ...interface{}) {
	l.entry.Fatalf(format, args...)
}

func (l *logrusLogEntry) Print(args ...interface{}) {
	l.entry.Print(args...)
}

func (l *logrusLogEntry) Printf(format string, args ...interface{}) {
	l.entry.Printf(format, args...)
}

func (l *logrusLogEntry) Instance() interface{} {
	panic("implement me")
}

func convertToLogrusFields(fields Fields) logrus.Fields {
	logrusFields := logrus.Fields{}
	for index, val := range fields {
		logrusFields[index] = val
	}
	return logrusFields
}