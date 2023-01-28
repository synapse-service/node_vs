package main

import (
	"flag"
	l "log"
	"os"

	"github.com/242617/core/application"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"

	"github.com/synapse-service/node-vs/pkg/settings"
	"github.com/synapse-service/node-vs/service"
)

type Config struct {
	Service service.Config `yaml:"service"`
}

func init() { l.SetFlags(l.Lshortfile | l.Ltime) }
func main() {
	configPath := flag.String("config", "/etc/service/config.yaml", "Config file path")
	flag.Parse()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05.000"})

	var cfg Config
	if b, err := os.ReadFile(*configPath); err != nil {
		l.Fatal(errors.Wrap(err, "load config"))
	} else if err := yaml.Unmarshal(b, &cfg); err != nil {
		l.Fatal(errors.Wrap(err, "unmarshal config"))
	}

	settings, err := settings.New(settings.WithFileName(cfg.Service.Settings))
	if err != nil {
		l.Fatal(errors.Wrap(err, "create settings"))
	}

	service, err := service.New(
		service.WithConfig(cfg.Service),
		service.WithSettings(settings),
	)
	if err != nil {
		l.Fatal(errors.Wrap(err, "create service"))
	}

	app, err := application.New(
		application.WithComponents(
			application.NewLifecycleComponent("service", service),
		),
	)
	if err != nil {
		l.Fatal(errors.Wrap(err, "create application"))
	}
	if err := app.Run(); err != nil {
		l.Fatal(errors.Wrap(err, "run application"))
	}
}
