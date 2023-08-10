package sderr

import (
	"runtime"
	"strings"

	"github.com/rotisserie/eris"
)

type StackFrame = eris.StackFrame

func StackOf(err error) []StackFrame {
	if err == nil {
		return []StackFrame{}
	}
	rawFrames := eris.StackFrames(err)
	if len(rawFrames) <= 0 {
		return []StackFrame{}
	}
	var frames []StackFrame
	callersFrames := runtime.CallersFrames(rawFrames)
	for {
		callerFrames, more := callersFrames.Next()
		i := strings.LastIndex(callerFrames.Function, "/")
		name := callerFrames.Function[i+1:]
		frames = append(frames, StackFrame{
			Name: name,
			File: callerFrames.File,
			Line: callerFrames.Line,
		})
		if !more {
			break
		}
	}
	return frames
}
