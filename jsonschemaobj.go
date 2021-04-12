package jsonschemaobj

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/santhosh-tekuri/jsonschema"
	"strings"
)

type (
	Obj       = map[string]interface{}
	Arr       = []interface{}
	ObjSchema map[string]interface{}
)

func (s ObjSchema) Compile() (*CompiledObjSchema, error) {
	schemaData, err := json.Marshal(s)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal ObjSchema")
	}
	compiler := jsonschema.NewCompiler()
	err = compiler.AddResource("", strings.NewReader(string(schemaData)))
	if err != nil {
		return nil, errors.Wrap(err, "failed to add resource to jsonschema compiler")
	}
	schema, err := compiler.Compile("")
	if err != nil {
		return nil, errors.Wrap(err, "failed to compile jsonschema")
	}
	return &CompiledObjSchema{
		schema: schema,
	}, nil
}

func (s ObjSchema) MustCompile() *CompiledObjSchema {
	compiledObjSchema, err := s.Compile()
	if err != nil {
		panic(err)
	}
	return compiledObjSchema
}

type CompiledObjSchema struct {
	schema *jsonschema.Schema
}

func (s *CompiledObjSchema) Validate(i interface{}) *ValidationError {
	err := s.schema.ValidateInterface(i)
	if err != nil {
		return &ValidationError{flattenErr(err.(*jsonschema.ValidationError))}
	}
	return nil
}