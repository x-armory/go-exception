package ex

import (
	"fmt"
	"os"
)

func Assert(condition bool, message ...interface{}) {
	if !condition {
		panic(buildException(AssertionErrorCode, nil, message...))
	}
}
func AssertNoError(err interface{}, message ...interface{}) {
	if err != nil {
		panic(buildException(AssertionErrorCode, err, message...))
	}
}
func ExitIf(condition bool, message ...interface{}) {
	if condition {
		buildException(ExitErrorCode, nil, message...).PrintErrorStack()
		os.Exit(1)
	}
}
func ExitIfError(err interface{}, message ...interface{}) {
	if err != nil {
		buildException(ExitErrorCode, err, message...).PrintErrorStack()
		os.Exit(1)
	}
}

func buildException(defaultCode string, cause interface{}, v ...interface{}) *exceptionClass {
	if len(v) > 0 {
		if e, ok := v[0].(*exceptionClass); ok {
			if cause != nil {
				e.cause = Wrap(cause)
			}
			return e
		}
	}
	return buildErrorWithCallerStack(stackStartDepth+1, defaultCode, buildMsg(v...), cause)
}

func buildMsg(v ...interface{}) string {
	if len(v) == 0 {
		return ""
	}
	return fmt.Sprintf(fmt.Sprintf("%v", v[0]), v[1:]...)
}
