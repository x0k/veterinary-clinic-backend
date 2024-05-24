//go:build js && wasm

package js_adapters

import (
	"context"
	"fmt"
	"syscall/js"

	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger"
	"github.com/x0k/veterinary-clinic-backend/internal/lib/logger/sl"
)

const name = "js_adapters.SimpleCache[T]"

type SimpleCacheConfig struct {
	Enabled bool      `js:"enabled"`
	Get     *js.Value `js:"get"`
	Add     *js.Value `js:"add"`
}

type SimpleCache[T any] struct {
	log    *logger.Logger
	cfg    SimpleCacheConfig
	toJs   func(T) (js.Value, error)
	fromJs func(js.Value) (T, error)
}

func NewSimpleCache[T any](
	log *logger.Logger,
	name string,
	cfg SimpleCacheConfig,
	toJs func(T) (js.Value, error),
	fromJs func(js.Value) (T, error),
) *SimpleCache[T] {
	return &SimpleCache[T]{
		log:    log.With(sl.Component(name)),
		cfg:    cfg,
		fromJs: fromJs,
		toJs:   toJs,
	}
}

func (c *SimpleCache[T]) Get(ctx context.Context) (T, bool) {
	const op = name + ".Get"
	if !c.cfg.Enabled {
		return *new(T), false
	}
	res, err := Await(ctx, c.cfg.Get.Invoke())
	if err != nil {
		c.log.Error(ctx, "failed to get from cache", sl.Op(op), sl.Err(err))
		return *new(T), false
	}
	if res.IsNull() {
		return *new(T), false
	}
	val, err := c.fromJs(res)
	if err != nil {
		c.log.Error(ctx, "failed to parse value from cache", sl.Op(op), sl.Err(err))
		return *new(T), false
	}
	return val, true
}

func (c *SimpleCache[T]) Add(ctx context.Context, value T) error {
	const op = name + ".Add"
	if !c.cfg.Enabled {
		return nil
	}
	jsValue, err := c.toJs(value)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	_, err = Await(ctx, c.cfg.Add.Invoke(jsValue))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
