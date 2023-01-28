package plugin

import (
	"bufio"
	"errors"
	"io"
	"io/fs"

	"github.com/242617/other/dispatcher"
	"google.golang.org/protobuf/proto"

	"github.com/synapse-service/node-vs/proto/plugin"
)

// Writing to stdin
func (p *Plugin) handleStdin(stdin io.Writer) {
	w := NewWriter(stdin)
	for b := range p.stdinCh {
		if _, err := w.Write(b); err != nil {
			p.Error().Err(err).Msg("write error")
		}
	}
}

// Reading from stdout
func (p *Plugin) handleStdout(stdout io.ReadCloser) {
	r := NewReader(stdout)
	for {
		b, err := r.ReadAll()
		if err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, fs.ErrClosed) {
				p.Warn().Err(err).Msg("graceful")
				break
			}
			p.Error().Err(err).Msg("read header from stdout")
			continue
		}
		p.stdoutCh <- b
	}
	close(p.stdoutCh)
}

func (p *Plugin) handleBytes() {
	for b := range p.stdoutCh {
		var message plugin.Output
		if err := proto.Unmarshal(b, &message); err != nil {
			p.Error().Err(err).Msg("unmarshal message from plugin")
			continue
		}

		switch data := message.Data.(type) {
		case *plugin.Output_Event:
			p.Dispatch(
				dispatcher.Event{
					Type:  EventTypeEvent,
					Value: Event{Type: data.Event.Type.String()},
				},
			)

		}
	}
}

// Reading from stderr for debug purpose
func (p *Plugin) handleStderr(stderr io.Reader) {
	s := bufio.NewScanner(stderr)
	for s.Scan() {
		p.Debug().Str("stderr", s.Text()).Msg("plugin")
	}
}
