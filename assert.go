package ex

import (
	"fmt"
	"os"
)

func Assert(condition bool, message ...interface{}) {
	if !condition {
		panic(buildErrorWithCallerStack(stackStartDepth, AssertionErrorCode, buildMsg(message...), nil))
	}
}
func AssertNoError(err interface{}, message ...interface{}) {
	if err != nil {
		panic(buildErrorWithCallerStack(stackStartDepth, AssertionErrorCode, buildMsg(message...), err))
	}
}
func ExitIf(condition bool, message ...interface{}) {
	if condition {
		buildErrorWithCallerStack(stackStartDepth, ExitErrorCode, buildMsg(message...), nil).PrintErrorStack()
		os.Exit(1)
	}
}
func ExitIfError(err interface{}, message ...interface{}) {
	if err != nil {
		buildErrorWithCallerStack(stackStartDepth, ExitErrorCode, buildMsg(message...), err).PrintErrorStack()
		os.Exit(1)
	}
}

func buildMsg(v ...interface{}) string {
	if len(v) == 0 {
		return ""
	}
	return fmt.Sprintf(fmt.Sprintf("%v", v[0]), v[1:]...)
}
