package ast

import (
	"fmt"
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

type Parent struct {
	parent Parenter
}

func (p Parent) GetParent() Parenter {
	return p.parent
}

func (p *Parent) SetParent(parent Parenter) {
	p.parent = parent
}

type Parenter interface {
	GetParent() Parenter
	SetParent(parent Parenter)
}

type BinaryExpr struct {
	Meta
	Parent
	Op string
	X  Expr
	Y  Expr
}

type UnaryExpr struct {
	Meta
	Parent
	Op string
	X  Expr
}

type FuncCall struct {
	Meta
	Parent
	Name string
	Args []Expr
}

type Number struct {
	Meta
	Parent
	Value interface{}
}

type List struct {
	Meta
	Parent
	Elements []Expr
}

type Ident struct {
	Meta
	Parent
	Name string
}

type FuncDef struct {
	Meta
	Parent
	// Name is "" in the case of a lambda
	Name string
	Args []string
	Help string
	Body Expr
}

type SetStmt struct {
	Meta
	Parent
	Name string
	Rhs  Expr
}

type Stmts struct {
	Meta
	Parent
	Stmts []interface{}
}

type Visitor func(node interface{}, depth int) bool

type WalkType int

const (
	Pre WalkType = iota
	Post
)

func SetParent(node interface{}, parent Parenter) {
	node.(Parenter).SetParent(parent)
}

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
	case *SetStmt:
		wk(t.Rhs)
	case *List:
		for _, v := range t.Elements {
			wk(v)
		}
	default:
		panic(fmt.Sprintf("ast.Walk: unexpected node type %T", node))
	}

	if t == Post {
		if !v(node, depth) {
			return
		}
	}
}

type Criteria func(Parenter) bool

func Ancestor(node Parenter, crit Criteria) Parenter {
	for ; node != nil; node = node.GetParent() {
		if crit(node) {
			return node
		}
	}
	return nil

}
