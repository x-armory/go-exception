package ex

import (
	"bytes"
	"fmt"
	"gopkg.in/go-playground/validator.v8"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

const (
	WrappedErrorCode    = "WRAPPED"
	AssertionErrorCode  = "ASSERTION_FAILED"
	ValidationErrorCode = "VALIDATION_FAILED"
	ExitErrorCode       = "EXIT"
)

type exceptionClass struct {
	code         string
	message      string
	stackTrace   []*runtime.Frame
	cause        *exceptionClass
	isWrapped    bool
	wrappedCause interface{}
}

func Exception(code string, message string, cause *exceptionClass) *exceptionClass {
	return buildErrorWithCallerStack(stackStartDepth, code, message, cause)
}

// wap any object without stack as root cause
func Wrap(err interface{}) *exceptionClass {
	if e, ok := err.(*exceptionClass); ok {
		return e
	}
	if e, ok := err.(exceptionClass); ok {
		return &e
	}
	e := buildErrorWithCallerStack(stackStartDepth, WrappedErrorCode, "", err)
	e.isWrapped = true
	return e
}

func (e *exceptionClass) Code() string {
	return e.code
}
func (e *exceptionClass) Message() string {
	return e.message
}
func (e *exceptionClass) WrappedCause() interface{} {
	return e.wrappedCause
}
func (e *exceptionClass) Error() string {
	return "<" + e.code + "> " + e.message
}
func (e *exceptionClass) IsWrapped() bool {
	return e.isWrapped
}
func (e *exceptionClass) StackTraceString() string {
	var buffer bytes.Buffer
	e.addStackTraceMessage(&buffer)
	return buffer.String()
}
func (e *exceptionClass) PrintErrorStack() {
	println(e.StackTraceString())
}
func (e *exceptionClass) RootCause() *exceptionClass {
	if e.cause == nil {
		return e
	} else {
		return e.cause.RootCause()
	}
}
func (e *exceptionClass) Cause() *exceptionClass {
	return e.cause
}
func (e *exceptionClass) StackTrace() []*runtime.Frame {
	return e.stackTrace
}
func (e *exceptionClass) Throw() {
	panic(e)
}
func (e *exceptionClass) addStackTraceMessage(buffer *bytes.Buffer) {
	buffer.WriteString("Exception: ")
	buffer.WriteString(e.Error())
	for _, t := range e.stackTrace {
		buffer.WriteString("\n   at: ")
		buffer.WriteString(t.File)
		buffer.WriteString(":")
		buffer.WriteString(strconv.Itoa(t.Line))
		buffer.WriteString(" ")
		buffer.WriteString(methodNameReg.FindString(t.Function))
		buffer.WriteString("()")
	}
	if e.IsWrapped() && e.wrappedCause != nil && e.wrappedCause != "" {
		buffer.WriteString(fmt.Sprintf("\ncaused by %+v", e.wrappedCause))
	}
	if e.cause != nil {
		buffer.WriteString("\ncaused by ")
		e.cause.addStackTraceMessage(buffer)
	}
}

func buildErrorWithCallerStack(skip int, code string, message string, cause interface{}) *exceptionClass {
	if code == "" {
		code = WrappedErrorCode
	}
	err := &exceptionClass{code: code, message: message, stackTrace: []*runtime.Frame{}}
	depth := rootErrorStackMaxDepth
	if cause != nil {
		if e2, ok := cause.(*exceptionClass); ok {
			err.cause = e2
			depth = middleErrorStackMaxDepth + 1
		} else if e2, ok := cause.(validator.ValidationErrors); ok {
			buff := bytes.NewBufferString("")
			i := 0
			for _, err := range e2 {
				i++
				if i > 1 {
					buff.WriteString("; ")
				}
				buff.WriteString("field ")
				buff.WriteString(strings.ToLower(err.Field[0:1]))
				buff.WriteString(err.Field[1:])
				buff.WriteString(" ")
				buff.WriteString(err.Tag)
			}
			err.code = ValidationErrorCode
			err.message = buff.String()
			err.isWrapped = true
			err.wrappedCause = cause
		} else if e2, ok := cause.(string); ok {
			if err.message == "" {
				err.message = e2
			}
		} else {
			if err.message == "" {
				err.message = fmt.Sprintf("%v", cause)
				if len(err.message) > 200 {
					err.message = err.message[:200] + "..."
				}
			}
			err.isWrapped = true
			err.wrappedCause = cause
		}
	}
	stack := make([]uintptr, depth)
	count := runtime.Callers(skip+1, stack[:])
	stack = stack[:count]
	frames := runtime.CallersFrames(stack)
	for true {
		if frame, more := frames.Next(); !more {
			break
		} else {
			err.stackTrace = append(err.stackTrace, &frame)
		}
	}
	return err
}

var methodNameReg = regexp.MustCompile(`(\w+[.])*\w+$`)

const (
	stackStartDepth = 2
)

var rootErrorStackMaxDepth = 32
var middleErrorStackMaxDepth = 5

func SetRootErrorStackMaxDepth(depth int) {
	rootErrorStackMaxDepth = depth
}

func SetMiddleErrorStackMaxDepth(depth int) {
	middleErrorStackMaxDepth = depth
}
