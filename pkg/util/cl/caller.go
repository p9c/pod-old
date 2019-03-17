package cl

import (
	"fmt"
	"runtime"
)

// Ine (cl.Ine) returns caller location in source code
var Ine = func() error {

	_, file, line, _ := runtime.Caller(1)
	return fmt.Errorf("[%s:%d]", file, line)
}
