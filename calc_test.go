package main

import (
	"math/big"
	"testing"
)

var smallFloat = big.NewFloat(0.00001)

func numEql(a, b interface{}) bool {
	aint, ok := a.(*big.Int)
	if ok {
		bint, ok := b.(*big.Int)
		if !ok {
			// a is an int and b is not.
			return false
		}
		return 0 == aint.Cmp(bint)
	} else {
		aflt := a.(*big.Float)
		bflt, ok := b.(*big.Float)
		if !ok {
			// a is a float and b is an int
			return false
		}
		//fmt.Printf("numEql: both is floats: (%v, %v). cmp == %v\n", aflt, bflt, aflt.Cmp(bflt))
		acpy := big.NewFloat(0).Copy(aflt)
		return acpy.Sub(aflt, bflt).Abs(acpy).Cmp(smallFloat) < 0
		//return 0 == aflt.Cmp(bflt)
	}
}

func TestCalc(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output interface{}
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
			input:  "1+1",
			output: big.NewInt(2),
		},
		{
			name:   "add_two_flts_lspace",
			input:  " 1+1",
			output: big.NewInt(2),
		},
		{
			name:   "add_two_flts_rspace",
			input:  "1+1 ",
			output: big.NewInt(2),
		},
		{
			name:   "add_two_flts_inspace",
			input:  "1 +1",
			output: big.NewInt(2),
		},
		{
			name:   "add_two_flts_inspace_r",
			input:  "1+ 1",
			output: big.NewInt(2),
		},
		{
			name:   "add_two_flts_manyspace",
			input:  "  1  +  1  ",
			output: big.NewInt(2),
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
			name:   "commas_in_numbers_flt",
			input:  "24,000.00+6,000.00",
			output: big.NewFloat(30000),
		},
		{
			name:   "commas_in_numbers_int",
			input:  "24,000+6,000",
			output: big.NewInt(30000),
		},

		{
			name:   "function_no_params",
			input:  "func()",
			output: big.NewInt(555),
		},
		{
			name:   "function_number_param",
			input:  "func(5)",
			output: big.NewInt(555),
		},
		{
			name:   "function_number_param_space",
			input:  "func( 5 )",
			output: big.NewInt(555),
		},
		{
			name:   "function_expr_param",
			input:  "func(5*2/3)",
			output: big.NewInt(555),
		},
		{
			name:   "function_expr_param_space",
			input:  "func( 5*2/3 )",
			output: big.NewInt(555),
		},
		{
			name:   "function_two_params",
			input:  "func(1,2)",
			output: big.NewInt(555),
		},
		{
			name:   "function_two_params_space",
			input:  "func( 1 , 2 )",
			output: big.NewInt(555),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parsed, err := Parse("test", []byte(tc.input))

			if err != nil {
				t.Fatalf("parsing '%s' failed: %v", tc.input, err)
			}

			//t.Logf("type: %v", reflect.TypeOf(parsed))
			if !numEql(parsed, tc.output) {
				t.Fatalf("expected '%v' (type %T) but got '%v' (type %T)", tc.output, tc.output, parsed, parsed)
			}

		})
	}
}
