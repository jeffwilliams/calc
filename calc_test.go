package main

import (
	"math/big"
	"testing"
)

var smallFloat = big.NewFloat(0.00001)

func teql(a, b interface{}) bool {
	switch t := a.(type) {
	case BigIntList:
		return listEql(a, b)
	case BigFloatList:
		return listEql(a, b)
	case []interface{}:
		bl, ok := b.([]interface{})
		if !ok {
			return false
		}
		if len(t) != len(bl) {
			return false
		}
		for i, v := range bl {
			if !teql(v, t[i]) {
				return false
			}
		}
		return true
	default:
		return numEql(a, b)
	}
}

func listEql(a, b interface{}) bool {
	switch t := a.(type) {
	case BigIntList:
		return t.Eql(b)
	case BigFloatList:
		return t.Eql(b)
	default:
		return false
	}

	return true
}

func numEql(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil && b != nil || a != nil && b == nil {
		return false
	}

	aint, ok := a.(*big.Int)
	if ok {
		bint, ok := b.(*big.Int)
		if !ok {
			// a is an int and b is not.
			return false
		}
		if aint == nil && bint == nil {
			return true
		}
		return 0 == aint.Cmp(bint)
	} else {
		aflt := a.(*big.Float)
		bflt, ok := b.(*big.Float)
		if !ok {
			// a is a float and b is an int
			return false
		}
		if aflt == nil && bflt == nil {
			return true
		}
		acpy := big.NewFloat(0).Copy(aflt)
		return acpy.Sub(aflt, bflt).Abs(acpy).Cmp(smallFloat) < 0
	}
}

func strSliceEql(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func TestCalc(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output interface{}
		err    bool
	}{
		{
			name:   "single_int",
			input:  "1",
			output: big.NewInt(1),
		},
		{
			name:   "single_int_hex",
			input:  "0xAa",
			output: big.NewInt(170),
		},
		{
			name:   "single_int_binary",
			input:  "0b1010",
			output: big.NewInt(10),
		},
		{
			name:   "single_int_lspace",
			input:  " 1",
			output: big.NewInt(1),
		},
		{
			name:   "single_int_rspace",
			input:  "1 ",
			output: big.NewInt(1),
		},
		{
			name:   "single_int_lspace_many",
			input:  "   1",
			output: big.NewInt(1),
		},
		{
			name:   "single_int_rspace_many",
			input:  "1   ",
			output: big.NewInt(1),
		},

		{
			name:   "single_flt",
			input:  "1.1",
			output: big.NewFloat(1.1),
		},
		{
			name:   "single_flt_lspace",
			input:  " 1.2",
			output: big.NewFloat(1.2),
		},
		{
			name:   "single_flt_rspace",
			input:  "1.2 ",
			output: big.NewFloat(1.2),
		},
		{
			name:   "single_flt_lspace_many",
			input:  "   1.2",
			output: big.NewFloat(1.2),
		},
		{
			name:   "single_flt_rspace_many",
			input:  "1.2   ",
			output: big.NewFloat(1.2),
		},

		{
			name:   "add_two_ints",
			input:  "1+1",
			output: big.NewInt(2),
		},
		{
			name:   "add_two_ints_lspace",
			input:  " 1+1",
			output: big.NewInt(2),
		},
		{
			name:   "add_two_ints_rspace",
			input:  "1+1 ",
			output: big.NewInt(2),
		},
		{
			name:   "add_two_ints_inspace",
			input:  "1 +1",
			output: big.NewInt(2),
		},
		{
			name:   "add_two_ints_inspace_r",
			input:  "1+ 1",
			output: big.NewInt(2),
		},
		{
			name:   "add_two_ints_manyspace",
			input:  "  1  +  1  ",
			output: big.NewInt(2),
		},
		{
			name:   "add_ints_hex_bin",
			input:  "10+0xa+0b1010",
			output: big.NewInt(30),
		},

		{
			name:   "add_two_flts",
			input:  "1.0+1.0",
			output: big.NewFloat(2),
		},
		{
			name:   "add_two_flts_lspace",
			input:  " 1.0+1.0",
			output: big.NewFloat(2),
		},
		{
			name:   "add_two_flts_rspace",
			input:  "1.0+1.0 ",
			output: big.NewFloat(2),
		},
		{
			name:   "add_two_flts_inspace",
			input:  "1.0 +1.0",
			output: big.NewFloat(2),
		},
		{
			name:   "add_two_flts_inspace_r",
			input:  "1.0+ 1.0",
			output: big.NewFloat(2),
		},
		{
			name:   "add_two_flts_manyspace",
			input:  "  1.0  +  1.0  ",
			output: big.NewFloat(2),
		},

		{
			name:   "add_int_flt",
			input:  "1+1.0",
			output: big.NewFloat(2.0),
		},

		{
			name:   "add_flt_int",
			input:  " 1.0 +1",
			output: big.NewFloat(2.0),
		},

		{
			name:   "add_int_list_flt_list",
			input:  "[1]+[1.0]",
			output: BigFloatList{big.NewFloat(2)},
		},

		{
			name:   "add_int_flt",
			input:  "1+1.0",
			output: big.NewFloat(2.0),
		},
		{
			name:   "add_three_ints",
			input:  "1+1+2",
			output: big.NewInt(4),
		},
		{
			name:   "add_three_ints_space",
			input:  " 1 + 1 + 2 ",
			output: big.NewInt(4),
		},

		{
			name:   "add_three_flts",
			input:  "1.0+1.0+2.0",
			output: big.NewFloat(4.0),
		},
		{
			name:   "add_three_flts_space",
			input:  " 1.0 + 1.0 + 2.0 ",
			output: big.NewFloat(4.0),
		},

		{
			name:   "add_sub_ints",
			input:  "1+2-3",
			output: big.NewInt(0),
		},

		{
			name:   "add_sub_flts",
			input:  "1.1+2.1-3.2",
			output: big.NewFloat(0),
		},

		{
			name:   "single_int_paren",
			input:  "(1)",
			output: big.NewInt(1),
		},
		{
			name:   "add_two_ints_paren",
			input:  "(1+1)",
			output: big.NewInt(2),
		},

		{
			name:   "add_sub_ints_paren",
			input:  "(1+2)-3",
			output: big.NewInt(0),
		},
		{
			name:   "add_sub_ints_paren_space",
			input:  " ( 1+2)-3",
			output: big.NewInt(0),
		},
		{
			name:   "add_sub_ints_paren_space2",
			input:  "(1+2 )-3",
			output: big.NewInt(0),
		},
		{
			name:   "add_sub_ints_paren_space3",
			input:  " (1+  2) -3",
			output: big.NewInt(0),
		},

		{
			name:   "add_sub_ints_paren_next1",
			input:  "5-(3+1)",
			output: big.NewInt(1),
		},
		{
			name:   "add_sub_ints_paren_next2",
			input:  "5-(3+1)+1",
			output: big.NewInt(2),
		},
		{
			name:   "add_sub_ints_paren_next3",
			input:  "(5) - (3+1) + (1+1) - 20",
			output: big.NewInt(-17),
		},
		{
			name:   "add_sub_ints_nested_paren",
			input:  "20-(5+(5-1-(1)))",
			output: big.NewInt(12),
		},
		{
			name:   "add_sub_ints_nested_paren_flt",
			input:  "20-(5.0+(5-1-(1.0)))",
			output: big.NewFloat(12.0),
		},

		{
			name:   "mul_div_int",
			input:  "5*4/10",
			output: big.NewInt(2),
		},
		{
			name:   "mul_div_flt",
			input:  "5.0*4.0/10.0",
			output: big.NewFloat(2),
		},
		{
			name:   "prec_mul_add",
			input:  "1+2*5",
			output: big.NewInt(11),
		},
		{
			name:   "prec_mul_add_2",
			input:  "2*5+1",
			output: big.NewInt(11),
		},
		{
			name:   "prec_exp_mul_add",
			input:  "1+2*2^3",
			output: big.NewInt(17),
		},
		{
			name:   "prec_exp_mul_add",
			input:  "1+2^3*2",
			output: big.NewInt(17),
		},
		{
			name:   "prec_and_or",
			input:  "0b1&0b1|0b0&0b0",
			output: big.NewInt(1),
		},

		{
			name:   "add_as_function_int",
			input:  "+(1,2)",
			output: big.NewInt(3),
		},
		{
			name:   "add_as_function_flt",
			input:  "+(1.0,2.0)",
			output: big.NewFloat(3),
		},
		{
			name:   "add_as_function_int_list",
			input:  "+([1,1],[2,2])",
			output: BigIntList{big.NewInt(3), big.NewInt(3)},
		},
		{
			name:   "add_as_function_flt_list",
			input:  "+([1,1],[2.0,2.0])",
			output: BigFloatList{big.NewFloat(3), big.NewFloat(3)},
		},

		{
			name:   "sub_as_function_int",
			input:  "-(2,1)",
			output: big.NewInt(1),
		},
		{
			name:   "sub_as_function_flt",
			input:  "-(2.0,1.0)",
			output: big.NewFloat(1),
		},
		{
			name:   "sub_as_function_int_list",
			input:  "-([2,2],[1,1])",
			output: BigIntList{big.NewInt(1), big.NewInt(1)},
		},
		{
			name:   "sub_as_function_flt_list",
			input:  "-([2.0,2.0],[1,1])",
			output: BigFloatList{big.NewFloat(1), big.NewFloat(1)},
		},

		{
			name:   "mul_as_function_int",
			input:  "*(1,2)",
			output: big.NewInt(2),
		},
		{
			name:   "mul_as_function_flt",
			input:  "*(1.0,2.0)",
			output: big.NewFloat(2),
		},
		{
			name:   "mul_as_function_int_list",
			input:  "*([1,1],[2,2])",
			output: BigIntList{big.NewInt(2), big.NewInt(2)},
		},
		{
			name:   "mul_as_function_flt_list",
			input:  "*([1,1],[2.0,2.0])",
			output: BigFloatList{big.NewFloat(2), big.NewFloat(2)},
		},

		{
			name:   "quo_as_function_int",
			input:  "/(1,2)",
			output: big.NewInt(0),
		},
		{
			name:   "quo_as_function_flt",
			input:  "/(1.0,2.0)",
			output: big.NewFloat(0.5),
		},
		{
			name:   "quo_as_function_int_list",
			input:  "/([1,1],[2,2])",
			output: BigIntList{big.NewInt(0), big.NewInt(0)},
		},
		{
			name:   "quo_as_function_flt_list",
			input:  "/([1,1],[2.0,2.0])",
			output: BigFloatList{big.NewFloat(0.5), big.NewFloat(0.5)},
		},

		{
			name:   "exp_as_function_int",
			input:  "^(2,3)",
			output: big.NewInt(8),
		},
		{
			name:   "exp_as_function_flt",
			input:  "^(2.0,3.0)",
			output: (*big.Float)(nil),
			err:    true,
		},
		{
			name:   "exp_as_function_int_list",
			input:  "^([2,1],[3,3])",
			output: BigIntList{big.NewInt(8), big.NewInt(1)},
		},
		{
			name:   "exp_as_function_flt_list",
			input:  "^([2.0,1.0],[3.0,3.0])",
			output: BigFloatList{},
			err:    true,
		},

		{
			name:   "order_of_ops_mul",
			input:  "1+2*3",
			output: big.NewInt(7),
		},
		{
			name:   "order_of_ops_div",
			input:  "10 - 10/2 + 2",
			output: big.NewInt(7),
		},
		{
			name:   "order_of_ops_exp",
			input:  " 2*3 + 3^3 + 4*1",
			output: big.NewInt(37),
		},
		{
			name:   "order_of_ops_exp_paren",
			input:  "(2*3 + 3)^2 + 4*1",
			output: big.NewInt(85),
		},

		{
			name:   "bitwise_not",
			input:  "~1",
			output: big.NewInt(-2), // Two's compliment
		},
		{
			name:   "bitwise_not_and_and",
			input:  "7 & (~1)",
			output: big.NewInt(6),
		},
		{
			name:   "bitwise_not_and_add",
			input:  "~4 + 5",
			output: big.NewInt(0),
		},
		{
			name:   "lt",
			input:  "4<5;4<4;4<3",
			output: []interface{}{big.NewInt(1), big.NewInt(0), big.NewInt(0)},
		},
		{
			name:   "lte",
			input:  "4<=5;4<=4;4<=3",
			output: []interface{}{big.NewInt(1), big.NewInt(1), big.NewInt(0)},
		},
		{
			name:   "gt",
			input:  "4>5;4>4;4>3",
			output: []interface{}{big.NewInt(0), big.NewInt(0), big.NewInt(1)},
		},
		{
			name:   "gte",
			input:  "4>=5;4>=4;4>=3",
			output: []interface{}{big.NewInt(0), big.NewInt(1), big.NewInt(1)},
		},
		{
			name:   "eql",
			input:  "4=5;4=4",
			output: []interface{}{big.NewInt(0), big.NewInt(1)},
		},

		{
			name:   "unary_negation",
			input:  "-4",
			output: big.NewInt(-4),
		},
		{
			name:   "unary_negation_mul_add",
			input:  "-4*2 + 7",
			output: big.NewInt(-1),
		},
		{
			name:   "unary_negation_sub",
			input:  "1--1",
			output: big.NewInt(2),
		},
		{
			name:   "left_shift",
			input:  "1<<1",
			output: big.NewInt(2),
		},
		{
			name:   "right_shift",
			input:  "8>>2",
			output: big.NewInt(2),
		},
		{
			name:   "left_shift_on_float",
			input:  "8.0<<1",
			output: (*big.Int)(nil),
			err:    true,
		},
		/* commas in numbers is not supported due to ambiguety
		   with function calls with multiple parameters.
		   		{
		   			name:   "commas_in_numbers_flt",
		   			input:  "24,000.00+6,000.00",
		   			output: big.NewFloat(30000),
		   		},
		   		{
		   			name:   "commas_in_numbers_int",
		   			input:  "24,000+6,000",
		   			output: big.NewInt(30000),
		   		},
		*/
		{
			name:   "function_no_params",
			input:  "funca()",
			output: big.NewInt(555),
		},
		{
			name:   "function_number_param",
			input:  "funcb(5)",
			output: big.NewInt(555),
		},
		{
			name:   "function_number_param_space",
			input:  "funcb( 5 )",
			output: big.NewInt(555),
		},
		{
			name:   "function_expr_param",
			input:  "funcb(5*2/3)",
			output: big.NewInt(555),
		},
		{
			name:   "function_expr_param_space",
			input:  "funcb( 5*2/3 )",
			output: big.NewInt(555),
		},
		{
			name:   "function_two_params",
			input:  "funcc(1,2)",
			output: big.NewInt(555),
		},
		{
			name:   "function_two_params_space",
			input:  "funcc( 1 , 2 )",
			output: big.NewInt(555),
		},
		{
			name:   "function_def",
			input:  "def f(x,y) 5",
			output: nil,
		},
		{
			name:   "function_def_one_param",
			input:  "def f(x) 5",
			output: nil,
		},
		{
			name:   "function_def_long_param_names",
			input:  "def func (parmA,parmB) 5+parmA/parmB",
			output: nil,
		},
		{
			name:   "empty_list",
			input:  "[]",
			output: BigIntList{},
		},
		{
			name:   "one_elem_list",
			input:  "[5]",
			output: BigIntList{big.NewInt(5)},
		},
		{
			name:   "two_elem_list",
			input:  " [ 5 , 2*6 ] ",
			output: BigIntList{big.NewInt(5), big.NewInt(12)},
		},
		{
			name:  "wrong_type_list",
			input: " [ 5 , 2.0 ] ",
			err:   true,
		},
		{
			name:  "wrong_type_list_2",
			input: " [ 5.0 , 2 ] ",
			err:   true,
		},
		{
			name:   "vec_add",
			input:  "[1]+[2]",
			output: BigIntList{big.NewInt(3)},
		},
		{
			name:   "vec_sub",
			input:  "[3,4]-[2,2]",
			output: BigIntList{big.NewInt(1), big.NewInt(2)},
		},
		{
			name:   "list_map",
			input:  "map([25.0,9.0,81.0], sqrt)",
			output: BigFloatList{big.NewFloat(5), big.NewFloat(3), big.NewFloat(9)},
		},
		{
			name:   "list_reduce",
			input:  "reduce([1,2,3], +, 0)",
			output: big.NewInt(6),
		},
		{
			name:   "wrong_list_len",
			input:  "[5,2]+[3]",
			output: BigIntList{},
			err:    true,
		},
		{
			name:   "multiple_statements",
			input:  "5;6",
			output: []interface{}{big.NewInt(5), big.NewInt(6)},
		},
		{
			name:   "var_cloned_when_read",
			input:  "v=1;v+1+v",
			output: []interface{}{big.NewInt(3)},
		},
		{
			name:   "var_cloned_when_read2",
			input:  "v=1;v=1+v;v",
			output: []interface{}{nil, big.NewInt(2)},
		},
		{
			name:   "var_cloned_in_fn_call",
			input:  "v=1;def f(x) 1+x; f(v); v",
			output: []interface{}{nil, big.NewInt(2), big.NewInt(1)},
		},
		{
			name:   "var_cloned_when_read_list",
			input:  "v=[1];v+[1]+v",
			output: []interface{}{BigIntList{big.NewInt(3)}},
		},

		{
			name:   "lambda",
			input:  "map([1,2,3],def(x){x+1})",
			output: BigIntList{big.NewInt(2), big.NewInt(3), big.NewInt(4)},
		},
		{
			name:   "lambda2",
			input:  "map([1,3,5],def(x){(x+1)/2})",
			output: BigIntList{big.NewInt(1), big.NewInt(2), big.NewInt(3)},
		},

		{
			name:   "if1",
			input:  "if(9)",
			output: big.NewInt(9),
		},
		{
			name:   "if2",
			input:  "if(1,4,5)",
			output: big.NewInt(4),
		},
		{
			name:   "if3",
			input:  "if(0,4,5)",
			output: big.NewInt(5),
		},
		{
			name:   "if4",
			input:  "if(0,4,1,2,5)",
			output: big.NewInt(2),
		},
		{
			name:   "if5",
			input:  "if(0,4,0,2,5)",
			output: big.NewInt(5),
		},
		{
			name:   "closures_1",
			input:  "def clamp(y) def(x){if(x>y,y,x)}; fn=clamp(3); fn(5); fn(2)",
			output: []interface{}{nil, big.NewInt(3), big.NewInt(2)},
		},
	}

	fn0 := func() (*big.Int, error) {
		return big.NewInt(555), nil
	}
	fn1 := func(a *big.Int) (*big.Int, error) {
		return big.NewInt(555), nil
	}
	fn2 := func(a, b *big.Int) (*big.Int, error) {
		return big.NewInt(555), nil
	}

	RegisterBuiltin("funca", fn0, "")
	RegisterBuiltin("funcb", fn1, "")
	RegisterBuiltin("funcc", fn2, "")

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parsed, err := Parse("test", []byte(tc.input))
			// Uncomment the below to print pigeon debug info
			//parsed, err := Parse("test", []byte(tc.input), Debug(true))

			if err != nil && !tc.err {
				t.Fatalf("parsing '%s' failed: %v", tc.input, err)
			}

			if tc.err && err == nil {
				t.Fatalf("expected error but none occurred")
			}

			if !teql(parsed, tc.output) {
				t.Fatalf("expected '%v' (type %T) but got '%v' (type %T)", tc.output, tc.output, parsed, parsed)
			}

		})
	}
}

func firstParseErr(err error) error {
	return err.(errList)[0].(*parserError).Inner
}

func TestUndefVarInOtherwiseValidExpr(t *testing.T) {
	_, err := Parse("test", []byte("1+X"))
	if err == nil {
		t.Fatalf("no error when variable unbound")
	}
	if len(err.(errList)) > 1 {
		t.Fatalf("there is more than one error: %v", err)
	}
	if err, ok := firstParseErr(err).(ErrUnboundVar); !ok {
		t.Fatalf("error is not ErrUnboundVar: %v %T", err, err)
	}
}

func TestTwoUndefVarInOtherwiseValidExpr(t *testing.T) {
	_, err := Parse("test", []byte("1+X+y"))
	if err == nil {
		t.Fatalf("no error when variable unbound")
	}
	if len(err.(errList)) > 2 {
		t.Fatalf("there is more than two errors: %v", err)
	}
	for _, e := range err.(errList) {
		inner := e.(*parserError).Inner
		if err, ok := inner.(ErrUnboundVar); !ok {
			t.Fatalf("one of the errors is not ErrUnboundVar: %v %T", err, err)
		}
	}
}

func TestUndefVarInInvalidExpr(t *testing.T) {
	_, err := Parse("test", []byte("1-+X"))
	if err == nil {
		t.Fatalf("no error when expected")
	}
	if _, ok := firstParseErr(err).(ErrUnboundVar); len(err.(errList)) == 1 && ok {
		t.Fatalf("there is only one error and its the unbound var (parse error was not detected): %v", err)
	}
}

func TestMultipleBlocks(t *testing.T) {
	r, err := Parse("test", []byte("5;1+2"))
	if err != nil {
		t.Fatalf("parsing failed: %v", err)
	}

	l := r.([]interface{})
	v, ok := l[0].(*big.Int)
	if !ok {
		t.Fatalf("First block did not evaluate to int")
	}
	if !numEql(v, big.NewInt(5)) {
		t.Fatalf("first block has wrong value: %v", v)
	}

	v, ok = l[1].(*big.Int)
	if !ok {
		t.Fatalf("First block did not evaluate to int")
	}
	if !numEql(v, big.NewInt(3)) {
		t.Fatalf("first block has wrong value: %v", v)
	}
}

func TestSetVar(t *testing.T) {
	_, err := Parse("test", []byte("baz = 6"))
	if err != nil {
		t.Fatalf("error when setting var: %v", err)
	}

	v, err := Parse("test", []byte("baz"))
	if err != nil {
		t.Fatalf("error when reading var: %v", err)
	}

	if !numEql(v, big.NewInt(6)) {
		t.Fatalf("variable has wrong value: %v", v)
	}

}

func TestFuncDef(t *testing.T) {

	tests := []struct {
		name       string
		text       string
		paramNames []string
		body       string
		help       string
	}{
		{
			name:       "no_params",
			text:       "def fobb () 1",
			paramNames: []string{},
			body:       "1",
		},
		{
			name:       "one_params",
			text:       "def fobb (x) x+1",
			paramNames: []string{"x"},
			body:       "x+1",
		},
		{
			name:       "two_params",
			text:       "def fobb (x,why) x+why+1",
			paramNames: []string{"x", "why"},
			body:       "x+why+1",
		},
		{
			name:       "no_params_space",
			text:       "def fobb (  ) 1",
			paramNames: []string{},
			body:       "1",
		},
		{
			name:       "two_params_space",
			text:       "def fobb ( x , why ) x+why+1",
			paramNames: []string{"x", "why"},
			body:       "x+why+1",
		},
		{
			name:       "help",
			text:       "def fobb ( x , why ) \"fobbert b\" x+why+1",
			paramNames: []string{"x", "why"},
			body:       "x+why+1",
			help:       "fobbert b",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Uncomment the below to print pigeon debug info
			//_, err := Parse("test", []byte(tc.text), Debug(true))
			_, err := Parse("test", []byte(tc.text))
			if err != nil {
				t.Fatalf("error when parsing: %v", err)
			}
			f, ok := Funcs["fobb"]
			if !ok {
				t.Fatalf("function `fobb` didn't get defined")
			}

			df, ok := f.(*DefinedFunc)
			if !ok {
				t.Fatalf("function `fobb` is not a DefinedFunc")
			}

			if !strSliceEql(df.paramNames, tc.paramNames) {
				t.Fatalf("function `fobb` has wrong param: expected %v, actual %v ", tc.paramNames, df.paramNames)
			}

			if string(df.body) != tc.body {
				t.Fatalf("function `fobb` body is wrong: expected %v, actual %v ", tc.body, df.body)
			}

			if df.help != tc.help {
				t.Fatalf("function `fobb` help is wrong: expected %v, actual %v ", tc.help, df.help)
			}

		})
	}

}
