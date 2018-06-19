package main

import (
	"fmt"
	"math/big"
	"reflect"
)

type Func interface {
	Call(parms []interface{}) (result interface{}, err error)
	Help() string
}

type BuiltinFunc struct {
	name string
	help string
	fn   reflect.Value
	typ  reflect.Type
}

func (f BuiltinFunc) Call(parms []interface{}) (result interface{}, err error) {

	// Validate arity
	if !f.typ.IsVariadic() {
		if len(parms) != f.typ.NumIn() {
			err = fmt.Errorf("Invalid number of params when calling %s: expected %d but got %d", f.name, f.typ.NumIn(), len(parms))
			return
		}

		// Validate parameter types
		for i := 0; i < f.typ.NumIn(); i++ {
			p := reflect.TypeOf(parms[i])
			t := f.typ.In(i)
			if !p.AssignableTo(t) {
				// try to upcast from int to float
				if pi, ok := parms[i].(*big.Int); ok {
					flt := big.NewFloat(0)
					flt.SetInt(pi)
					parms[i] = flt
					p = reflect.TypeOf(parms[i])
				}

				if !p.AssignableTo(t) {
					err = fmt.Errorf("Parameter %d is invalid: expected %s but got %s", i+1, t, p)
					return
				}
			}
		}
	}

	// Make values
	vals := make([]reflect.Value, len(parms))
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

func (f BuiltinFunc) Help() string {
	return f.help
}

func (f BuiltinFunc) NumParams() int {
	if f.typ.IsVariadic() {
		return -1
	} else {
		return f.typ.NumIn()
	}
}

type DefinedFunc struct {
	name       string
	help       string
	paramNames []string
	body       []byte
	bound      map[string]interface{}
}

func (f DefinedFunc) Call(parms []interface{}) (result interface{}, err error) {
	defer ClearLocals()

	if f.bound != nil {
		for i, bvar := range f.bound {
			LocalVars[i] = bvar
		}
	}

	for i, parm := range parms {
		if i > len(f.paramNames) {
			break
		}
		LocalVars[f.paramNames[i]] = parm
	}

	return Parse("function call", f.body)
}

func (f DefinedFunc) Help() string {
	return f.help
}

func (f DefinedFunc) NumParams() int {
	return len(f.paramNames)
}

var Funcs map[string]Func = map[string]Func{}

// Create a Func that wraps the passed function `fn` and store it in the Funcs map so that it may be used in
// calculations. The created Func is returned.
func RegisterBuiltin(name string, fn interface{}, help string) Func {

	f := &BuiltinFunc{
		name: name,
		help: help,
		typ:  reflect.TypeOf(fn),
		fn:   reflect.ValueOf(fn),
	}

	Funcs[f.name] = f

	updateAutocomplete()

	return f
}

func RegisterDefined(name string, paramNames []string, body []byte, help string) Func {

	f := &DefinedFunc{
		name:       name,
		help:       help,
		paramNames: paramNames,
		body:       body,
	}

	Funcs[f.name] = f

	updateAutocomplete()

	return f
}

func Call(name string, parms []interface{}) (result interface{}, err error) {
	f, ok := Funcs[name]
	if ok {
		return f.Call(parms)
	}

	v, err := ResolveStrict(name)
	if err == nil {
		if f, ok := v.(Func); ok {
			result, err = f.Call(parms)
			//result = clone(result)
			return
		}
	}

	return nil, fmt.Errorf("No such function %s", name)
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
