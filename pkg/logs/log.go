package log

import (
	syslog "log"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

const (
	msgfield = "msg"
)

// Level represents a logging filter level
type Level string

var (
	// LevelNone No logging
	LevelNone Level = "none"
	// LevelError Only Error level events
	LevelError Level = "error"
	// LevelWarn Only Warn | Error level events
	LevelWarn Level = "warn"
	// LevelInfo Only Info | Warn | Error level events
	LevelInfo Level = "info"
	// LevelDebug Only Debug | Info | Warn | Error level events
	LevelDebug Level = "debug"
)

// Log is a light wrapper over go-kit logging, to simplify use and reduce duplication in calls.
type Log struct {
	logger log.Logger
}

// NewLog returns a new Log instance using a default go-kit logfmt logger on stderr
func NewLog() *Log {
	w := log.NewSyncWriter(os.Stdout)
	logger := log.NewLogfmtLogger(w)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.Caller(6))
	option := level.AllowDebug()
	logger = level.NewFilter(logger, option)
	return &Log{logger: logger}
}

// NewLogFromLogger returns a new Log instance using the provided go-kit logger instance
func NewLogFromLogger(logger log.Logger) *Log {
	return &Log{logger: logger}
}

// Debug log event
func (l *Log) Debug(msg string, keyvalues ...interface{}) {
	kvs := getKeyValues(msg, keyvalues)
	_ = level.Debug(l.logger).Log(kvs...)
}

// Info log event
func (l *Log) Info(msg string, keyvalues ...interface{}) {
	kvs := getKeyValues(msg, keyvalues)
	_ = level.Info(l.logger).Log(kvs...)
}

// Warn log event
func (l *Log) Warn(msg string, keyvalues ...interface{}) {
	kvs := getKeyValues(msg, keyvalues)
	_ = level.Warn(l.logger).Log(kvs...)
}

// Error log event
func (l *Log) Error(msg string, keyvalues ...interface{}) {
	kvs := getKeyValues(msg, keyvalues)
	_ = level.Error(l.logger).Log(kvs...)
}

func (l *Log) Fatal(msg string, keyvalues ...interface{}) {
	kvs := getKeyValues(msg, keyvalues)
	_ = level.Error(l.logger).Log(kvs...)
	syslog.Fatal(msg)
}

// With creates a log instance from the current one with some default k/v fields
func (l *Log) With(keyvalues ...interface{}) *Log {
	logger := log.With(l.logger, keyvalues...)
	return NewLogFromLogger(logger)
}

// GetLevelFilteredLogger creates a new log instance, but one that filters events to the level specified.
func GetLevelFilteredLogger(l *Log, logLevel Level) *Log {
	option := level.AllowAll()
	switch logLevel {
	case LevelDebug:
		option = level.AllowDebug()
	case LevelInfo:
		option = level.AllowInfo()
	case LevelWarn:
		option = level.AllowWarn()
	case LevelError:
		option = level.AllowError()
	case LevelNone:
		option = level.AllowNone()
	}

	logger := level.NewFilter(l.logger, option)

	return NewLogFromLogger(logger)
}

func getKeyValues(msg string, keyvalues []interface{}) []interface{} {
	n := len(keyvalues) + 2
	if len(keyvalues)%2 != 0 {
		n++
	}
	kvs := make([]interface{}, 0, n)
	kvs = append(kvs, msgfield)
	kvs = append(kvs, msg)
	kvs = append(kvs, keyvalues...)

	return kvs
}
