package service

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	l "github.com/rs/zerolog/log"

	"github.com/synapse-service/node-vs/pkg/process"
)

const (
	ProcessClassSystem = "system"
	ProcessClassCamera = "camera"
)

func New(options ...option) (*Service, error) {
	s := Service{
		log:       l.With().Str("component", "service").Logger(),
		processes: process.New(),
	}

	for _, option := range options {
		if err := option(&s); err != nil {
			return nil, errors.Wrap(err, "apply option")
		}
	}

	if s.settings == nil {
		return nil, errors.New("empty settings")
	}

	return &s, nil
}

type Service struct {
	log       zerolog.Logger
	processes *process.Processes
	cfg       Config
	settings  Settings
}

func (s *Service) Start(ctx context.Context) error {
	// Start camera processes
	s.settings.OnUpdate(s.updateCameraProcesses)
	s.updateCameraProcesses()
	// Start syncing to gateway
	s.processes.Start(ProcessNameSync, ProcessClassSystem, s.syncProcess)
	return nil
}

func (s *Service) Stop(ctx context.Context) error {
	s.processes.StopAll()
	time.Sleep(500 * time.Millisecond)
	return nil
}
