package util

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

var (
	once     sync.Once
	instance *Logger
)

func GetLogger() *Logger {
	once.Do(func() {
		instance = &Logger{logrus.New()}

		customFormatter := &logrus.TextFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
			FullTimestamp:   true,
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				return "", fmt.Sprintf("[%s:%d]", f.File, f.Line)
			},
		}

		instance.SetFormatter(customFormatter)
		instance.SetLevel(logrus.DebugLevel)
		instance.SetReportCaller(true)

		instance.Debug("logger initialized")
	})

	return instance
}
