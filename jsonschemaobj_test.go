package jsonschemaobj_test

import (
	"github.com/deemson/jsonschemaobj"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestString(t *testing.T) {
	s := jsonschemaobj.ObjSchema{
		"type": "string",
	}
	c := s.MustCompile()
	err := c.Validate("hello")
	require.Nil(t, err)
	err = c.Validate(2)
	require.NotNil(t, err)
	require.Equal(t, "validation error: [@root: expected string, but got number]", err.Error())
	require.Len(t, err.Errors(), 1)
	pErr := err.Errors()[0]
	require.Equal(t, "@root", pErr.StringPath("."))
	require.Equal(t, "expected string, but got number", pErr.Message())
}

func TestSimpleObject(t *testing.T) {
	s := jsonschemaobj.ObjSchema{
		"type":                 "object",
		"additionalProperties": false,
		"properties": jsonschemaobj.Obj{
			"field": jsonschemaobj.Obj{
				"type": "string",
			},
		},
	}
	c := s.MustCompile()
	err := c.Validate(jsonschemaobj.Obj{"field": 1})
	require.NotNil(t, err)
	require.Equal(t, "validation error: [@root.field: expected string, but got number]", err.Error())
	require.Len(t, err.Errors(), 1)
	pErr := err.Errors()[0]
	require.Equal(t, []string{"field"}, pErr.Path())
	err = c.Validate(jsonschemaobj.Obj{"badField": "hello"})
	require.NotNil(t, err)
	require.Equal(t, `validation error: [@root: additional property "badField" not allowed]`, err.Error())
}

func TestBadJSON(t *testing.T) {
	s := jsonschemaobj.ObjSchema{
		"type": "object",
		// Obviously non-JSON serializable object
		"properties": make(chan int),
	}
	c, err := s.Compile()
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "unsupported type: chan int")
	require.Nil(t, c)
}

func TestBadJSONPanic(t *testing.T) {
	s := jsonschemaobj.ObjSchema{
		"type": "object",
		// Obviously non-JSON serializable object
		"properties": make(chan int),
	}
	defer func() {
		r := recover()
		require.NotNil(t, r)
		rErr, ok := r.(error)
		require.True(t, ok)
		require.Equal(t, "failed to marshal ObjSchema: json: unsupported type: chan int", rErr.Error())
	}()
	s.MustCompile()
}

func TestBadSchema(t *testing.T) {
	s := jsonschemaobj.ObjSchema{
		"type": "bad-type",
	}
	c, err := s.Compile()
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "I[#/type] S[#/definitions/simpleTypes/enum] value must be one of")
	require.Nil(t, c)
}