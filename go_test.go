package fun

import (
	"sync/atomic"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_GoMaker(t *testing.T) {
	t.Parallel()
	Convey("GoMaker", t, func() {
		f := 0
		s := 0
		err := Parallel(0, func(g GoMaker) {
			g.GoNamed("first", func() error { f = 4; return nil }, nil)
			g.GoNamed("second", func() error { s = 6; return nil }, nil)
			g.GoNamed("first,second>last", func() error { f *= 2; s *= 3; return nil }, nil)
		}).GetError()
		So(err, ShouldBeNil)
		So(f, ShouldEqual, 8)
		So(s, ShouldEqual, 18)
	})
	Convey("MaxParallel", t, func() {
		var max int32
		var ct int32
		err := Parallel(5, func(g GoMaker) {
			for a := 0; a < 100; a++ {
				g.Go(func() error {
					defer atomic.AddInt32(&max, -1)
					atomic.AddInt32(&ct, 1)
					return If(atomic.AddInt32(&max, 1) > 5, "max too big").GetError()
				})
			}
		}).GetError()
		So(err, ShouldBeNil)
		So(ct, ShouldEqual, 100)
	})
}
