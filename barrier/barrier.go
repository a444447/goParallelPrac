package main

import "sync"

type Barrier struct {
	count int
	cond  *sync.Cond
	mutex *sync.Mutex
	total int
}

func NewBarrier(total int) *Barrier {
	lockToUse := &sync.Mutex{}
	condToUse := sync.NewCond(lockToUse)
	return &Barrier{total, condToUse, lockToUse, total}
}

func (b *Barrier) Wait() {
	b.mutex.Lock()
	b.count--
	if b.count == 0 {
		b.count = b.total
		b.cond.Broadcast()
	} else {
		b.cond.Wait()
	}
	b.mutex.Unlock()
}
