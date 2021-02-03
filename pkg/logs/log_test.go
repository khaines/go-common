package log

import (
	"bytes"
	"io"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/stretchr/testify/assert"
)

func TestBasicLogOutput(t *testing.T) {

	buf := new(bytes.Buffer)

	log := createBufferedLogger(buf)

	log.Debug("hello")

	// get the log entry
	content := buf.String()

	assert.Equal(t, "level=debug msg=hello\n", content)
}

func createBufferedLogger(stream io.Writer) *Log {
	w := log.NewSyncWriter(stream)
	logger := log.NewLogfmtLogger(w)
	option := level.AllowDebug()
	logger = level.NewFilter(logger, option)
	return &Log{logger: logger}
}
