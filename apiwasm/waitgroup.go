package apiwasm

import (
	"context"
	"log"
	"sync"

	"github.com/iansmith/parigot/apishared/id"
)

type MustRequireFunc func(context.Context, id.ServiceId)

var waitGroupVerbose = true

type ParigotWaitGroup struct {
	*sync.WaitGroup
	name  string
	count int
	lock  *sync.Mutex
}

func NewParigotWaitGroup(name string) *ParigotWaitGroup {
	return &ParigotWaitGroup{
		WaitGroup: &sync.WaitGroup{},
		name:      name,
		lock:      &sync.Mutex{},
	}
}

func (w *ParigotWaitGroup) Add(delta int) {
	w.lock.Lock()
	defer w.lock.Unlock()

	w.count += delta
	if waitGroupVerbose {
		log.Printf("%s, add called: %d is new value ", w.name, w.count)
	}
	w.WaitGroup.Add(delta)
}

func (w *ParigotWaitGroup) Done() {
	w.lock.Lock()
	defer w.lock.Unlock()
	w.count -= 1
	if waitGroupVerbose {
		log.Printf("%s, done called: %d is new value ", w.name, w.count)
	}
	// uncomment if you want to see WHERE the Done was called from
	//debug.PrintStack()
	w.WaitGroup.Done()
}

func (w *ParigotWaitGroup) Wait() {
	if waitGroupVerbose {
		log.Printf("%s, wait started with value: %d ", w.name, w.count)
	}
	w.WaitGroup.Wait()
	if waitGroupVerbose {
		log.Printf("%s, wait completed: %d ", w.name, w.count)
	}
}

func (w *ParigotWaitGroup) Value() int {
	w.lock.Lock()
	defer w.lock.Unlock()
	return w.count
}
