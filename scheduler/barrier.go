package scheduler

import (
	"sync"
)

type Barrier struct {
	n       int
	workers int
	wait    *sync.Cond
}

func NewBarrier(n int) *Barrier {
	return &Barrier{n, 0, sync.NewCond(&sync.Mutex{})}
}

func (b *Barrier) Arrive() {
	b.wait.L.Lock()
	b.workers++

	if b.workers < b.n {
		b.wait.Wait()
	} else {
		b.workers = 0
		b.wait.Broadcast()
	}

	b.wait.L.Unlock()
}
