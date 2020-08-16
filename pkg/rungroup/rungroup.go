package rungroup

import (
	"context"
	"sync"
)

// A RunGroup mirrors the interface of an error group.
type RunGroup interface {
	Go(f func(ctx context.Context) error)
	Wait() error
}

// group is akin to an error group, but will stop all routines when any routine
// exits instead of waiting for each to finish.
type group struct {
	ctx    context.Context
	cancel func()

	waitgroup sync.WaitGroup
	once      sync.Once
	err       error
}

// New creates a new RunGroup.
func New(ctx context.Context) *group {
	ctx, cancel := context.WithCancel(ctx)
	return &group{ctx: ctx, cancel: cancel}
}

// Wait waits for the wait group, and will cancel the context after the
// waitgroup stops blocking.
func (g *group) Wait() error {
	g.waitgroup.Wait()
	g.cancel()

	return g.err
}

// Go starts the given function in a new routine, adds it to the waitgroup, and
// cancels the group's context if/when it exits.
func (g *group) Go(f func(ctx context.Context) error) {
	g.waitgroup.Add(1)

	go func() {
		defer g.waitgroup.Done()

		if err := f(g.ctx); err != nil {
			g.once.Do(func() { g.err = err })
		}

		g.cancel()
	}()
}
