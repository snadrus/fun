package fun

import (
	"errors"
	"runtime"
	"strings"
)

type ErrBridge interface {
	If(condition bool, thenErrText string) ErrBridge
	Then(error) ErrBridge
	AugmentError(s string) ErrBridge
	GetError() error
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

func If(b bool, errText string) ErrBridge {
	u := &usd{}
	return u.If(b, errText)
}

func (u *usd) If(b bool, errText string) ErrBridge {
	if u.err != nil && b {
		u.err = errors.New(errText)
	}
	return u
}

func Then(e error) ErrBridge {
	return &usd{err: e}
}

func (u *usd) Then(e error) ErrBridge {
	if u.err != nil {
		u.err = e
	}
	return u
}

func (u *usd) AugmentError(s string) ErrBridge {
	if u.err != nil {
		u.err = errors.New(strings.Replace(s, "%e", u.err.Error(), 1))
	}
	return u
}
