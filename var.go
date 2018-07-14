package main

import (
	"fmt"
	"math/big"
)

var GlobalVars = map[string]interface{}{}

// Parameters of DefinedFunctions are local vars
var LocalVars = map[string]interface{}{}

type ErrUnboundVar string

func NewErrUnboundVar(name string) ErrUnboundVar {
	return ErrUnboundVar(fmt.Sprintf("Unbound variable %s", name))
}

func (e ErrUnboundVar) Error() string {
	return string(e)
}

func Resolve(varName string) (interface{}, error) {
	v, err := ResolveStrict(varName)

	if err == nil {
		return v, nil
	}

	if v, ok := Funcs[varName]; ok {
		return v, nil
	}

	return big.NewInt(1), err
}

// ResolveStrict only resolves variables, not functions.
func ResolveStrict(varName string) (interface{}, error) {
	if v, ok := LocalVars[varName]; ok {
		v, _ = clone(v)
		return v, nil
	}

	if v, ok := GlobalVars[varName]; ok {
		v, _ = clone(v)
		return v, nil
	}

	return big.NewInt(1), NewErrUnboundVar(varName)
}

func SetGlobal(name string, val interface{}) {
	GlobalVars[name] = val
	updateAutocomplete()
}

func ClearLocals() {
	for k := range LocalVars {
		delete(LocalVars, k)
	}
}
