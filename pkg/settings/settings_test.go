package settings_test

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"

	"github.com/synapse-service/node-vs/pkg/settings"
)

type Caller struct {
	mu sync.RWMutex
	n  int
}

func (c *Caller) Call() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.n++
}
func (c *Caller) Number() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.n
}

func TestBasic(t *testing.T) {
	var (
		filename       = fmt.Sprintf("%ssettings_%d", os.TempDir(), time.Now().UnixNano())
		initialAddress = "0.0.0.0:8081"
		initialValue   = []byte("gateway_api:\n  address: " + initialAddress)
		changedAddress = "0.0.0.0:8080"
	)
	assert.NoError(t, os.WriteFile(filename, initialValue, 0644), "unexpected error on writing to temporary file")
	defer func() {
		assert.NoError(t, os.Remove(filename), "unexpected error on removing temporary file")
	}()

	s, err := settings.New(
		settings.WithFileName(filename),
	)
	assert.NoError(t, err, "unexpected error on creating cache")

	// Initial retrieve
	assert.Equal(t, initialAddress, s.Get().GatewayAPI.Address, "retrieved value must be the same as initial")

	// Update
	onUpdate := new(Caller)
	s.OnUpdate(onUpdate.Call)
	assert.Equal(t, 0, onUpdate.Number(), "onUpdate must not be called without updates")
	var value settings.Value
	value.GatewayAPI.Address = changedAddress
	assert.NoError(t, s.Update(value), "unexpected error on updating value")
	assert.Equal(t, changedAddress, s.Get().GatewayAPI.Address, "retrieved value must be the same as set")
	assert.Equal(t, 1, onUpdate.Number(), "onUpdate must not be called once")

	// Update file
	b, err := os.ReadFile(filename)
	assert.NoError(t, err, "unexpected error on reading file")
	var changed settings.Value
	assert.NoError(t, yaml.Unmarshal(b, &changed), "unexpected error on unmarshaling file")
	assert.Equal(t, changedAddress, changed.GatewayAPI.Address, "persisted to file value must be the same as set")
}
