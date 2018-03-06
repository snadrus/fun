// Fun error chains
package fun

import (
	"errors"
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"
)

type ErrorChain interface {
	// If this condition, produce an error with the following text
	If(condition bool, thenErrText string) ErrorChain

	// Recover allows enclosed function's PANICs to be captured as errors in the chain.
	// Works great with fun.PanicOnError
	Recover(f func()) ErrorChain

	// Improve the error. %e in the string is replaced with the error
	// Ex: If(sqlQuery().Scan()).AugmentError("Sql was unhappy: %e").ReturnError()
	Explain(s string) ErrorChain

	// Gets the error with this (enclosing) function name included.
	GetError() error

	// For a simple range (0 to i-1):   fun.For(len(z), func(i idx)error{ fmt.Println(z[i]) })
	For(i int, f func(i int) error) ErrorChain

	// Run commands in parallel and wait for them to complete before proceeding.
	// limit is how many concurrent Go calls (not GoNamed) are allowed. 0==unlimited
	Parallel(limit uint, f func(g GoMaker)) ErrorChain
}

type usd struct {
	err       error
	explained bool
}

func (u *usd) GetError() error { // must GetError so nil is possible
	if u.err != nil {
		u.err = errors.New(getCaller(3) + " " + u.err.Error())
	}
	return u.err
}

func getCaller(i int) string {
	res := make([]uintptr, i*2)
	runtime.Callers(i, res)
	f, _ := runtime.CallersFrames(res).Next()
	return f.Function
}

// If this condition, produce an error with the following text
func If(b bool, errText string) ErrorChain {
	u := &usd{}
	return u.If(b, errText)
}

func (u *usd) If(b bool, errText string) ErrorChain {
	if u.err != nil && b {
		u.err = errors.New(errText)
	}
	return u
}

func Recover(f func()) ErrorChain {
	return (&usd{}).Recover(f)
}
func (u *usd) Recover(f func()) (u2 ErrorChain) {
	if u.err != nil {
		return u
	}
	defer func() {
		if v := recover(); v != nil {
			u.err = fmt.Errorf("%v \n Stack: %s", v, string(debug.Stack()))
		}
		u2 = u
	}()
	f()
	return u
}

// Helper for Recover scenarios
func Err2Panic(err error) {
	if err != nil {
		panic(err)
	}
}

func (u *usd) Explain(s string) ErrorChain {
	if u.err != nil && !u.explained {
		u.err = errors.New(strings.Replace(s, "%e", u.err.Error(), 1))
		u.explained = true
	}
	return u
}

func For(i int, f func(i int) error) ErrorChain {
	return (&usd{}).For(i, f)
}
func (u *usd) For(i int, f func(i int) error) ErrorChain {
	if u.err != nil {
		return u
	}
	for a := 0; a < i; a++ {
		if u.err = f(a); u.err != nil {
			return u
		}
	}
	return u
}
