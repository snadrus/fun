package fun

import (
	"errors"
	"runtime"
)

type ErrBridge interface {
	So(bool) ErrBridge
	Then(func() error) ErrBridge
	Error() error
}

type usd struct {
	err error
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
