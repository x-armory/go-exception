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
