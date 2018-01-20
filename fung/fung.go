// Fun GoRoutine chained API to return errors only
package fung

import "github.com/snadrus/fun"

type GoChain interface {
	fun.ErrorChain

	Wait() fun.ErrorChain

	WaitLater(wait *func() error) GoChain

	// Go routine run, catching panic. Waits if max is set.
	// Errors do not block until wait.
	Go(func() error) GoChain

	// Name & run a function, or run it after stated dependencies.
	// Ex: "a"
	// Ex: "a,b>c"  depends on a and be before running c. Errors in those terminate early.
	Name(s string, f func() error) GoChain

	// PoolMax concurrency. less than 1 == infinite. Do early.
	PoolMax(int) GoChain
}

func PoolMax(i int) GoChain {
	u := &usdGo{ErrorChain: fun.Then(nil)}
	return u.PoolMax(i)
}

// Goroutine run, catching panic. Waits if max is set.
func Go(f func() error) GoChain {
	u := &usdGo{ErrorChain: fun.Then(nil)}
	return u.Go(f)
}

func Name(s string, f func() error) GoChain {
	u := &usdGo{ErrorChain: fun.Then(nil)}
	return u.Name(s, f)
}

type usdGo struct {
	fun.ErrorChain
	ct  int
	max int
}

func (u *usdGo) PoolMax(i int) GoChain {
	if i < 0 {
		i = 0
	}
	u.max = i
	return u
}

func (u *usdGo) WaitLater(wait *func() error) GoChain {
	*wait = u.Wait
	return u
}

func (u *usdGo) Go(func() error) GoChain {
	// TODO
	// TODO FUTURE for loop. This is best composed and not in either. funr
	return u
}

func (u *usdGo) Wait() fun.ErrorChain {
	// TODO
	return u
}

// .BackToFung() -- traverse the tree? Can we have trees here?
// funr.Range(ErrorChain, length, func(idx int)error)
// funr.Range(length, func(idx int)error)
/*
fung.PoolMax(20).Then(funr.Range(len(sqlResults), func(idx int)error{
	fung.FromStack().Go(func()error{
		fmt.Println(sqlResults[idx])
	})
})).Wait()

Problems:
  - FromStack hard!
  - .Wait missing on a .Then (ErrorChain). Is parent lost?
*/

func (u *usdGo) Name(s string, f func() error) GoChain {
	// TODO
	return u
}
