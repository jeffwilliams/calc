package ast

import (
	"fmt"
	"reflect"
)

type Expr interface{}

type Meta struct {
	Data interface{}
}

func (m Meta) GetMeta() interface{} {
	return m.Data
}

func (m *Meta) SetMeta(d interface{}) {
	m.Data = d
}

type Metaer interface {
	GetMeta() interface{}
	SetMeta(d interface{})
}

type BinaryExpr struct {
	Meta
	Op string
	X  Expr
	Y  Expr
}

type UnaryExpr struct {
	Meta
	Op string
	X  Expr
}

type FuncCall struct {
	Meta
	Name string
	Args []Expr
}

type Number struct {
	Meta
	Value interface{}
}

type List struct {
	Meta
	Type     reflect.Type
	Elements []Expr
}

type Ident struct {
	Meta
	Name string
}

type FuncDef struct {
	Meta
	// Name is "" in the case of a lambda
	Name string
	Args []string
	Help string
	Body Expr
}

type SetStmt struct {
	Meta
	Name string
	Rhs  Expr
}

type Stmts struct {
	Meta
	Stmts []interface{}
}

type Visitor func(node interface{}, depth int) bool

type WalkType int

const (
	Pre WalkType = iota
	Post
)

func Walk(v Visitor, t WalkType, node interface{}) {
	walk(v, t, node, 0)
}

func walk(v Visitor, t WalkType, node interface{}, depth int) {

	if t == Pre {
		if !v(node, depth) {
			return
		}
	}

	wk := func(node interface{}) {
		walk(v, t, node, depth+1)
	}

	switch t := node.(type) {
	case *Stmts:
		for _, s := range t.Stmts {
			wk(s)
		}
	case *BinaryExpr:
		wk(t.X)
		wk(t.Y)
	case *UnaryExpr:
		wk(t.X)
	case *FuncCall:
		for _, a := range t.Args {
			wk(a)
		}
	case *FuncDef:
		wk(t.Body)
	case *Number:
		// Leaf
	case *Ident:
		// Leaf

	default:
		panic(fmt.Sprintf("ast.Walk: unexpected node type %T", node))
	}

	if t == Post {
		if !v(node, depth) {
			return
		}
	}
}
