package di

import (
	"context"
	"fmt"
	"sync"
)

// Container defines a minimalistic dependency injection container.
// It supports both singleton and scoped lifetime dependencies.
type Container interface {
	// AddSingleton registers a dependency with a singleton lifetime.
	AddSingleton(key Key, fn DepFactoryFunc)
	// AddScoped registers a dependency with a scoped lifetime.
	AddScoped(key Key, fn DepFactoryFunc)
	// Scoped returns a new context with a child container that inherits the parent container's dependencies.
	Scoped(ctx context.Context) context.Context
	// Get retrieves a depentency (instance) by the provided key.
	Get(key Key) any
}

// DepFactoryFunc builds an instance of a dependency using the provided container.
type DepFactoryFunc func(c Container) (any, error)

type scope int

const (
	singleton scope = iota + 1
	scoped
)

var _ Container = (*container)(nil)

type container struct {
	parent     *container
	deps       map[Key]depDef
	instances  map[Key]any
	cycleGuard cycleGuard
	mu         sync.Mutex
}

var _ Container = (*container)(nil)

// depDef defines how to construct a dependency (instance) associated with a specific key.
type depDef struct {
	key         Key
	scope       scope
	factoryFunc DepFactoryFunc
}

// New creates a new, empty container.
func New() *container {
	return &container{
		deps:      make(map[Key]depDef),
		instances: make(map[Key]any),
	}
}

func (c *container) AddSingleton(key Key, fn DepFactoryFunc) {
	c.deps[key] = depDef{
		key:         key,
		scope:       singleton,
		factoryFunc: fn,
	}
}

func (c *container) AddScoped(key Key, fn DepFactoryFunc) {
	c.deps[key] = depDef{
		key:         key,
		scope:       scoped,
		factoryFunc: fn,
	}
}

func (c *container) Scoped(ctx context.Context) context.Context {
	return context.WithValue(ctx, containerKey, c.scoped())
}

func (c *container) Get(key Key) any {
	info, exists := c.deps[key]
	if !exists {
		panic(fmt.Sprintf("there is no dependency registered with `%s`", key))
	}

	// catch cases of: building Foo needs Bar and building Bar needs Foo :boom:
	if _, exists := c.cycleGuard[info.key]; exists {
		panic(fmt.Sprintf("cyclic dependencies encountered while building `%s`, tracked: %s", info.key, c.cycleGuard))
	}

	if info.scope == singleton {
		return c.getFromParent(info)
	}

	return c.getInstance(info)
}

func (c *container) getFromParent(info depDef) any {
	if c.parent != nil {
		return c.parent.getFromParent(info)
	}

	return c.getInstance(info)
}

type builtDone = chan struct{}

func (c *container) getInstance(dep depDef) any {
	c.mu.Lock()

	// check if the value is already built
	ie, exists := c.instances[dep.key]
	if !exists {
		done := make(builtDone)
		c.instances[dep.key] = done // assign an in-progress channel to indicate that the build is in progress
		c.mu.Unlock()
		return c.buildInstance(dep, done)
	}

	c.mu.Unlock()
	done, isDone := ie.(builtDone)
	// if the value is not a done channel, it means it was already built
	if !isDone {
		return ie
	}

	// wait for the build to finish
	<-done

	// after the build is done, we can safely get the value again
	return c.getInstance(dep)
}

func (c *container) buildInstance(dep depDef, done builtDone) any {
	ie, err := dep.factoryFunc(c.builder(dep))

	c.mu.Lock()

	// check if the build was successful
	if err != nil {
		// if the build failed, we need to clean up and panic
		delete(c.instances, dep.key)
		c.mu.Unlock()
		close(done)
		panic(fmt.Sprintf("error building dependency `%s`: %s", dep.key, err))
	}

	// if the build was successful, we can store the value
	c.instances[dep.key] = ie
	c.mu.Unlock()
	close(done)

	return ie
}

// scoped returns a child container that inherits the parent container's dependencies
func (c *container) scoped() *container {
	return &container{
		parent:    c,
		deps:      c.deps,
		instances: make(map[Key]any),
	}
}

func (c *container) builder(dep depDef) *container {
	return &container{
		parent:     c.parent,
		deps:       c.deps,
		instances:  c.instances,
		cycleGuard: c.cycleGuard.add(dep),
	}
}
