package options

import "sync"

type RWLock interface {
	RLock()
	RUnlock()
	sync.Locker
}

type options struct {
	Mux RWLock
}

func New() options {
	return options{Mux: &fakeLock{}}
}

type Option func(opt *options)

func WithLocker() Option {
	return func(opt *options) {
		opt.Mux = &sync.RWMutex{}
	}
}

type fakeLock struct {
}

func (f fakeLock) RLock() {

}

func (f fakeLock) RUnlock() {
}

func (f fakeLock) Lock() {
}

func (f fakeLock) Unlock() {
}
