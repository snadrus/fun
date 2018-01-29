// Fun error chains
package fun

import (
	"errors"
	"runtime"
	"strings"
)

type ErrorChain interface {
	// If this condition, produce an error with the following text
	If(condition bool, thenErrText string) ErrorChain

	// Then IIFE here
	// Ex:  Then(func()error{ return errors.New("hi")}())
	Then(error) ErrorChain

	// If b, run f, else run g. Either can be nil to skip.
	// Ex:    ...IfElse(!u.InTest, u.DoNextStep, nil)
	IfElse(b bool, f func() error, g func() error) ErrorChain

	// Improve the error. %e in the string is replaced with the error
	// Ex: If(sqlQuery().Scan()).AugmentError("Sql was unhappy: %e").ReturnError()
	Explain(s string) ErrorChain

	// Gets the error with this (enclosing) function name included.
	GetError() error

	// For a simple range (0 to i-1):   fun.For(len(z), func(i idx)error{ fmt.Println(z[i]) })
	For(i int, f func(i int) error) ErrorChain

	// Run commands in parallel and wait for them to complete before proceeding
	Parallel(limit uint, f func(g GoMaker)) ErrorChain
}

type usd struct {
	err error
}

func (u *usd) GetError() error { // must GetError so nil is possible
	if u.err != nil {
		u.err = errors.New(getCaller(3) + u.err.Error())
	}
	return u.err
}

func getCaller(i int) string {
	res := make([]uintptr, i*2)
	runtime.Callers(2, res)
	f, _ := runtime.CallersFrames(res).Next()
	return f.Function
}

// If this condition, produce an error with the following text
func If(b bool, errText string) ErrorChain {
	u := &usd{}
	return u.If(b, errText)
}

// If b, run f, else run g. Either can be nil to skip.
func IfElse(b bool, f func() error, g func() error) ErrorChain {
	if b && f != nil {
		return Then(f())
	} else if g != nil {
		return Then(g())
	}
	return &usd{}
}

// Then IIFE here:  Then(func()error{ return errors.New("hi")}())
func Then(e error) ErrorChain {
	return &usd{err: e}
}

func (u *usd) If(b bool, errText string) ErrorChain {
	if u.err != nil && b {
		u.err = errors.New(errText)
	}
	return u
}

func (u *usd) IfElse(b bool, f func() error, g func() error) ErrorChain {
	if u.err != nil {
		return u
	}
	if b && f != nil {
		return u.Then(f())
	} else if g != nil {
		return u.Then(g())
	}
	return u
}

func (u *usd) Then(e error) ErrorChain {
	if u.err == nil {
		u.err = e
	}
	return u
}

func (u *usd) Explain(s string) ErrorChain {
	if u.err != nil {
		u.err = errors.New(strings.Replace(s, "%e", u.err.Error(), 1))
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
