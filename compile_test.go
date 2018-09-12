package main

import (
	"math/big"
	"strconv"
	"testing"

	"github.com/jeffwilliams/calc/ast"
	"github.com/jeffwilliams/calc/compiler"
	"github.com/jeffwilliams/calc/vm"
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
		{
			"test_list",
			"li([1,5,6],1)",
			nil,
			big.NewInt(5),
		},
		{
			"call_lambda_from_var",
			"l = def(x){x+1}; l(5)",
			nil,
			big.NewInt(6),
		},
		{
			"call_lambda_as_fn_param",
			"def call_w_3(fn) fn(3); call_w_3(def(x){x+1})",
			nil,
			big.NewInt(4),
		},
		{
			"call_lambda_from_var_as_fn_param",
			"l = def(x){x+1}; def call_w_3(fn) fn(3); call_w_3(l)",
			nil,
			big.NewInt(4),
		},
		{
			"call_builtin_as_lambda_from_var",
			"a=+; a(1,2)",
			nil,
			big.NewInt(3),
		},
		{
			"call_builtin_as_lambda_as_fn_param",
			"def call_w_2_1(fn) fn(2,1); call_w_2_1(+)",
			nil,
			big.NewInt(3),
		},
		{
			"call_closure",
			"def parent(x) def(){x+1}; a = parent(5); a()",
			nil,
			big.NewInt(3),
		},
	}

	for i, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {
			m.State().Data = tc.data

			parsed, err := Parse("last line", []byte(tc.program))
			if err != nil {
				t.Fatalf("parsing failed: %v\n", err)
			}

			if _, ok := parsed.(*ast.Stmts); !ok {
				t.Fatalf("parsing result is not an AST\n")
			}

			obj, err := compiler.Compile(strconv.Itoa(i), parsed, builtinIndices, nil)
			if err != nil {
				t.Fatalf("compilation failed: %v\n", err)
			}

			linked, err := obj.Linked()
			if err != nil {
				t.Fatalf("linking failed: %v\n", err)
			}

			code := linked.Code

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

	tests := []struct {
		name     string
		programs []string
		data     []interface{}
		expected *big.Int
	}{
		{
			"fn_link",
			[]string{
				"def f(x) x",
				"f(5)",
			},
			nil,
			big.NewInt(5),
		},
		{
			"var_link",
			[]string{
				"x=2",
				"y=3",
				"x",
			},
			nil,
			big.NewInt(2),
		},
		{
			"var_link",
			[]string{
				"x=2",
				"y=3",
				"y",
			},
			nil,
			big.NewInt(3),
		},
	}

	for _, tc := range tests {
		var world *compiler.Compiled

		for i, x := range tc.programs {

			s := (*compiler.Shared)(nil)
			if sharedCode != nil {
				s = &sharedCode.Shared
			}

			code, err := compile(strconv.Itoa(i), x, builtinIndices, s)
			if err != nil {
				t.Fatalf("compilation failed: %v\n", err)
			}

			world = world.Link(code)
		}

		linked, err := world.Linked()
		if err != nil {
			t.Fatalf("linking failed: %v\n", err)
		}

		code := linked.Code

		err = m.Run(code, nil)
		if err != nil {
			t.Fatalf("execution error: %v\nexecution state at error:\n%v", err, err.(vm.ExecError).Details())
		}

		if len(m.State().Stack) == 0 {
			t.Fatalf("stack is empty after program")
		}

		if m.State().Stack.Top().(*big.Int).Cmp(tc.expected) != 0 {
			t.Fatalf("expected %v but got %v", 5, m.State().Stack.Top())

		}
	}

}
