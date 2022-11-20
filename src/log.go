package src

import (
	"fmt"
	"os"
	"path"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type TMessageFLoggerOption struct {
	LogDirector string
}

type TMessageFLogger struct {
	option *TMessageFLoggerOption
	*logrus.Logger
	logmutex *sync.RWMutex
	file     *os.File
}

func NewTMessageFLogger(option TMessageFLoggerOption) (tl *TMessageFLogger) {
	tl = &TMessageFLogger{
		Logger:   logrus.New(),
		option:   &option,
		logmutex: &sync.RWMutex{},
	}
	tl.Logger.SetFormatter(&logrus.JSONFormatter{})
	tl.Logger.SetOutput(tl)
	return
}

func (t *TMessageFLogger) GetCurrentLogPath() string {
	return path.Join(t.option.LogDirector, fmt.Sprintf("%s_log", time.Now().Format("2006-01-02")))
}

func (t *TMessageFLogger) Fabric() (err error) {
	if t.option.LogDirector != "" {
		if _, err := os.Stat(t.option.LogDirector); err != nil {
			if err := os.Mkdir(t.option.LogDirector, os.ModePerm); err != nil {
				return err
			}
		}
	}
	if t.file, err = os.OpenFile(t.GetCurrentLogPath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666); err == nil {
		time.AfterFunc(howLongToTomorrow(), func() {
			t.Warning("log file updating process. Access restricted temporarily.")
			t.logmutex.Lock()
			t.file.Close()
			defer t.logmutex.Unlock()

			if err := t.Fabric(); err != nil {
				t.file = os.Stdout
				fmt.Fprintf(t.file, "\033[31m!!!! LOGOUT CHANGED AS A STDOUT PLEASE REPORT THIS !!!!!\033[0m")
				go t.Errorf("the log file couldn't be updated (%s) err= %v", t.GetCurrentLogPath(), err)
			}

		})
	}

	return

}

func (t TMessageFLogger) Write(bin []byte) (size int, err error) {
	fmt.Printf("%s", bin) // TODO: DEL? Maybe

	t.logmutex.Lock()
	defer t.logmutex.Unlock()

	size, err = t.file.Write(bin)
	return

}
