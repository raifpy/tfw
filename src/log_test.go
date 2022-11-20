package src

import (
	"testing"
)

func TestTMessageFLog(t *testing.T) {

	logger := NewTMessageFLogger(TMessageFLoggerOption{
		LogDirector: "debug_log",
	})
	if err := logger.Fabric(); err != nil {
		t.Fatal(err)
	}

	logger.Info("Merhaba Dünya!")
}
