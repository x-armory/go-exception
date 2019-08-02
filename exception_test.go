package ex

import (
	"errors"
	"fmt"
	"gopkg.in/go-playground/assert.v1"
	"io"
	"testing"
	"time"
)

func TestSub(t *testing.T) {
	println(fmt.Sprintf("%v", 12345)[:3] + "...")
}
func TestErrorStack(t *testing.T) {
	e1 := Exception("code1", "message1", nil)
	e2 := Exception("code2", "message2", e1)
	e3 := Exception("code3", "message3", e2)
	e4 := Exception("code4", "message4", e3)
	println(e4.StackTraceString())

	assert.Equal(t, e4.RootCause(), e1)
	assert.Equal(t, e4.Cause(), e3)
}
func TestWrappedErrorStack(t *testing.T) {
	e1 := Wrap(io.EOF)
	e2 := Exception("code2", "message2", e1)
	e3 := Exception("code3", "message3", e2)
	e4 := Exception("code4", "message4", e3)
	println(e4.StackTraceString())

	assert.Equal(t, e4.RootCause(), e1)
	assert.Equal(t, e4.Cause(), e3)
}

// 10万次嵌套异常
// 性能主要消耗在获取调用堆栈上
//
// PanicWrapError			4层  12,998,830,000 nano 13s
// PanicWrapError			3层   8,376,735,000
// PanicWrapError			2层   4,597,791,000
// PanicWrapError			1层   1,472,826,000
//
// ReturnWrapError			4层   2,568,737,000
// ReturnWrapError			3层   2,047,857,000
// ReturnWrapError			2层   1,547,794,000
// ReturnWrapError			1层   1,170,963,000
//
// PanicBuildinError		4层     143,883,000
// PanicBuildinError		3层     103,607,000
// PanicBuildinError		2层      59,551,000
// PanicBuildinError		1层      14,780,000
//
// ReturnBuildinError		4层         496,000
// ReturnBuildinError		3层         356,000
// ReturnBuildinError		2层         191,000
// ReturnBuildinError		1层          31,000

func TestPanicWrapErrorPerformance(t *testing.T) {
	t1 := time.Now()
	for i := 0; i < 100000; i++ {
		Try(func() {
			Try(func() {
				Try(func() {
					Try(func() {
						Wrap(io.EOF).Throw()
					}).Catch(func(e interface{}) {
						Exception("code1", "message1", Wrap(e)).Throw()
					})
				}).Catch(func(e interface{}) {
					Exception("code2", "message2", Wrap(e)).Throw()
				})
			}).Catch(func(e interface{}) {
				Exception("code2", "message3", Wrap(e)).Throw()
			})
		}).Catch(func(e interface{}) {
			if Wrap(e).StackTraceString() == "" {
			}
		})
	}
	t2 := time.Now()
	println("cost", t2.UnixNano()-t1.UnixNano())
}
func TestPanicBuildinErrorPerformance(t *testing.T) {
	t1 := time.Now()
	for i := 0; i < 100000; i++ {
		Try(func() {
			Try(func() {
				Try(func() {
					Try(func() {
						panic(io.EOF)
					}).Catch(func(e interface{}) {
						panic(errors.New(fmt.Sprintf("message2 %v", e)))
					})
				}).Catch(func(e interface{}) {
					panic(errors.New(fmt.Sprintf("message3 %v", e)))
				})
			}).Catch(func(e interface{}) {
				panic(errors.New(fmt.Sprintf("message4 %v", e)))
			})
		}).Catch(func(e interface{}) {
			if e == nil {
			}
			//fmt.Printf("%v", e)
		})
	}
	t2 := time.Now()
	println("cost", t2.UnixNano()-t1.UnixNano())
}

func TestReturnWrapErrorPerformance(t *testing.T) {
	t1 := time.Now()
	for i := 0; i < 100000; i++ {
		e := func() *ExceptionClass {
			//return Wrap(func() *ExceptionClass {
			//return Wrap(func() *ExceptionClass {
			//return Wrap(func() *ExceptionClass {
			return Wrap(io.EOF)
			//}())
			//}())
			//}())
		}()
		if e.StackTraceString() == "" {
		}
	}
	t2 := time.Now()
	println("cost", t2.UnixNano()-t1.UnixNano())
}

func TestReturnBuildinErrorPerformance(t *testing.T) {
	t1 := time.Now()
	for i := 0; i < 100000; i++ {
		e := func() error {
			return func() error {
				return func() error {
					return func() error {
						return io.EOF
					}()
				}()
			}()
		}()
		if e == nil {
		}
	}
	t2 := time.Now()
	println("cost", t2.UnixNano()-t1.UnixNano())
}
