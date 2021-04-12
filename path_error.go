package jsonschemaobj

import (
	"fmt"
	"strings"
)

type PathError struct {
	path    []string
	message string
}

func (e *PathError) Path() []string {
	path := make([]string, len(e.path))
	for index, pathPart := range e.path {
		path[index] = pathPart
	}
	return path
}

func (e *PathError) Message() string {
	return e.message
}

func (e *PathError) StringPath(separator string) string {
	path := append([]string{"@root"}, e.path...)
	return strings.Join(path, separator)
}

func (e *PathError) Error() string {
	path := e.StringPath(".")
	return fmt.Sprintf(`%s: %s`, path, e.message)
}
