package process

import (
	"context"
	"log"
	"time"
)

type Process = func(context.Context) error

func New() *Processes {
	ctx, cancel := context.WithCancel(context.Background())
	return &Processes{
		ctx:    ctx,
		cancel: cancel,
		infos:  map[string]processInfo{},
	}
}

type Processes struct {
	ctx    context.Context
	cancel context.CancelFunc
	infos  map[string]processInfo
}

type processInfo struct {
	name      string
	class     string
	process   Process
	cancel    context.CancelFunc
	startedAt time.Time
}

func (ps *Processes) Start(name, class string, process Process) {
	ctx, cancel := context.WithCancel(ps.ctx)
	info := processInfo{
		name:      name,
		class:     class,
		process:   process,
		cancel:    cancel,
		startedAt: time.Now(),
	}
	ps.infos[name] = info
	go func(ctx context.Context, info processInfo) {
		if err := process(ctx); err != nil {
			log.Println(err)
		}
	}(ctx, info)
}

func (ps *Processes) Stop(names ...string) {
	m := map[string]struct{}{}
	for _, name := range names {
		m[name] = struct{}{}
	}
	for _, info := range ps.infos {
		if _, ok := m[info.name]; ok {
			info.cancel()
			delete(ps.infos, info.name)
		}
	}
}

func (ps *Processes) Filter(class string) []string {
	var names []string
	for _, info := range ps.infos {
		if info.class == class {
			names = append(names, info.name)
		}
	}
	return names
}

func (ps *Processes) StopAll() { ps.cancel() }
