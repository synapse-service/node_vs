package plugin_test

import (
	"log"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/synapse-service/node-vs/pkg/plugin"
)

// Generate `sample_plugin` for tests (use `go generate ./...`)

func TestBasic(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	exitCh := make(chan struct{})
	eventCallback, exitCallback := &checkFunc{}, &checkFunc{f: func() { exitCh <- struct{}{} }}
	p, err := plugin.New("../../sample_plugin",
		plugin.WithSettings(plugin.Settings{ID: "007"}),
		plugin.WithEventCallback(func(plugin.Event) { eventCallback.Run() }),
		plugin.WithExitCallback(func() { exitCallback.Run() }),
	)
	assert.NoError(t, err, "new plugin")
	assert.NoError(t, p.Kill(), "no errors on kill")
	<-exitCh
	assert.Equal(t, 0, eventCallback.Number(), "event callback never called")
	assert.Equal(t, 1, exitCallback.Number(), "exit callback called once")
}

type checkFunc struct {
	sync.RWMutex
	f func()
	n int
}

func (f *checkFunc) Run() {
	f.Lock()
	defer f.Unlock()
	f.n++
	if f.f != nil {
		f.f()
	}
}

func (f *checkFunc) Number() int {
	f.RLock()
	defer f.RUnlock()
	return f.n
}
