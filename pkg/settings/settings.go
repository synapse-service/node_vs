package settings

import (
	"os"
	"sync"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func New(options ...option) (*Settings, error) {
	var s Settings

	for _, option := range options {
		if err := option(&s); err != nil {
			return nil, errors.Wrap(err, "apply option")
		}
	}

	if s.filename == "" {
		return nil, errors.New("empty filename")
	}

	b, err := os.ReadFile(s.filename)
	if err != nil {
		return nil, errors.Wrap(err, "read file")
	}
	if err := yaml.Unmarshal(b, &s.value); err != nil {
		return nil, errors.Wrap(err, "decode settings")
	}

	return &s, nil
}

type Settings struct {
	mu                sync.RWMutex
	filename          string
	value             Value
	onUpdateCallbacks []func()
}

func (s *Settings) Get() Value {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.value
}

func (s *Settings) Update(value Value) error {
	defer func() {
		for _, callback := range s.onUpdateCallbacks {
			callback()
		}
	}()

	s.mu.Lock()
	defer s.mu.Unlock()
	s.value = value

	b, err := yaml.Marshal(s.value)
	if err != nil {
		return errors.Wrap(err, "decode settings")
	}
	if err := os.WriteFile(s.filename, b, 0644); err != nil {
		return errors.Wrap(err, "read file")
	}
	return nil
}

func (s *Settings) OnUpdate(callback func()) {
	s.onUpdateCallbacks = append(s.onUpdateCallbacks, callback)
}
