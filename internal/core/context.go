package core

import (
	"fmt"
	"reflect"
	"sync"
)

type ApplicationContext struct {
	mu    sync.RWMutex
	beans map[reflect.Type]any
}

var (
	instance *ApplicationContext
	once     sync.Once
)

// Singleton
func GetContext() *ApplicationContext {
	once.Do(func() {
		instance = &ApplicationContext{
			beans: make(map[reflect.Type]any),
		}
	})
	return instance
}

func (c *ApplicationContext) register(t reflect.Type, bean any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.beans[t] = bean
}

func (c *ApplicationContext) get(t reflect.Type) (any, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	bean, ok := c.beans[t]
	if !ok {
		return nil, fmt.Errorf("bean of type %s not found", t.String())
	}
	return bean, nil
}

func Register[T any](bean T) {
	t := reflect.TypeOf((*T)(nil)).Elem()
	GetContext().register(t, bean)
}

func Get[T any]() (T, error) {
	var zero T
	t := reflect.TypeOf((*T)(nil)).Elem()

	bean, err := GetContext().get(t)
	if err != nil {
		return zero, err
	}
	return bean.(T), nil
}

func MustGet[T any]() T {
	b, err := Get[T]()
	if err != nil {
		panic(err)
	}
	return b
}
