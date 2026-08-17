package sugar

import (
	"sync"
	"sync/atomic"
)

type Singleton[T any] interface {
	Load() T
	Inject(newFn func() T)
	Replace(newObj T)
}

type singleton[T any] struct {
	newFunc func() T
	once    sync.Once
	obj     atomic.Pointer[T]
}

func NewSingleton[T any](newFn func() T) Singleton[T] {
	return &singleton[T]{newFunc: newFn}
}

func (s *singleton[T]) Load() T {
	s.once.Do(func() {
		vv := s.newFunc()
		s.obj.Store(&vv)
	})

	obj := s.obj.Load()
	if obj == nil {
		panic("singleton object is nil")
	}

	return *obj
}

func (s *singleton[T]) Reload() {
	vv := s.newFunc()
	s.obj.Store(&vv)
}

func (s *singleton[T]) Inject(newFn func() T) {
	vv := newFn()
	s.obj.Store(&vv)
}

func (s *singleton[T]) Replace(newObj T) {
	s.obj.Store(&newObj)
}

type singletonWithValue[T any, I any] struct {
	newFunc func(I) T
	once    sync.Once
	obj     atomic.Pointer[T]
}

type SingletonWithValue[T any, I any] interface {
	Singleton[T]

	Init(value I) T
	Reload(value I)
}

func NewSingletonWithValue[T any, I any](newFn func(I) T) SingletonWithValue[T, I] {
	return &singletonWithValue[T, I]{newFunc: newFn}
}

func (s *singletonWithValue[T, I]) Init(value I) T {
	s.once.Do(func() {
		vv := s.newFunc(value)
		s.obj.Store(&vv)
	})

	obj := s.obj.Load()
	if obj == nil {
		panic("singleton object is nil")
	}

	return *obj
}

func (s *singletonWithValue[T, I]) Load() T {
	obj := s.obj.Load()
	if obj == nil {
		panic("singleton object is nil")
	}

	return *obj
}

func (s *singletonWithValue[T, I]) Reload(value I) {
	vv := s.newFunc(value)
	s.obj.Store(&vv)
}

func (s *singletonWithValue[T, I]) Inject(newFn func() T) {
	vv := newFn()
	s.obj.Store(&vv)
}

func (s *singletonWithValue[T, I]) Replace(newObj T) {
	s.obj.Store(&newObj)
}
