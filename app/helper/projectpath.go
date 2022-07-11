package helper

import (
	"path/filepath"
	"runtime"
)

var (
	_, path, _, _ = runtime.Caller(0)
	BasePath      = filepath.Join(filepath.Dir(path), "../..")
)
