package jsonschemaobj

import (
	"fmt"
	"github.com/santhosh-tekuri/jsonschema"
	"strings"
)

type ValidationError struct {
	errors []*PathError
}

func (e *ValidationError) Errors() []*PathError {
	errors := make([]*PathError, len(e.errors))
	for index, err := range e.errors {
		errors[index] = err
	}
	return errors
}

func (e *ValidationError) Error() string {
	errStrings := make([]string, len(e.errors))
	for index, err := range e.errors {
		errStrings[index] = err.Error()
	}
	return fmt.Sprintf("validation error: [%s]", strings.Join(errStrings, ", "))
}

func instancePtrToPath(ptr string) []string {
	return strings.Split(ptr, "/")[1:]
}

func flattenErr(err *jsonschema.ValidationError) []*PathError {
	var pathErrors []*PathError
	if err.Causes == nil {
		var message string
		if strings.HasPrefix(err.Message, "additionalProperties") {
			message = strings.Replace(err.Message, "additionalProperties", "additional property", 1)
		} else {
			message = err.Message
		}
		pathErrors = append(pathErrors, &PathError{
			path:    instancePtrToPath(err.InstancePtr),
			message: message,
		})
	} else {
		for _, cause := range err.Causes {
			childErrors := flattenErr(cause)
			for _, childError := range childErrors {
				pathErrors = append(pathErrors, childError)
			}
		}
	}
	return pathErrors
}
