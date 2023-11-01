package ghttp

import "github.com/sirupsen/logrus"

var logger = logrus.New()

func ResetLogger(_logger *logrus.Logger) {
	logger = _logger
}
