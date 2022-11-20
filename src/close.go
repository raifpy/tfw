package src

import (
	"os"
	"os/signal"
	"syscall"
)

func (s *TMessageF) CloseWait() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM)

	<-ch
	s.Log.Warnf("SIGNAL TERM DETECTED. EXITING")
	os.Exit(1)
}
