package mapsync

import (
	"sync"

	"github.com/cornelk/hashmap"
)

type Mutex struct {
	lockmap hashmap.HashMap
}
type mapped struct {
	ctx *Mutex
	v   interface{}
}

func (mutex *Mutex) Map(v interface{}) *mapped {
	return &mapped{mutex, v}
}

func (hashed *mapped) Lock() {
	hashed.ctx.Lock(hashed.v)
}

func (hashed *mapped) Unlock() {
	hashed.ctx.Unlock(hashed.v)
}

func (mutex *Mutex) Lock(v interface{}) {
	locker, _ := mutex.lockmap.GetOrInsert(v, &sync.Mutex{})
	locker.(interface{ Lock() }).Lock()
}

func (mutex *Mutex) Unlock(v interface{}) {
	locker, found := mutex.lockmap.Get(v)
	if !found {
		panic("mutex: unlock of unlocked mutex")
	}
	locker.(interface{ Unlock() }).Unlock()
}
