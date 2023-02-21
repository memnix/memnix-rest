package infrastructures

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/memnix/memnix-rest/config"
	"github.com/pkg/errors"
)

func getSentryDSN() string {
	return config.EnvHelper.GetEnv("SENTRY_DSN")
}

func ConnectSentry() error {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              getSentryDSN(),
		TracesSampleRate: 1.0,
		EnableTracing:    true,
		BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			if hint.Context != nil {
				if c, ok := hint.Context.Value(sentry.RequestContextKey).(*fiber.Ctx); ok {
					// You have access to the original Context if it panicked
					fmt.Println(utils.CopyString(c.Hostname()))
				}
			}
			fmt.Println(event)
			return event
		},
		Debug:            true,
		AttachStacktrace: true,
	})
	if err != nil {
		return errors.Wrap(err, "error initializing sentry")
	}

	return nil
}
