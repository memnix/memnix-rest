package infrastructures

import (
	"github.com/memnix/memnix-rest/config"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/pkg/errors"
)

var relicApp *newrelic.Application

// NewRelic creates a newrelic application
func NewRelic() error {
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName(config.EnvHelper.GetEnv("NEWRELIC_NAME")),
		newrelic.ConfigLicense(config.EnvHelper.GetEnv("NEWRELIC_KEY")),
		newrelic.ConfigAppLogForwardingEnabled(true),
	)
	if err != nil {
		return errors.Wrap(err, "error creating newrelic application")
	}

	relicApp = app

	return nil
}

// GetRelicApp returns the newrelic application
func GetRelicApp() *newrelic.Application {
	return relicApp
}
