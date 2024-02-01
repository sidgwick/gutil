package logger

import "github.com/sirupsen/logrus"

var logger = logrus.New()

func ResetLogger(_logger *logrus.Logger) {
	logger = _logger
}

func Tracef(format string, args ...interface{}) {
	logger.Tracef(format, args...)
}
