package fun

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
)

// We already have a For loop function,

type GoMaker interface {
	// Go (as part of a wait group)
	Go(f func() error)

	// GoNamed exmaple using deps:
	// GoNamed("first", blaFirst, firstCleanup)    // starts immediately
	// 			firstCleanup  always runs once everything else is done (if first succeeds)
	// GoNamed("second", bla,nil)   // starts immediately
	// GoNamed("first,second>third", bla,nil)  // starts after others succeed
	GoNamed(csvDepsGtName string, f func() error, deferFunc func())
}

type goMaker struct {
	signal        sync.WaitGroup // only used for signaling
	atomicCount   int64
	atomicWaiting uint64
	err           error
	mutex         sync.RWMutex
	todos         map[string]chan bool
	defers        []func()
}

func Parallel(limit uint, f func(g GoMaker)) ErrorChain {
	return (&usd{}).Parallel(limit, f)
}
func (u *usd) Parallel(limit uint, f func(g GoMaker)) ErrorChain {
	if u.err != nil {
		return u
	}
	signal := sync.WaitGroup{}
	signal.Add(1)
	g := &goMaker{signal: signal, todos: map[string]chan bool{}, defers: []func(){}}
	f(g)

	atomic.StoreUint64(&g.atomicWaiting, 1)
	signal.Wait()
	for _, d := range g.defers {
		d()
	}
	u.err = g.err
	return u
}

func (g *goMaker) setErr(s string) {
	g.mutex.Lock()
	g.mutex.Unlock()
	g.err = fmt.Errorf("%s", s)
}

func (g *goMaker) wgDone() {
	tmp := atomic.AddInt64(&g.atomicCount, -1)
	if tmp == 0 { // either we are done, or something else is yet-to-be-added
		if atomic.LoadUint64(&g.atomicWaiting) == 1 {
			g.signal.Done()
		}
	}
}
func (g *goMaker) Go(f func() error) {
	atomic.AddInt64(&g.atomicCount, 1)
	go func() {
		defer g.wgDone()

		g.mutex.RLock()
		tmp := g.err
		g.mutex.Lock()
		if tmp != nil {
			return
		}

		defer func() {
			if v := recover(); v != nil {
				g.setErr(fmt.Sprint(v))
			}
		}()

		err := f()
		if err != nil {
			g.setErr(err.Error())
		}
	}()
}

func (g *goMaker) GoNamed(csvDepsGtName string, f func() error, fDefer func()) {
	res := strings.SplitN(csvDepsGtName, ">", 2)
	name := res[len(res)-1]
	waitFor := []chan bool{}
	g.mutex.Lock()
	defer g.mutex.Unlock()
	if nil != g.todos[name] {
		panic("bad code: cannot have 2 gonames the same")
	}
	g.todos[name] = make(chan bool)
	depAry := res[:len(res)-1]
	if len(depAry) > 0 {
		for _, dep := range strings.Split(depAry[0], ",") {
			ch, ok := g.todos[dep]
			if !ok {
				ch = make(chan bool)
				g.todos[dep] = ch
			}
			waitFor = append(waitFor, ch)
		}
	}

	atomic.AddInt64(&g.atomicCount, 1)
	go func() {
		defer g.wgDone()
		for _, depCh := range waitFor {
			<-depCh
		}
		g.mutex.RLock()
		err := g.err
		g.mutex.RUnlock()
		if err != nil {
			return
		}
		if f != nil {
			err = f()
			if err != nil {
				g.mutex.Lock()
				if g.err != nil {
					g.err = err
				}
				g.mutex.Unlock()
			}
		}
		if err != nil && fDefer != nil {
			g.mutex.Lock()
			g.defers = append([]func(){fDefer}, g.defers...)
			g.mutex.Unlock()
		}
	}()
}
