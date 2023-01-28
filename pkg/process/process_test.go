package process_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/synapse-service/node-vs/pkg/process"
)

type testProcess struct {
	sync.RWMutex
	finished bool
	err      error
}

func (p *testProcess) Err() error {
	p.RLock()
	err := p.err
	p.RUnlock()
	return err
}

func (p *testProcess) Finished() bool {
	p.RLock()
	finished := p.finished
	p.RUnlock()
	return finished
}

func (p *testProcess) Do(ctx context.Context) error {
	<-ctx.Done()
	err := ctx.Err()
	p.Lock()
	p.finished = true
	p.err = err
	p.Unlock()
	return err
}

func TestProcessBasic(t *testing.T) {
	ps := process.New()

	first := new(testProcess)
	second := new(testProcess)
	third := new(testProcess)
	fourth := new(testProcess)
	fifth := new(testProcess)

	ps.Start("first", "sample", first.Do)
	ps.Start("second", "sample", second.Do)
	ps.Start("third", "abra", third.Do)
	ps.Start("fourth", "cadabra", fourth.Do)
	ps.Start("fifth", "cadabra", fifth.Do)
	delay()

	names := ps.Filter("sample")
	assert.Equal(t, 2, len(names), "incorrect names length")
	ps.Stop(names...)
	delay()

	assert.True(t, first.Finished(), "first finished")
	assert.ErrorIs(t, first.Err(), context.Canceled, "unexpected first error")
	assert.True(t, second.Finished(), "second finished")
	assert.ErrorIs(t, second.Err(), context.Canceled, "unexpected second error")
	assert.False(t, third.Finished(), "third not finished")
	assert.NoError(t, third.Err(), "unexpected third error")

	ps.StopAll()
	delay()

	assert.True(t, third.Finished(), "third finished")
	assert.ErrorIs(t, third.Err(), context.Canceled, "unexpected third error")
	assert.True(t, fourth.Finished(), "fourth finished")
	assert.ErrorIs(t, fourth.Err(), context.Canceled, "unexpected fourth error")
	assert.True(t, fifth.Finished(), "fifth finished")
	assert.ErrorIs(t, fifth.Err(), context.Canceled, "unexpected fifth error")
}

func TestFilterProcessesAfterStop(t *testing.T) {
	ps := process.New()

	first := new(testProcess)
	second := new(testProcess)

	ps.Start("first", "sample", first.Do)
	ps.Start("second", "sample", second.Do)
	delay()

	ps.Stop("first")
	names := ps.Filter("sample")
	assert.Equal(t, 1, len(names), "length of filtered processes after removing")
}

func delay() { time.Sleep(10 * time.Millisecond) }
