package log

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/x-cray/logrus-prefixed-formatter"
)

// Level represents the logging level.
type Level int32

const (
	// Info level. General operational entries about what's going on inside the
	// application.
	Info = Level(logrus.WarnLevel)
	// Debug level. Usually only enabled when debugging. Very verbose logging.
	Debug = Level(logrus.DebugLevel)
	// Warning level. Non-critical entries that deserve eyes.
	Warning = Level(logrus.WarnLevel)
	// Error level. Used for errors that should definitely be noted.
	Error = Level(logrus.ErrorLevel)
)

// Format represents the format of logs to be use.
type Format int32

const (
	// Text is the default human readable log format.
	Text Format = iota
	// JSON format, mainly for machine to machine comunication.
	JSON
)

// LoggerFactory is a logger factory used to instanciate new Loggers, from
// string configuration, mainly comming from console flags.
type LoggerFactory struct {
	// Leves as string, values are "info", "debug", "warning" or "error".
	Level string
	// Format as string, values are "text" or "json", by default "text" is used.
	Format string
	// Fields in JSON format to be used by configured in the new Logger.
	Fields string
}

// New returns a new logger based on the LoggerFactory values.
func (f LoggerFactory) New() (Logger, error) {
	l := logrus.New()
	if err := f.setLevel(l); err != nil {
		return nil, err
	}

	if err := f.setFormat(l); err != nil {
		return nil, err
	}

	return f.setFields(l)
}

// ApplyToLogrus configures the standard logrus Logger with the LoggerFactory
// values. Useful to propagate the configuration to third-party libraries using
// logrus.
func (f LoggerFactory) ApplyToLogrus() error {
	if err := f.setLevel(logrus.StandardLogger()); err != nil {
		return err
	}

	return f.setFormat(logrus.StandardLogger())
}

func (f LoggerFactory) setLevel(l *logrus.Logger) error {
	level, err := logrus.ParseLevel(f.Level)
	if err != nil {
		return err
	}

	l.Level = level
	return nil
}

func (f LoggerFactory) setFormat(l *logrus.Logger) error {
	switch f.Format {
	case "text":
		f := new(prefixed.TextFormatter)
		f.ForceColors = true
		f.FullTimestamp = true
		l.Formatter = f
	case "json":
		l.Formatter = new(logrus.JSONFormatter)
	default:
		return fmt.Errorf("unknown logger format: %q", f.Format)
	}

	return nil
}

func (f *LoggerFactory) setFields(l *logrus.Logger) (Logger, error) {
	var fields logrus.Fields
	if f.Fields != "" {
		if err := json.Unmarshal([]byte(f.Fields), &fields); err != nil {
			return nil, err
		}
	}

	e := l.WithFields(fields)
	return &logger{*e}, nil
}
