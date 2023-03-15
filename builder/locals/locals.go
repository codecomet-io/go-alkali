package locals

import (
	"strconv"
)

var localSingleton *LocalsRegister

type LocalsRegister struct {
	x       int
	locals  map[string]string
	reverse map[string]string
}

func Get(path string) string {
	if localSingleton == nil {
		localSingleton = &LocalsRegister{
			x:       0,
			locals:  make(map[string]string),
			reverse: make(map[string]string),
		}
	}
	if _, ok := localSingleton.locals[path]; !ok {
		localSingleton.x++
		localSingleton.locals[path] = "folder" + strconv.Itoa(localSingleton.x)
		localSingleton.reverse["folder"+strconv.Itoa(localSingleton.x)] = path
	}

	return localSingleton.locals[path]
}

func Dump() map[string]string {
	if localSingleton == nil {
		localSingleton = &LocalsRegister{}
	}
	return localSingleton.reverse
}
