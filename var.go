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

func tst() {
	var e error
	e = NewErrUnboundVar("m")
	_ = e
}

//var ErrUnboundVar error = fmt.Errorf("Unbound variable")

func Resolve(varName string) (interface{}, error) {
	if v, ok := LocalVars[varName]; ok {
		return v, nil
	}

	if v, ok := GlobalVars[varName]; ok {
		return v, nil
	} else {
		return big.NewInt(1), NewErrUnboundVar(varName)
	}
}

func ClearLocals() {
	for k := range LocalVars {
		delete(LocalVars, k)
	}
}
