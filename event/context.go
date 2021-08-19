package event

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/fatih/structs"
)

type ContextType string

const (
	ContextInvalid ContextType = "invalid"
)

func GetContextType(typ string) ContextType {
	for k := range ctxreg {
		if strings.EqualFold(typ, string(k)) {
			return k
		}
	}

	return ContextInvalid
}

type Context interface {
	Type() ContextType
	Values() map[string]interface{}
	Interface() interface{}
	Validate() bool
}

type ContextValidator interface {
	Validate() bool
}

type context struct {
	typ ContextType
	v   interface{}
}

func NewContext(typ ContextType, v interface{}) Context {
	return &context{
		typ: typ,
		v:   v,
	}
}

func (ctx context) Type() ContextType {
	return ctx.typ
}

func (ctx context) Values() map[string]interface{} {
	if structs.IsStruct(ctx.v) {
		s := structs.New(ctx.v)
		s.TagName = "json"

		return s.Map()
	}

	switch vv := ctx.v.(type) {
	case map[string]interface{}:
		return vv
	case map[string]string:
		m := make(map[string]interface{}, len(vv))
		for k, v := range vv {
			m[k] = v
		}

		return m
	default:
		return map[string]interface{}{}
	}
}

func (ctx *context) Interface() interface{} {
	return ctx.v
}

func (ctx *context) Validate() bool {
	v, ok := ctx.v.(ContextValidator)
	if ok {
		return v.Validate()
	}

	return true
}

func emptyContext(typ ContextType) (Context, error) {
	ctor, ok := ctxreg[typ]
	if !ok {
		return nil, errors.New("context type unknown")
	}

	return ctor(), nil
}

func decodeContext(typ string, vals json.RawMessage) (Context, error) {
	ct := GetContextType(typ)
	if ct == ContextInvalid {
		return nil, errors.New("context type invalid")
	}

	ec, err := emptyContext(ct)
	if err != nil {
		return nil, err
	}

	nv := ec.Interface()
	if err := json.Unmarshal(vals, &nv); err != nil {
		return nil, err
	}

	ctx, ok := nv.(Context)
	if !ok {
		return nil, errors.New("invalid context")
	}
	return ctx, nil
}
