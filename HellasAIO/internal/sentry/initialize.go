package sentry

import (
	"github.com/HellasAIO/HellasAIO/internal/utils"
	"github.com/HellasAIO/HellasAIO/internal/version"
	"github.com/getsentry/sentry-go"
	"log"
	"time"
)

func Initialize() {
	var env string
	if utils.Debug {
		env = "environment"
	} else {
		env = "release"
	}

	err := sentry.Init(sentry.ClientOptions{
		Dsn:              "sentryurl",
		Environment:      env,
		Release:          "hellasaio@" + version.Version,
		Debug:            utils.Debug,
		TracesSampleRate: 0.01,
		AttachStacktrace: true,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}

	defer sentry.Recover()
	defer sentry.Flush(2 * time.Second)
}
