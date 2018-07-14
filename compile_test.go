package main

import (
	"github.com/jeffwilliams/calc/ast"
	"github.com/jeffwilliams/calc/compiler"
	"github.com/jeffwilliams/calc/vm"
	"math/big"
	"testing"
)

func TestCompiledProgs(t *testing.T) {

	m, builtinIndices, err := NewVM()
	if err != nil {
		t.Fatalf("creating vm failed: %v", err)
	}

	tests := []struct {
		name     string
		program  string
		data     []interface{}
		expected *big.Int
	}{
		{
			"add_1_2",
			"1+2",
			nil,
			big.NewInt(3),
		},
		{
			"add_1_2_1",
			"1+2+1",
			nil,
			big.NewInt(4),
		},
		{
			"func_1",
			"def f(x) x; f(4)",
			nil,
			big.NewInt(4),
		},
		{
			"func_2",
			"def f(x,y) x; f(6+1,2)",
			nil,
			big.NewInt(7),
		},
		{
			"func_3",
			"def f(x) x*2; f(3)",
			nil,
			big.NewInt(6),
		},
		{
			// The parameter x is referenced twice. This tests if the parameter
			// is properly duplicated before being used in operations. If not
			// the result will be 8 instead of 6.
			"func_param_ref_twice",
			"def f(x) x*2+x; f(2)",
			nil,
			big.NewInt(6),
		},
		{
			"set_var_simple",
			"y = 5; y",
			nil,
			big.NewInt(5),
		},
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {
			m.State().Data = tc.data

			parsed, err := Parse("last line", []byte(tc.program))
			if err != nil {
				t.Fatalf("parsing failed: %v\n", err)
			}

			if _, ok := parsed.(*ast.Stmts); !ok {
				t.Fatalf("parsing result is not an AST\n")
			}

			obj, err := compiler.Compile(parsed, builtinIndices, nil)
			if err != nil {
				t.Fatalf("compilation failed: %v\n", err)
			}

			code, err := obj.Linked()
			if err != nil {
				t.Fatalf("linking failed: %v\n", err)
			}

			err = m.Run(code, nil)
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

func TestIncrementalLink(t *testing.T) {
	m, builtinIndices, err := NewVM()
	if err != nil {
		t.Fatalf("creating vm failed: %v", err)
	}

	// Define func, and link it to new code
	code1, err := compile("def f(x) x", builtinIndices, nil)
	t.Fatalf("%v", err)

	code2, err := compile("f(5)", builtinIndices, &code1.Shared)
	t.Fatalf("%v", err)

	linked, err := code2.Linked()
	if err != nil {
		t.Fatalf("linking failed: %v\n", err)
	}

	err = m.Run(linked, nil)
	if err != nil {
		t.Fatalf("execution error: %v\nexecution state at error:\n%v", err, err.(vm.ExecError).Details())
	}

	if len(m.State().Stack) == 0 {
		t.Fatalf("stack is empty after program")
	}

	if m.State().Stack.Top().(*big.Int).Cmp(big.NewInt(5)) != 0 {
		t.Fatalf("expected %v but got %v", 5, m.State().Stack.Top())
	}

}
