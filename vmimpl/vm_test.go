package vmimpl

import (
	"fmt"
	"github.com/jeffwilliams/calc/vm"
	"math/big"
	"testing"
)

func TestProgs(t *testing.T) {

	m, err := NewVM()
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
				I(PushOp, big.NewInt(1)),
				I(PushOp, big.NewInt(2)),
				I(IAddOp, nil),
			},
			big.NewInt(3),
		},
		{
			"add_3",
			[]vm.Instruction{
				I(PushOp, big.NewInt(1)),
				I(PushOp, big.NewInt(2)),
				I(PushOp, big.NewInt(1)),
				I(IAddOp, nil),
				I(IAddOp, nil),
			},
			big.NewInt(4),
		},
		{
			"call_fn",
			[]vm.Instruction{
				// Call "add 1" with argument 5
				I(PushOp, big.NewInt(5)),
				I(CallOp, 5),

				// Continue after call and add 2 to result of call
				I(PushOp, big.NewInt(2)),
				I(IAddOp, nil),
				I(HaltOp, nil),

				// Definition of function "add 1" that adds 1 to the passed parameter
				I(EnterOp, nil),
				I(PushOp, big.NewInt(1)),
				I(PushParmOp, 0),
				I(IAddOp, nil),
				I(LeaveOp, 1),
				I(ReturnOp, nil),
			},
			big.NewInt(8),
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
