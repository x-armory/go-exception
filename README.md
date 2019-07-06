# go-exception
go-exception is stacktraceable error, like java exception

Installation
------------

Use go get.

	go get github.com/x-armory/go-exception

or to update

	go get -u github.com/x-armory/go-exception

Then import the validator package into your own code.

	import "github.com/x-armory/go-exception"


##### Examples:

```go
package main

import (
	"fmt"
	"github.com/x-armory/go-exception"
	"io/ioutil"
)

func main() {
	ex.SetMiddleErrorStackMaxDepth(5) // default is 5
	ex.SetRootErrorStackMaxDepth(15)  // default is 32

	// in global exception handler
	ex.Try(controllerMethod).SafeCatch(func(err interface{}) {
		// err3, type is ex.Exception, cause is err2, root cause is err1
		e := ex.Wrap(err)
		e.PrintErrorStack() // print callers

		rootError := e.RootCause().WrappedCause()
		fmt.Printf("root error type is: %T\nroot error is: %+v\n\n", rootError, rootError)

		println("\nexception print again")
		e.Throw() // will be catch and print again
	})

}

func controllerMethod() {
	// bind & validate request
	// ...

	// do biz mothod
	ex.Try(bizMethod).Catch(func(err interface{}) {
		// err2, type is ex.Exception, cause is err1
		e2 := ex.Wrap(err)

		// err3, type is ex.Exception, cause is err2, root cause is err1
		ex.Exception("BizErrorCode_001", "do biz failed", e2).Throw() // in catch func, will throw out
	})
}

func bizMethod() {
	defer func() {
		println("release resources for biz") //will print finally
	}()

	println("start do biz")

	// will throw err1, type is ex.Exception, wrappedCause is *os.PathError
	ex.AssertNoError(thirdPartyMethod(), "do thirdPartyMethod failed")

	// do next steps if thirdPartyMethod not return error
	// ...
	println("finish do biz") // will not print
}

func thirdPartyMethod() error {
	_, e := ioutil.ReadFile("/unknown/file")
	return e
}

```

###### output:

```go
start do biz
release resources for biz

root error type is: *os.PathError
root error is: open /unknown/file: no such file or directory

Exception: <BizErrorCode_001> do biz failed
   at: /Volumes/C/go/src/github.com/x-armory/go-exception/sample/sample_try_catch_exception_assert_app.go:38 main.controllerMethod.func1()
   at: /Volumes/C/go/src/github.com/x-armory/go-exception/try.go:20 Catch()
   at: /Volumes/C/go/src/github.com/x-armory/go-exception/sample/sample_try_catch_exception_assert_app.go:33 main.controllerMethod()
   at: /Volumes/C/go/src/github.com/x-armory/go-exception/try.go:10 exception.Try()
   at: /Volumes/C/go/src/github.com/x-armory/go-exception/sample/sample_try_catch_exception_assert_app.go:14 main.main()
caused by Exception: <ASSERTION_FAILED> do thirdPartyMethod failed
   at: /Volumes/C/go/src/github.com/x-armory/go-exception/sample/sample_try_catch_exception_assert_app.go:50 main.bizMethod()
   at: /Volumes/C/go/src/github.com/x-armory/go-exception/try.go:10 exception.Try()
   at: /Volumes/C/go/src/github.com/x-armory/go-exception/sample/sample_try_catch_exception_assert_app.go:33 main.controllerMethod()
   at: /Volumes/C/go/src/github.com/x-armory/go-exception/try.go:10 exception.Try()
   at: /Volumes/C/go/src/github.com/x-armory/go-exception/sample/sample_try_catch_exception_assert_app.go:14 main.main()
   at: /usr/local/Cellar/go/1.12.4/libexec/src/runtime/proc.go:200 runtime.main()
caused by open /unknown/file: no such file or directory

exception print again
Exception: <BizErrorCode_001> do biz failed
   at: /Volumes/C/go/src/github.com/x-armory/go-exception/sample/sample_try_catch_exception_assert_app.go:38 main.controllerMethod.func1()
   at: /Volumes/C/go/src/github.com/x-armory/go-exception/try.go:20 Catch()
   at: /Volumes/C/go/src/github.com/x-armory/go-exception/sample/sample_try_catch_exception_assert_app.go:33 main.controllerMethod()
   at: /Volumes/C/go/src/github.com/x-armory/go-exception/try.go:10 exception.Try()
   at: /Volumes/C/go/src/github.com/x-armory/go-exception/sample/sample_try_catch_exception_assert_app.go:14 main.main()
caused by Exception: <ASSERTION_FAILED> do thirdPartyMethod failed
   at: /Volumes/C/go/src/github.com/x-armory/go-exception/sample/sample_try_catch_exception_assert_app.go:50 main.bizMethod()
   at: /Volumes/C/go/src/github.com/x-armory/go-exception/try.go:10 exception.Try()
   at: /Volumes/C/go/src/github.com/x-armory/go-exception/sample/sample_try_catch_exception_assert_app.go:33 main.controllerMethod()
   at: /Volumes/C/go/src/github.com/x-armory/go-exception/try.go:10 exception.Try()
   at: /Volumes/C/go/src/github.com/x-armory/go-exception/sample/sample_try_catch_exception_assert_app.go:14 main.main()
   at: /usr/local/Cellar/go/1.12.4/libexec/src/runtime/proc.go:200 runtime.main()
caused by open /unknown/file: no such file or directory
```