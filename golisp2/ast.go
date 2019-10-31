package golisp2

import (
	"fmt"
	"strings"
)

type (
	// Expr is the fundamental unit of lisp - it represents anything that can be
	// evaluated to a value.
	Expr interface {
		Eval(ExprContext) Value
	}

	// ExprContext is the context on evaluation. It contains a resolvable set of
	// identifiers->values that can be chained.
	ExprContext struct {
		vals map[string]Value
	}

	// Value represents any arbitrary value within the lisp interpreting
	// environment.
	Value interface {
		Expr

		// PrintStr returns a printable version of the value. Note that this is not
		// the same as casting to a string.
		PrintStr() string
	}

	// NumberValue is a representation of a number value within the interpreted
	// environment.
	NumberValue struct {
		val float64
	}

	// NilValue is a representation of an null value within the interpreted
	// environment.
	NilValue struct {
	}

	// StringValue is a representation of a string within the interpreted
	// environment.
	StringValue struct {
		val string
	}

	// BoolValue is a representation of a boolean within the interpreted
	// environment.
	BoolValue struct {
		val bool
	}

	// FuncValue is a representation of a basic function within the interpreted
	// environment.
	FuncValue struct {
		// ques (bs): should this basic function type exist outside of this context?
		// Maybe.

		fn func(ExprContext, ...Expr) (Value, error)
	}

	// CellValue is a representation of a pair of values within the interpreted
	// environment. This can be composed to represent lists with standard car/cdr
	// operators.
	CellValue struct {
		left, right Value
	}

	// ques (bs): what other primitives, if any, do I want to add here? I'd like
	// to have maps, preferably with raw syntax support, and perhaps more
	// conventional vectors. Typed vectors would be great, but require types
	// before becoming practical.

)

func NewNumberValue(v float64) *NumberValue {
	return &NumberValue{
		val: v,
	}
}

func (nv *NumberValue) PrintStr() string {
	return fmt.Sprintf("%f", nv.val)
}

func (nv *NumberValue) Eval(ExprContext) Value {
	return nv
}

func (nv *NumberValue) Get() float64 {
	return nv.val
}

func NewNilValue() *NilValue {
	// note (bs): given that there are no values here to modify, consider using a
	// singleton here
	return &NilValue{}
}

func (nv *NilValue) PrintStr() string {
	return "nil"
}

func (nv *NilValue) Eval(ExprContext) Value {
	// note (bs): not sure about this. In general, I feel like eval needs to be
	// more intelligent
	return nv
}

func (nv *NilValue) Get() interface{} {
	return nil
}

func NewStringValue(str string) *StringValue {
	return &StringValue{
		val: str,
	}
}

func (nv *StringValue) PrintStr() string {
	return nv.val
}

func (nv *StringValue) Eval(ExprContext) Value {
	return nv
}

func (nv *StringValue) Get() string {
	return nv.val
}

func NewBoolValue(v bool) *BoolValue {
	// todo (bs): consider making this return singleton values
	return &BoolValue{
		val: v,
	}
}

func (nv *BoolValue) PrintStr() string {
	return fmt.Sprintf("%t", nv.val)
}

func (nv *BoolValue) Eval(ExprContext) Value {
	return nv
}

func (nv *BoolValue) Get() bool {
	return nv.val
}

func NewFuncValue(fn func(ExprContext, ...Expr) (Value, error)) *FuncValue {
	return &FuncValue{
		fn: fn,
	}
}

func (nv *FuncValue) PrintStr() string {
	// note (bs): probably want to customize this to print some details about the
	// function itself
	return fmt.Sprintf("<func>")
}

func (nv *FuncValue) Eval(ExprContext) Value {
	return nv
}

func (nv *FuncValue) Get() func(ExprContext, ...Expr) (Value, error) {
	return nv.fn
}

func NewCellValue(left, right Value) *CellValue {
	return &CellValue{
		left:  left,
		right: right,
	}
}

func (nv *CellValue) PrintStr() string {
	l, r := nv.left, nv.right
	if l == nil {
		l = NewNilValue()
	}
	if r == nil {
		r = NewNilValue()
	}
	// todo (bs): if second cell is a node, treat this different
	return fmt.Sprintf("(%s . %s)", l.PrintStr(), r.PrintStr())
}

func (nv *CellValue) Eval(ExprContext) Value {
	return nv
}

func (nv *CellValue) Get() (left, right Value) {
	return nv.left, nv.right
}

// note (bs): I'm o.k. with this for the immediate future, but functions should
// likely exist in a separate file. Let's worry about that later: I think
// packages should be broken up and all this should be moved out.

func addFn(c ExprContext, exprs ...Expr) (Value, error) {
	total := float64(0)
	for _, e := range exprs {
		v := e.Eval(c)
		asNum, isNum := v.(*NumberValue)
		if !isNum {
			// note (bs): eventually, try to make a version of this error that's more
			// portable, obvious, and a little more resilient to nil values.
			return nil, fmt.Errorf("non-number value in add: %v", asNum.PrintStr())
		}
		total += asNum.Get()
	}
	return &NumberValue{
		val: total,
	}, nil
}

func subFn(c ExprContext, exprs ...Expr) (Value, error) {
	// ques (bs): should I still enforce minimum airity requirements here? I'm
	// sorta inclined to say yes; but not sure how much I care about that right
	// now. Particularly: that seems to get into deeper questions of type
	// enforcement. Something like this could just be reduced to a set of number
	// values, and an error returned if
	//
	// That all seems like a "later" task. I'd like to just grind a bit on the
	// core language; some better limitations or even

	total := float64(0)
	for i, e := range exprs {
		v := e.Eval(c)
		asNum, isNum := v.(*NumberValue)
		if !isNum {
			// note (bs): eventually, try to make a version of this error that's more
			// portable, obvious, and a little more resilient to nil values.
			return nil, fmt.Errorf("non-number value in add: %v", v.PrintStr())
		}
		if i == 0 {
			total = asNum.Get()
		} else {
			total -= asNum.Get()
		}
	}

	return &NumberValue{
		val: total,
	}, nil
}

func concatFn(c ExprContext, exprs ...Expr) (Value, error) {
	var sb strings.Builder
	for _, e := range exprs {
		v := e.Eval(c)
		asStr, isStr := v.(*StringValue)
		if !isStr {
			return nil, fmt.Errorf("non-number value in add: %v", v.PrintStr())
		}
		sb.WriteString(asStr.Get())
	}
	return &StringValue{
		val: sb.String(),
	}, nil
}

func consFn(c ExprContext, exprs ...Expr) (Value, error) {
	if len(exprs) > 2 {
		return nil, fmt.Errorf("cons expects 0-2 argument; got %d", len(exprs))
	}
	v1 := exprs[0].Eval(c)
	v2 := exprs[1].Eval(c)
	return NewCellValue(v1, v2), nil
}

func carFn(c ExprContext, exprs ...Expr) (Value, error) {
	if len(exprs) != 1 {
		return nil, fmt.Errorf("car expects 1 argument; got %d", len(exprs))
	}
	v := exprs[0].Eval(c)
	asNode, isNode := v.(*CellValue)
	if !isNode {
		// note (bs): this was already commented on elsewhere, but I don't think
		// this is quite right. Need a better way to assemble type-error messages.
		return nil, fmt.Errorf("car expects a cell type, got %v", asNode)
	}
	leftV, _ := asNode.Get()
	return leftV, nil
}

func cdrFn(c ExprContext, exprs ...Expr) (Value, error) {
	if len(exprs) != 1 {
		return nil, fmt.Errorf("cdr expects 1 argument; got %d", len(exprs))
	}
	v := exprs[0].Eval(c)
	asNode, isNode := v.(*CellValue)
	if !isNode {
		return nil, fmt.Errorf("cdr expects a cell type, got %v", asNode)
	}
	_, rightV := asNode.Get()
	return rightV, nil
}

func andFn(c ExprContext, exprs ...Expr) (Value, error) {
	if len(exprs) == 0 {
		return nil, fmt.Errorf("and expects at least 1 argument; got %d", len(exprs))
	}
	total := true
	for _, e := range exprs {
		v := e.Eval(c)
		asBool, isBool := v.(*BoolValue)
		if !isBool {
			return nil, fmt.Errorf("and expects bool types, got %v", v)
		}
		total = total && asBool.Get()
	}
	return NewBoolValue(total), nil
}

func orFn(c ExprContext, exprs ...Expr) (Value, error) {
	if len(exprs) == 0 {
		return nil, fmt.Errorf("or expects at least 1 argument; got %d", len(exprs))
	}
	total := false
	for _, e := range exprs {
		v := e.Eval(c)
		asBool, isBool := v.(*BoolValue)
		if !isBool {
			return nil, fmt.Errorf("or expects bool types, got %v", v)
		}
		total = total || asBool.Get()
	}
	return NewBoolValue(total), nil
}

func notFn(c ExprContext, exprs ...Expr) (Value, error) {
	if len(exprs) != 1 {
		return nil, fmt.Errorf("not expects 1 argument; got %d", len(exprs))
	}
	v := exprs[0].Eval(c)
	asBool, isBool := v.(*BoolValue)
	if !isBool {
		return nil, fmt.Errorf("not expects bool argument, got %v", v)
	}
	return NewBoolValue(!asBool.Get()), nil
}
