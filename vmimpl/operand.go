package vmimpl

import (
	"bytes"
	"fmt"

	"github.com/jeffwilliams/calc/vm"
)

type CallBuiltinOperand struct {
	Index, NumParms int
}

func (op CallBuiltinOperand) StringWithState(s *vm.State) string {
	return fmt.Sprintf("builtin '%s', num parms %d", s.Builtins.Name(op.Index), op.NumParms)
}

type CopyStackOperand struct {
	Offset, Len int
}

func (op CopyStackOperand) StringWithState(s *vm.State) string {
	return fmt.Sprintf("offset (from end) %d, len %d", op.Offset, op.Len)
}

// LambdaClosureOperand represents a pointer to a lambda and it's corresponding closure environment
type LambdaClosureOperand struct {
	// LambdaAddr points into the code segment
	LambdaAddr int
	// ClosureEnv points into the data segment
	ClosureEnv int
}

func (op LambdaClosureOperand) StringWithState(s *vm.State) string {
	return fmt.Sprintf("lambda at %d, closure env at %d", op.LambdaAddr, op.ClosureEnv)
}

type TableOperand []interface{}

func (op TableOperand) StringWithState(s *vm.State) string {
	var buf bytes.Buffer

	fmt.Fprintf(&buf, "[")
	for i, v := range op {
		fmt.Fprintf(&buf, "%d->%v, ", i, v)
	}
	fmt.Fprintf(&buf, "]")

	return buf.String()
}
