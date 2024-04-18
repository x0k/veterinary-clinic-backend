package module

import (
	"context"
	"log/slog"
	"sync"
	"sync/atomic"
)

type Module struct {
	name      string
	log       *slog.Logger
	wg        sync.WaitGroup
	services  []Service
	postStart []Hook
	preStop   []Hook
	fatal     chan error
	stopped   atomic.Bool
}

func New(log *slog.Logger, name string) *Module {
	return &Module{
		log:   log.With(slog.String("module_name", name)),
		name:  name,
		fatal: make(chan error, 1),
	}
}

func (m *Module) Name() string {
	return m.name
}

func (m *Module) Append(services ...Service) {
	m.services = append(m.services, services...)
}

func (m *Module) PostStart(hooks ...Hook) {
	m.postStart = append(m.postStart, hooks...)
}

func (m *Module) PreStop(hooks ...Hook) {
	m.preStop = append(m.preStop, hooks...)
}

func (m *Module) Fatal(ctx context.Context, err error) {
	if m.stopped.Swap(true) {
		m.log.LogAttrs(ctx, slog.LevelError, "fatal error", slog.String("error", err.Error()))
		return
	}
	m.fatal <- err
	close(m.fatal)
}

func (m *Module) awaiter(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return nil
	case err := <-m.fatal:
		return err
	}
}

func (m *Module) start(ctx context.Context, awaiter func(context.Context) error) error {
	if len(m.services) == 0 {
		return nil
	}

	if m.stopped.Load() {
		return <-m.fatal
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, service := range m.services {
		m.wg.Add(1)
		go func() {
			defer m.wg.Done()
			m.log.LogAttrs(ctx, slog.LevelInfo, "starting", slog.String("service", service.Name()))
			if err := service.Start(ctx); err != nil {
				m.Fatal(ctx, err)
			}
			m.log.LogAttrs(ctx, slog.LevelInfo, "stopped", slog.String("service", service.Name()))
		}()
	}

	for _, hook := range m.postStart {
		m.log.LogAttrs(ctx, slog.LevelInfo, "run post start", slog.String("hook", hook.Name()))
		if err := hook.Run(ctx); err != nil {
			m.Fatal(ctx, err)
		}
	}

	err := awaiter(ctx)

	for _, hook := range m.preStop {
		m.log.LogAttrs(ctx, slog.LevelInfo, "run pre stop", slog.String("hook", hook.Name()))
		if err := hook.Run(ctx); err != nil {
			m.Fatal(ctx, err)
		}
	}

	m.log.LogAttrs(ctx, slog.LevelInfo, "stopping")
	m.stopped.Store(true)
	cancel()

	m.wg.Wait()

	return err
}

func (m *Module) Start(ctx context.Context) error {
	return m.start(ctx, m.awaiter)
}
