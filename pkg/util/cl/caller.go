package cl

import (
	"fmt"
	"runtime"
)

// Ine (cl.Ine) prefixes error string with  location in source code
var Ine = func(e *error) error {
	s := (*e).Error()
	_, file, line, ok := runtime.Caller(1)
	if ok {
		*e = fmt.Errorf("[%s:%d]:'%s'", file, line, s)
	}
	return *e
}
