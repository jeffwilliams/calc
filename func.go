package main

import (
	"fmt"
	"reflect"
)

type Func struct {
	name string
	fn   reflect.Value
	typ  reflect.Type
}

func (f Func) Call(parms []interface{}) (result interface{}, err error) {
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

var Funcs map[string]*Func = map[string]*Func{}

// Create a Func that wraps the passed function `fn` and store it in the Funcs map so that it may be used in
// calculations. The created Func is returned.
func Register(name string, fn interface{}) *Func {

	f := &Func{
		name: name,
		typ:  reflect.TypeOf(fn),
		fn:   reflect.ValueOf(fn),
	}

	Funcs[f.name] = f
	return f
}
