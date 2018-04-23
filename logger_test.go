package log

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestLoggerNew(t *testing.T) {
	require := require.New(t)

	f := &LoggerFactory{Format: "text", Level: "debug"}
	l, err := f.New()
	require.NoError(err)

	l = l.New(Fields{"foo": "qux"})
	l1, ok := l.(*logger)
	require.True(ok)
	require.Equal(logrus.Fields{"foo": "qux"}, l1.Entry.Data)

	l = l.New(Fields{"bar": "baz"})
	l2, ok := l.(*logger)
	require.True(ok)
	require.Equal(logrus.Fields{
		"foo": "qux",
		"bar": "baz",
	}, l2.Entry.Data)
}
