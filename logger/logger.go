package logger

import (
	"context"

	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

func ResetLogger(_logger *logrus.Logger) {
	logger = _logger
}

func Tracef(ctx context.Context, format string, args ...interface{}) {
	logger.WithContext(ctx).Tracef(format, args...)
}
