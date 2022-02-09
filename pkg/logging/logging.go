package logging

import (
	"fmt"
	"runtime"

	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.DebugLevel)
	formatter := &logrus.TextFormatter{
		TimestampFormat: "02-01-2006 15:04:05",
		FullTimestamp:   true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			return "", fmt.Sprintf("")
		},
	}
	logrus.SetFormatter(formatter)
}
