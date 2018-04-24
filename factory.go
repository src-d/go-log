package log

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/onrik/logrus/filename"
	"github.com/sirupsen/logrus"
	"github.com/x-cray/logrus-prefixed-formatter"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	// DefaultLevel is the level used by LoggerFactory when Level is omitted.
	DefaultLevel = "info"
	// DefaultFormat is the format used by LoggerFactory when Format is omitted.
	DefaultFormat = "text"
	// DefaultTimeFormat is a handy timestamp (Jan _2 15:04:05.000000)
	// with microsecond precision
	DefaultTimeFormat = time.StampMicro
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
	// Level as string, values are "info", "debug", "warning" or "error".
	Level string
	// Format as string, values are "text" or "json", by default "text" is used.
	// when a terminal is not detected "json" is used instead.
	Format string
	// TimeFormat is used for marshaling timestamps,
	// by default:"Jan _2 15:04:05.000000"
	TimeFormat string
	// Fields in JSON format to be used by configured in the new Logger.
	Fields string
	// ForceFormat if true the fact of being in a terminal or not is ignored.
	ForceFormat bool
}

// New returns a new logger based on the LoggerFactory values.
func (f *LoggerFactory) New() (Logger, error) {
	l := logrus.New()
	f.setHook(l)

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
func (f *LoggerFactory) ApplyToLogrus() error {
	std := logrus.StandardLogger()
	f.setHook(std)

	if err := f.setLevel(std); err != nil {
		return err
	}

	return f.setFormat(std)
}

func (f *LoggerFactory) setLevel(l *logrus.Logger) error {
	if err := f.setDefaultLevel(); err != nil {
		return err
	}

	level, err := logrus.ParseLevel(f.Level)
	if err != nil {
		return err
	}

	l.Level = level
	return nil
}

func (f *LoggerFactory) setDefaultLevel() error {
	if f.Level == "" {
		f.Level = DefaultLevel
	}

	f.Level = strings.ToLower(f.Level)
	if validLevels[f.Level] {
		return nil
	}

	return fmt.Errorf(
		"invalid level %s, valid levels are: %v",
		f.Level, getKeysFromMap(validLevels),
	)
}

func (f *LoggerFactory) setFormat(l *logrus.Logger) error {
	if err := f.setDefaultFormat(); err != nil {
		return err
	}

	switch f.Format {
	case "text":
		fmt := new(prefixed.TextFormatter)
		fmt.ForceColors = true
		fmt.FullTimestamp = true
		fmt.TimestampFormat = f.TimeFormat
		l.Formatter = fmt
	case "json":
		fmt := new(logrus.JSONFormatter)
		fmt.TimestampFormat = f.TimeFormat
		l.Formatter = fmt
	}

	return nil
}

func (f *LoggerFactory) setDefaultFormat() error {
	if f.Format == "" {
		f.Format = DefaultFormat
	}

	if f.TimeFormat == "" {
		f.TimeFormat = DefaultTimeFormat
	}

	f.Format = strings.ToLower(f.Format)
	if validFormats[f.Format] {
		return nil
	}

	if !f.ForceFormat && isTerminal() {
		f.Format = "json"
	}

	return fmt.Errorf(
		"invalid format %s, valid formats are: %v",
		f.Format, getKeysFromMap(validFormats),
	)
}

func (f *LoggerFactory) setHook(l *logrus.Logger) {
	l.AddHook(filename.NewHook(
		logrus.DebugLevel,
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel),
	)
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

func isTerminal() bool {
	return terminal.IsTerminal(int(os.Stdout.Fd()))
}
