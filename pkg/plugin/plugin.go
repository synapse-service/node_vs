// Package plugin is a wrapper over process with communication methods.

package plugin

import (
	"os/exec"

	"github.com/242617/other/dispatcher"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	l "github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"

	"github.com/synapse-service/node-vs/proto/plugin"
)

func New(file string, options ...option) (*Plugin, error) {
	if file == "" {
		return nil, errors.New("empty file")
	}

	p := Plugin{
		Logger:          l.With().Str("component", "plugin").Logger(),
		EventDispatcher: dispatcher.NewEventDispatcher(),
		file:            file,
	}

	for _, option := range options {
		if err := option(&p); err != nil {
			return nil, errors.Wrap(err, "apply option")
		}
	}

	if p.settings.IsEmpty() {
		return nil, errors.New("empty settings")
	}

	// Settings

	s := plugin.Settings{
		Id: p.settings.ID,
	}
	b, err := proto.Marshal(&s)
	if err != nil {
		return nil, errors.Wrap(err, "marshal settings")
	}

	// Start

	p.cmd = exec.Command(p.file, string(b))

	stdin, err := p.cmd.StdinPipe()
	if err != nil {
		return nil, errors.Wrap(err, "create stdin pipe")
	}
	p.stdinCh = make(chan []byte)
	go p.handleStdin(stdin)

	stdout, err := p.cmd.StdoutPipe()
	if err != nil {
		return nil, errors.Wrap(err, "create stdout pipe")
	}
	p.stdoutCh = make(chan []byte)
	go p.handleStdout(stdout)
	go p.handleBytes()

	stderr, err := p.cmd.StderrPipe()
	if err != nil {
		return nil, errors.Wrap(err, "create stderr pipe")
	}
	go p.handleStderr(stderr)

	if err := p.cmd.Start(); err != nil {
		return nil, errors.Wrap(err, "start")
	}

	runCh := make(chan struct{})
	go func() {
		defer stdin.Close()
		defer stdout.Close()
		defer stderr.Close()
		defer func() { close(p.stdinCh) }()
		runCh <- struct{}{}
		if err := p.cmd.Wait(); err != nil {
			p.Error().Err(err).Msg("plugin run")
		}
		p.Dispatch(
			dispatcher.Event{
				Type: EventTypeExit,
			},
		)
	}()
	<-runCh

	return &p, nil
}

// Plugin is a OS process wrapper with convenient and simplified API
type Plugin struct {
	zerolog.Logger
	dispatcher.EventDispatcher
	file     string
	settings Settings
	cmd      *exec.Cmd
	stdinCh  chan []byte
	stdoutCh chan []byte
}

func (p *Plugin) Kill() error { return p.kill() }
func (p *Plugin) kill() error {
	if err := p.cmd.Process.Kill(); err != nil {
		return errors.Wrap(err, "kill process")
	}
	return nil
}
