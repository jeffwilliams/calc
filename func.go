package main

import (
	"fmt"
	"reflect"
)

type Func interface {
	Call(parms []interface{}) (result interface{}, err error)
}

type BuiltinFunc struct {
	name string
	fn   reflect.Value
	typ  reflect.Type
}

func (f BuiltinFunc) Call(parms []interface{}) (result interface{}, err error) {

	// Validate arity
	if len(parms) != f.typ.NumIn() {
		err = fmt.Errorf("Invalid number of params when calling %s: expected %d but got %d", f.name, f.typ.NumIn(), len(parms))
		return
	}

	// Validate parameter types
	for i := 0; i < f.typ.NumIn(); i++ {
		p := reflect.TypeOf(parms[i])
		t := f.typ.In(i)
		if !p.AssignableTo(t) {
			err = fmt.Errorf("Parameter %d is invalid: expected %s but got %s", i+1, t.Name(), p.Name())
			return
		}
	}

	// Make values
	vals := make([]reflect.Value, f.typ.NumIn())
	for i, p := range parms {
		vals[i] = reflect.ValueOf(p)
	}

	resultVals := f.fn.Call(vals)
	result = resultVals[0].Interface()
	if !resultVals[1].IsNil() {
		err = resultVals[1].Interface().(error)
	}

	return result, err

}

type DefinedFunc struct {
	name       string
	paramNames []string
	body       []byte
}

func (f DefinedFunc) Call(parms []interface{}) (result interface{}, err error) {
	defer ClearLocals()

	for i, parm := range parms {
		if i > len(f.paramNames) {
			break
		}
		LocalVars[f.paramNames[i]] = parm
	}

	return Parse("function call", f.body)
}

var Funcs map[string]Func = map[string]Func{}

// Create a Func that wraps the passed function `fn` and store it in the Funcs map so that it may be used in
// calculations. The created Func is returned.
func RegisterBuiltin(name string, fn interface{}) Func {

	f := &BuiltinFunc{
		name: name,
		typ:  reflect.TypeOf(fn),
		fn:   reflect.ValueOf(fn),
	}

	Funcs[f.name] = f
	return f
}

func RegisterDefined(name string, paramNames []string, body []byte) Func {

	f := &DefinedFunc{
		name:       name,
		paramNames: paramNames,
		body:       body,
	}

	Funcs[f.name] = f
	return f
}

func Call(name string, parms []interface{}) (result interface{}, err error) {
	f, ok := Funcs[name]
	if !ok {
		return nil, fmt.Errorf("No such function %s", name)
	}
	return f.Call(parms)
}

var funcParse func(filename string, b []byte, opts ...Option) (interface{}, error)

func init() {
	funcParse = Parse
}

func validateFuncDef(parms []interface{}, body []byte) error {

	// Here we basically run the function body through
	// the parser to validate it is a valid expression.
	// Unless the function uses no variables/parameters
	// in the body (i.e. is a constant expression) then
	// we will get an "unbound variable" error, which
	// is acceptable. Any other errors indicate
	// a parse error though.

	if funcParse == nil {
		return fmt.Errorf("funcParse was not set")
	}
	_, err := funcParse("function def", body)

	if err == nil {
		return nil
	}

	n := len(parms)
	el := err.(errList)
	m := len(el)
	if m < n {
		return fmt.Errorf("Not all parameters are used in the function body")
	}
	// go through the error list, and remove any that are ErrUnboundVar.
	// if the list is not empty, return it as an error.
	nl := make(errList, 0, m)
	for _, e := range el {
		inner := e.(*parserError).Inner
		if _, ok := inner.(ErrUnboundVar); !ok {
			nl = append(nl, inner)
		}
	}

	if len(nl) > 0 {
		return nl
	}

	return nil

}