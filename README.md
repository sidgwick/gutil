# Golang util set

## Init logger example

```go
package main

import (
	"context"

	"github.com/sidgwick/gutil/logger"
	"github.com/sirupsen/logrus"
)


const (
	TRACE_KEY = "trace-id"
)

type logContext struct {
}

func (l logContext) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (l logContext) Fire(entry *logrus.Entry) error {
	if entry.Context != nil {
		entry.Data[TRACE_KEY] = entry.Context.Value(TRACE_KEY)
	}

	return nil
}

func initLogger() {
	_logger := logrus.New()

	_logger.SetLevel(logrus.TraceLevel)
	_logger.SetFormatter(&logrus.JSONFormatter{})

	_logger.AddHook(logContext{})

	logger.ResetLogger(_logger)
}

func main() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, TRACE_KEY, "this-is-trace-id")

	initLogger()
	logger.Tracef(ctx, "log output testing")
}
```
