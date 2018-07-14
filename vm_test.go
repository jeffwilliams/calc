package main

import (
	"github.com/jeffwilliams/calc/vm"
	"github.com/jeffwilliams/calc/vmimpl"
	"math/big"
	"testing"
)

var I = vmimpl.I

func TestProgs(t *testing.T) {

	if _, ok := Funcs["+"]; !ok {
		t.Fatalf("Function 'add' is not registered")
	}

	m, builtinIndices, err := NewVM()
	if err != nil {
		t.Fatalf("creating vm failed: %v", err)
	}

	tests := []struct {
		name     string
		program  []vm.Instruction
		data     []interface{}
		expected *big.Int
	}{
		{
			"add_2",
			[]vm.Instruction{
				I("push", big.NewInt(1)),
				I("push", big.NewInt(2)),
				I("iadd", nil),
			},
			nil,
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
			nil,
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
			nil,
			big.NewInt(8),
		},
		{
			"call_builtin_fn",
			[]vm.Instruction{
				// Call "add" with argument 5 and 6
				I("push", big.NewInt(5)),
				I("push", big.NewInt(6)),
				I("callb", vmimpl.CallBuiltinOperand{Index: builtinIndices["+"], NumParms: 2}),
			},
			nil,
			big.NewInt(11),
		},
		{
			"call_builtin_fn_div",
			[]vm.Instruction{
				// Call "div" with argument 6 and 3
				I("push", big.NewInt(3)),
				I("push", big.NewInt(6)),
				I("callb", vmimpl.CallBuiltinOperand{Index: builtinIndices["/"], NumParms: 2}),
			},
			nil,
			big.NewInt(2),
		},
		{
			"call_indirect_fn",
			[]vm.Instruction{
				// Call "add 1" with argument 5, indirectly by the contents of variable 0
				I("push", big.NewInt(5)),
				I("calli", 0),

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
			[]interface{}{
				// Offset 0: a variable containing the address of function "add 1"
				5,
			},
			big.NewInt(8),
		},
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {
			m.State().Data = tc.data
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
		})
	}
}
