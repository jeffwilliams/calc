package main

import (
	"fmt"
	"github.com/jeffwilliams/calc/vm"
	. "github.com/jeffwilliams/calc/vmimpl"
	"math/big"
	"testing"
)

func TestProgs(t *testing.T) {

	if _, ok := Funcs["+"]; !ok {
		t.Fatalf("Function 'add' is not registered")
	}

	builtinDesc := []vm.BuiltinDescr{
		{"+", Funcs["+"]},
	}

	builtinTable, err := vm.NewBuiltinTable(builtinDesc)
	if err != nil {
		t.Fatalf("creating builtin table failed: %v", err)
	}

	m, err := NewVM(builtinTable)
	if err != nil {
		t.Fatalf("creating vm failed: %v", err)
	}

	tests := []struct {
		name     string
		program  []vm.Instruction
		expected *big.Int
	}{
		{
			"add_2",
			[]vm.Instruction{
				I("push", big.NewInt(1)),
				I("push", big.NewInt(2)),
				I("iadd", nil),
			},
			big.NewInt(3),
		},
		{
			"add_3",
			[]vm.Instruction{
				I("push", big.NewInt(1)),
				I("push", big.NewInt(2)),
				I("push", big.NewInt(1)),
				I("iadd", nil),
				I("iadd", nil),
			},
			big.NewInt(4),
		},
		{
			"call_fn",
			[]vm.Instruction{
				// Call "add 1" with argument 5
				I("push", big.NewInt(5)),
				I("call", 5),

				// Continue after call and add 2 to result of call
				I("push", big.NewInt(2)),
				I("iadd", nil),
				I("halt", nil),

				// Definition of function "add 1" that adds 1 to the passed parameter
				I("enter", nil),
				I("push", big.NewInt(1)),
				I("pushparm", 0),
				I("iadd", nil),
				I("leave", 1),
				I("return", nil),
			},
			big.NewInt(8),
		},
		{
			"call_builtin_fn",
			[]vm.Instruction{
				// Call "add" with argument 5 and 6
				I("push", big.NewInt(5)),
				I("push", big.NewInt(6)),
				I("callb", CallBuiltinOperand{Index: 0, NumParms: 2}),
			},
			big.NewInt(11),
		},
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {
			err = m.Run(tc.program, nil)
			if err != nil {
				t.Fatalf("execution error: %v\nexecution state at error:\n%v", err, err.(vm.ExecError).Details())
			}

			if len(m.State().Stack) == 0 {
				t.Fatalf("stack is empty after program")
			}

			if m.State().Stack.Top().(*big.Int).Cmp(tc.expected) != 0 {
				t.Fatalf("expected %v but got %v", tc.expected, m.State().Stack.Top())
			}

			fmt.Printf("%v\n", m.State().Stack.Top())
		})
	}
}
