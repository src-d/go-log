package log

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	require := require.New(t)

	os.Setenv("LOG_LEVEL", "DEBUG")

	l, err := New()
	require.NoError(err)

	logger, ok := l.(*logger)
	require.True(ok)
	require.Equal(logrus.DebugLevel, logger.Entry.Logger.Level)
}
