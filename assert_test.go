package ex

import (
	"gopkg.in/go-playground/assert.v1"
	"testing"
)

func TestAssertErrorStack1(t *testing.T) {
	Try(func() {
		func() {
			func() {
				func() {
					Try(func() {
						Assert(false, "error message")
					}).Catch(func(e interface{}) {
						Wrap(e).Throw()
					})
				}()
			}()
		}()
	}).Catch(func(e interface{}) {
		assert.Equal(t, Wrap(e).Code(), AssertionErrorCode)
		assert.Equal(t, Wrap(e).Message(), "error message")
	})
}

func TestAssertErrorStackWithoutMessage(t *testing.T) {
	Try(func() {
		func() {
			func() {
				func() {
					Try(func() {
						Assert(false)
					}).Catch(func(e interface{}) {
						Wrap(e).Throw()
					})
				}()
			}()
		}()
	}).Catch(func(e interface{}) {
		assert.Equal(t, Wrap(e).Code(), AssertionErrorCode)
		assert.Equal(t, Wrap(e).Message(), "")
	})
}
