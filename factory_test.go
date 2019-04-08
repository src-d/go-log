package log

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/x-cray/logrus-prefixed-formatter"
)

func TestLoggerFactoryNew_TextWithForce(t *testing.T) {
	require := require.New(t)

	f := &LoggerFactory{Format: TextFormat, ForceFormat: true}
	l, err := f.New(nil)
	require.NoError(err)

	logger, ok := l.(*logger)
	require.True(ok)
	require.IsType(&prefixed.TextFormatter{}, logger.Entry.Logger.Formatter)
}

func TestLoggerFactoryNew_JSON(t *testing.T) {
	require := require.New(t)

	f := &LoggerFactory{Format: JSONFormat, Level: InfoLevel}
	l, err := f.New(nil)
	require.NoError(err)

	logger, ok := l.(*logger)
	require.True(ok)
	require.IsType(&logrus.JSONFormatter{}, logger.Entry.Logger.Formatter)
	require.Equal(logrus.InfoLevel, logger.Entry.Logger.Level)
}

func TestLoggerFactoryNew_fluentd(t *testing.T) {
	require := require.New(t)

	f := &LoggerFactory{Format: FluentdFormat, Level: InfoLevel}
	l, err := f.New(nil)
	require.NoError(err)

	logger, ok := l.(*logger)
	require.True(ok)

	buf := &bytes.Buffer{}
	logger.Entry.Logger.Out = buf
	logger.Warn("something")

	var data map[string]interface{}
	json.Unmarshal(buf.Bytes(), &data)

	message, ok := data["message"]
	require.True(ok)
	require.Equal("something", message)

	severity, ok := data["severity"]
	require.True(ok)
	require.Equal("warning", severity)
}

func TestLoggerFactoryNew_NewFields(t *testing.T) {
	require := require.New(t)

	f := &LoggerFactory{Format: TextFormat, Level: DebugLevel}
	l, err := f.New(Fields{"foo": "bar"})
	require.NoError(err)

	logger, ok := l.(*logger)
	require.True(ok)
	require.Equal(logrus.DebugLevel, logger.Entry.Logger.Level)
	require.Equal(logrus.Fields{"foo": "bar"}, logger.Entry.Data)

}

func TestLoggerFactoryNew_Fields(t *testing.T) {
	require := require.New(t)

	js := `{"foo":"bar"}`
	f := &LoggerFactory{Format: TextFormat, Level: DebugLevel, Fields: js}
	l, err := f.New(nil)
	require.NoError(err)

	logger, ok := l.(*logger)
	require.True(ok)
	require.Equal(logrus.DebugLevel, logger.Entry.Logger.Level)
	require.Equal(logrus.Fields{"foo": "bar"}, logger.Entry.Data)

}

func TestLoggerFactoryNew_Error(t *testing.T) {
	require := require.New(t)

	// invalid level
	f := &LoggerFactory{Level: "text"}
	_, err := f.New(nil)
	require.Error(err)

	// invalid format
	f = &LoggerFactory{Level: InfoLevel, Format: "qux"}
	_, err = f.New(nil)
	require.Error(err)

	// invalid json
	f = &LoggerFactory{Level: InfoLevel, Format: TextFormat, Fields: "qux"}
	_, err = f.New(nil)
	require.Error(err)
}

func TestLoggerFactoryApply(t *testing.T) {
	require := require.New(t)

	f := &LoggerFactory{Format: TextFormat, ForceFormat: true, Level: DebugLevel}
	err := f.ApplyToLogrus()
	require.NoError(err)

	require.IsType(&prefixed.TextFormatter{}, logrus.StandardLogger().Formatter)
	require.Equal(logrus.DebugLevel, logrus.StandardLogger().Level)
}
