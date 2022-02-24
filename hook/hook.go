package hook

import (
	"context"

	"github.com/newrelic/newrelic-client-go/pkg/config"
	"github.com/newrelic/newrelic-client-go/pkg/logs"
	"github.com/sirupsen/logrus"
)

const (
	BatchTimeout = 10
	BatchSize    = 20
)

type h struct {
	Client *logs.Logs
}

// Levels returns the log levels for the hook to be fired.
func (hk *h) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.DebugLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.InfoLevel,
	}
}

// Fire is triggered on every new log entry.
func (hk *h) Fire(l *logrus.Entry) error {
	if err := hk.Client.EnqueueLogEntry(l.Context, l); err != nil {
		return err
	}

	return nil
}

// Hook is triggered on every new log entry.
func Hook(key string, account int) logrus.Hook {
	cfg := config.New()
	cfg.LicenseKey = key
	cfg.LogLevel = "info"
	client := logs.New(cfg)

	err := client.BatchMode(
		context.Background(),
		account,
		logs.BatchConfigQueueSize(BatchSize),
		logs.BatchConfigTimeout(BatchTimeout),
	)

	if err != nil {
		panic(err)
	}

	return &h{
		Client: &client,
	}
}
