package log

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/x-cray/logrus-prefixed-formatter"
)

const (
	// DefaultLevel is the level used by LoggerFactory when Level is omitted.
	DefaultLevel = "info"
	// DefaultFormat is the format used by LoggerFactory when Format is omitted.
	DefaultFormat = "text"
)

var (
	validLevels = map[string]bool{
		"info": true, "debug": true, "warning": true, "error": true,
	}
	validFormats = map[string]bool{
		"text": true, "json": true,
	}
)

// LoggerFactory is a logger factory used to instanciate new Loggers, from
// string configuration, mainly coming from console flags.
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
	lvl := DefaultLevel
	if f.Level != "" {
		lvl = strings.ToLower(f.Level)
	}

	if !validLevels[lvl] {
		return fmt.Errorf(
			"invalid level %s, valid levels are: %v", lvl, getKeysFromMap(validLevels),
		)
	}

	level, err := logrus.ParseLevel(lvl)
	if err != nil {
		return err
	}

	l.Level = level
	return nil
}

func (f LoggerFactory) setFormat(l *logrus.Logger) error {
	format := DefaultFormat
	if f.Format != "" {
		format = strings.ToLower(f.Format)
	}

	if !validFormats[format] {
		return fmt.Errorf(
			"invalid format %s, valid formats are: %v", format, getKeysFromMap(validFormats),
		)
	}

	switch format {
	case "text":
		f := new(prefixed.TextFormatter)
		f.ForceColors = true
		f.FullTimestamp = true
		l.Formatter = f
	case "json":
		l.Formatter = new(logrus.JSONFormatter)
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

func getKeysFromMap(m map[string]bool) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}

	return keys
}
