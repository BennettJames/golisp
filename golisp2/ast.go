package golisp2

import (
	"fmt"
)

type (
	// Expr is the fundamental unit of lisp - it represents anything that can be
	// evaluated to a value.
	Expr interface {
		Eval(*ExprContext) Value
	}

	// ExprContext is the context on evaluation. It contains a resolvable set of
	// identifiers->values that can be chained.
	ExprContext struct {
		parent *ExprContext
		vals   map[string]Value
	}

	// Value represents any arbitrary value within the lisp interpreting
	// environment.
	Value interface {
		Expr

		// PrintStr returns a printable version of the value. Note that this is not
		// the same as casting to a string.
		PrintStr() string
	}

	// IdentValue is a representation of an identifier in the interpreted
	// environment, whose value is resolved by the context it is evaluated in.
	IdentValue struct {
		// note (bs): I'd like to eventually make it so that identifiers could be
		// "compound lookups"; e.g. "Foo.Bar.A"; in which case I think this should
		// not just be a string. Arguably, that should have it's own datatype
		// anyway.
		ident string
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

		fn func(*ExprContext, ...Expr) (Value, error)
	}

	// CellValue is a representation of a pair of values within the interpreted
	// environment. This can be composed to represent lists with standard car/cdr
	// operators.
	CellValue struct {
		left, right Value
	}

	// CallExpr is a function call. The first expression is treated as a function,
	// with the remaining elements passed to it.
	CallExpr struct {
		exprs []Expr
	}

	// IfExpr is an if expression. The condition is evaluated: if true, case1 is
	// evaluated and returned; if false
	IfExpr struct {
		cond         Expr
		case1, case2 Expr
	}

	// FnExpr is a function definition expression. It has a set of arguments and a
	// body, and will evaluate the body with the given arguments when called.
	FnExpr struct {
		args []Arg
		body []Expr
	}

	// Arg is a single element in a function list.
	Arg struct {
		Ident string
	}
)

// NewContext returns a new context with no parent. initialVals contains any
// values that the context should be initialized with; it can be left nil.
func NewContext(initialVals map[string]Value) *ExprContext {
	vals := map[string]Value{}
	for k, v := range initialVals {
		vals[k] = v
	}
	return &ExprContext{
		vals: map[string]Value{},
	}
}

// SubContext creates a new context with the current context as it's parent.
func (ec *ExprContext) SubContext(initialVals map[string]Value) *ExprContext {
	sub := NewContext(initialVals)
	sub.parent = ec
	return sub
}

// Add extends the current context with the provided value.
func (ec *ExprContext) Add(ident string, val Value) {
	ec.vals[ident] = val
}

// Resolve traverses the expr for the given ident. Will return it if found;
// otherwise a nil value and "false".
func (ec *ExprContext) Resolve(ident string) (Value, bool) {
	if ec == nil {
		return NewNilValue(), false
	}
	if v, ok := ec.vals[ident]; ok {
		return v, true
	}
	return ec.parent.Resolve(ident)
}

// NewIdentValue instantiates a new identifier value with the given identifier
// token.
func NewIdentValue(ident string) *IdentValue {
	return &IdentValue{
		ident: ident,
	}
}

// PrintStr will output the name of the identifier.
func (iv *IdentValue) PrintStr() string {
	return fmt.Sprintf("'%s'", iv.ident)
}

// Eval will traverse the context for the identifier and return nil if the value
// is not defined.
//
// todo (bs): consider making failed resolution an error. In this case, it
// should be a "severe error" that bubbles back and most likely halts execution.
// It's *possible* the right way to handle that is by creating a modified value
// interface that can directly support the notion of error.
func (iv *IdentValue) Eval(ec *ExprContext) Value {
	v, ok := ec.Resolve(iv.ident)
	if !ok {
		return NewNilValue()
	}
	return v
}

// Get just returns the underlying ident string.
func (iv *IdentValue) Get() string {
	return iv.ident
}

// NewNumberValue instantiates a new number with the given value.
func NewNumberValue(v float64) *NumberValue {
	return &NumberValue{
		val: v,
	}
}

// PrintStr prints the number.
func (nv *NumberValue) PrintStr() string {
	return fmt.Sprintf("%f", nv.val)
}

// Eval just returns itself.
func (nv *NumberValue) Eval(*ExprContext) Value {
	return nv
}

// Get just returns the underlying number.
func (nv *NumberValue) Get() float64 {
	return nv.val
}

// NewNilValue creates a new nil value.
//
// todo (bs): this should return a singleton; no need for duplcates given that
// it's unmodifiable.
func NewNilValue() *NilValue {
	return &NilValue{}
}

// PrintStr outputs "nil".
func (nv *NilValue) PrintStr() string {
	return "nil"
}

// Eval returns the nil value.
func (nv *NilValue) Eval(*ExprContext) Value {
	// note (bs): not sure about this. In general, I feel like eval needs to be
	// more intelligent
	return nv
}

// NewStringValue creates a new string value from the given string.
func NewStringValue(str string) *StringValue {
	return &StringValue{
		val: str,
	}
}

// PrintStr prints the string.
func (sv *StringValue) PrintStr() string {
	return fmt.Sprintf("\"%s\"", sv.val)
}

// Eval returns the string value.
func (sv *StringValue) Eval(*ExprContext) Value {
	return sv
}

// Get returns the raw string value.
func (sv *StringValue) Get() string {
	return sv.val
}

// NewBoolValue creates a bool with the given value.
//
// todo (bs): this probably should return singletons for true/false
func NewBoolValue(v bool) *BoolValue {
	return &BoolValue{
		val: v,
	}
}

// PrintStr prints "true"/"false" based on the value.
func (bv *BoolValue) PrintStr() string {
	return fmt.Sprintf("%t", bv.val)
}

// Eval returns the bool value.
func (bv *BoolValue) Eval(*ExprContext) Value {
	return bv
}

// Get returns the raw bool value.
func (bv *BoolValue) Get() bool {
	return bv.val
}

// NewFuncValue creates a function with the given value.
func NewFuncValue(fn func(*ExprContext, ...Expr) (Value, error)) *FuncValue {
	return &FuncValue{
		fn: fn,
	}
}

// PrintStr outputs some information about the function.
func (fv *FuncValue) PrintStr() string {
	// note (bs): probably want to customize this to print some details about the
	// function itself. That will involve (optionally) retaining the declaration
	// name of the function.
	return fmt.Sprintf("<func>")
}

// Eval evaluates the function using the provided context.
func (fv *FuncValue) Eval(ec *ExprContext) Value {
	return fv
}

// Get returns the function value.
func (fv *FuncValue) Get() func(*ExprContext, ...Expr) (Value, error) {
	return fv.fn
}

// Exec executes the underlying function with the given context and arguments.
func (fv *FuncValue) Exec(ec *ExprContext, exprs ...Expr) (Value, error) {
	return fv.fn(ec, exprs...)
}

// NewCellValue creates a cell with the given left/right values. Either can be
// 'nil'.
func NewCellValue(left, right Value) *CellValue {
	if left == nil {
		left = NewNilValue()
	}
	if right == nil {
		right = NewNilValue()
	}
	return &CellValue{
		left:  left,
		right: right,
	}
}

// PrintStr outputs the contents of all the cells.
func (nv *CellValue) PrintStr() string {
	// todo (bs): if second cell is a node, treat this different
	return fmt.Sprintf("(%s . %s)", nv.left.PrintStr(), nv.right.PrintStr())
}

// Eval returns the cell.
func (nv *CellValue) Eval(*ExprContext) Value {
	return nv
}

// Get returns the cell values.
func (nv *CellValue) Get() (left, right Value) {
	return nv.left, nv.right
}

// NewCallExpr creates a new CallExpr out of the given sub-expressions. Will
// treat the first argument as the function, and the remaining arguments as the
// arguments.
func NewCallExpr(exprs ...Expr) *CallExpr {
	return &CallExpr{
		exprs: exprs,
	}
}

// Eval will evaluate the expression and return its results.
func (ce *CallExpr) Eval(ec *ExprContext) Value {
	if len(ce.exprs) == 0 {
		return NewNilValue()
	}

	v1 := ce.exprs[0].Eval(ec)
	asFn, isFn := v1.(*FuncValue)
	if !isFn {
		if len(ce.exprs) == 1 {
			return v1
		}
		// fixme (bs): this needs to return an error. Again, either eval needs to be
		// modified to explicitly return (value, error), or Value's
		// signature/behavior needs to be modified to support the notion of an error.
		return nil
	}

	value, err := asFn.Exec(ec, ce.exprs[1:]...)
	var _ = err // fixme (bs): again, need to handle error passback
	return value
}

// Get returns the underlying set of expressions in the call.
func (ce *CallExpr) Get() []Expr {
	return ce.exprs
}

// NewIfExpr builds a new if statement with the given condition and cases. The
// cases may be left nil.
func NewIfExpr(cond Expr, case1, case2 Expr) *IfExpr {
	if case1 == nil {
		case1 = NewNilValue()
	}
	if case2 == nil {
		case2 = NewNilValue()
	}
	return &IfExpr{
		cond:  cond,
		case1: case1,
		case2: case2,
	}
}

// Eval evaluates the if and returns the evaluated contents of the according
// case.
func (ie *IfExpr) Eval(ec *ExprContext) Value {
	condV := ie.cond.Eval(ec)
	asBool, isBool := condV.(*BoolValue)
	if !isBool {
		// fixme (bs): this should return an error
		return NewNilValue()
	}
	if asBool.Get() {
		return ie.case1.Eval(ec)
	}
	return ie.case2.Eval(ec)
}

// NewFnExpr builds a new function expression with the given arguments and body.
func NewFnExpr(args []Arg, body []Expr) *FnExpr {
	return &FnExpr{
		args: args,
		body: body,
	}
}

// Eval returns an evaluate-able function value. Note that this does *not*
// execute the function; it must be evaluated within a call to be actually
// executed.
func (fe *FnExpr) Eval(parentEc *ExprContext) Value {
	return NewFuncValue(func(
		callEc *ExprContext,
		callExprs ...Expr,
	) (Value, error) {

		if len(fe.args) != len(callExprs) {
			return nil, fmt.Errorf("expected %d arguments in call; got %d",
				len(fe.args), len(callExprs))
		}
		evalEc := parentEc.SubContext(nil)
		for i, arg := range fe.args {
			evalEc.Add(arg.Ident, callExprs[i].Eval(callEc))
		}

		var evalV Value
		for _, e := range fe.body {
			evalV = e.Eval(evalEc)
		}
		if evalV == nil {
			evalV = NewNilValue()
		}
		return evalV, nil
	})
}
