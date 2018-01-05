package fun

import (
	"errors"
	"runtime"
	"strings"
)

type ErrBridge interface {
	So(bool) ErrBridge
	Then(func() error) ErrBridge
	ElseErr(s string) ErrBridge
	GetError() error
}

type usd struct {
	err error
}

func (u *usd) GetError() error {
	return u.err
}

// so you can just return one of these
func (u *usd) Error() error {
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

func So(b bool) ErrBridge {
	u := &usd{}
	return u.So(b)
}

func (u *usd) So(b bool) ErrBridge {
	if u.err != nil && !b {
		u.err = errors.New("Validation failed")
	}
	return u
}

func Then(f func() error) ErrBridge {
	return &usd{err: f()}
}

func (u *usd) Then(f func() error) ErrBridge {
	if u.err != nil {
		u.err = f()
	}
	return u
}

func (u *usd) ElseErr(s string) ErrBridge {
	if u.err != nil {
		u.err = errors.New(strings.Replace(s, "%e", u.err.Error(), 1))
	}
	return u
}
