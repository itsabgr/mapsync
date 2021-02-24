package mapsync

import (
	"sync"

	"github.com/cornelk/hashmap"
)

type RWMutex struct {
	lockmap hashmap.HashMap
}
type rwmapped struct {
	ctx *RWMutex
	v   interface{}
}

func (mutex *RWMutex) Map(v interface{}) *rwmapped {
	return &rwmapped{mutex, v}
}

func (hashed *rwmapped) Lock() {
	hashed.ctx.Lock(hashed.v)
}

func (hashed *rwmapped) Unlock() {
	hashed.ctx.Unlock(hashed.v)
}

func (hashed *rwmapped) RUnlock() {
	hashed.ctx.RUnlock(hashed.v)
}

func (hashed *rwmapped) RLock() {
	hashed.ctx.RLock(hashed.v)
}

func (hashed *rwmapped) RLocker() {
	hashed.ctx.RLocker(hashed.v)
}

func (mutex *RWMutex) Lock(v interface{}) {
	locker, _ := mutex.lockmap.GetOrInsert(v, &sync.RWMutex{})
	locker.(interface{ Lock() }).Lock()
}

func (mutex *RWMutex) RLock(v interface{}) {
	locker, _ := mutex.lockmap.GetOrInsert(v, &sync.RWMutex{})
	locker.(interface{ RLock() }).RLock()
}
func (mutex *RWMutex) RLocker(v interface{}) sync.Locker {
	locker, _ := mutex.lockmap.GetOrInsert(v, &sync.RWMutex{})
	return locker.(interface{ RLocker() sync.Locker }).RLocker()
}
func (mutex *RWMutex) RUnlock(v interface{}) {
	locker, found := mutex.lockmap.Get(v)
	if !found {
		panic("mapsync: RUnlock of unlocked RWMutex")
	}
	locker.(interface{ RUnlock() }).RUnlock()
}

func (mutex *RWMutex) Unlock(v interface{}) {
	locker, found := mutex.lockmap.Get(v)
	if !found {
		panic("lockmap: unlock of unlocked mutex")
	}
	locker.(interface{ Unlock() }).Unlock()
}
