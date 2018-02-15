package fun

import "sync"

type wgStackEl struct {
	wg       *sync.WaitGroup
	waitedOn bool
}
type WgStack struct {
	s     []wgStackEl
	mutex sync.Mutex
}

func (w *WgStack) Add(i int) (doner func()) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	if len(w.s) == 0 || w.s[0].waitedOn == false {
		w.s = append([]wgStackEl{{wg: &sync.WaitGroup{}}}, w.s...)
	}
	w.s[0].wg.Add(i)
	return w.s[0].wg.Done
}

func (w *WgStack) Wait() {
	for {
		w.mutex.Lock()
		l := len(w.s)
		w.mutex.Unlock()
		if l == 0 {
			break
		}
		w.mutex.Lock()
		el := w.s[len(w.s)-1]
		el.waitedOn = true
		wg := el.wg
		w.s = w.s[:len(w.s)-1]
		w.mutex.Unlock()
		wg.Wait()
	}
}
