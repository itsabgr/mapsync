package mapsync

import (
	"sync"

	"github.com/cornelk/hashmap"
)

type WaitGroup struct {
	lockmap hashmap.HashMap
}
type wgmapped struct {
	ctx *WaitGroup
	v   interface{}
}

func (mutex *WaitGroup) Map(v interface{}) *wgmapped {
	return &wgmapped{mutex, v}
}

func (hashed *wgmapped) Add(delta int) {
	hashed.ctx.Add(hashed.v, delta)
}
func (hashed *wgmapped) Done() {
	hashed.ctx.Done(hashed.v)
}
func (hashed *wgmapped) Wait() {
	hashed.ctx.Wait(hashed.v)
}

func (mutex *WaitGroup) Add(v interface{}, delta int) {
	wg, _ := mutex.lockmap.GetOrInsert(v, &sync.WaitGroup{})
	wg.(interface{ Add(int) }).Add(delta)
}

func (mutex *WaitGroup) Wait(v interface{}) {
	wg, _ := mutex.lockmap.Get(v)
	wg.(interface{ Wait() }).Wait()
}

func (mutex *WaitGroup) Done(v interface{}) {
	wg, _ := mutex.lockmap.Get(v)
	wg.(interface{ Done() }).Done()
}
