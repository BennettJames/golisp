package golisp2

import (
	"fmt"
	"strings"
)

type (
	// Expr is the fundamental unit of lisp - it represents anything that can be
	// evaluated to a value.
	Expr interface {
		// Eval will evaluate the underlying expression, and return the value (if
		// any) that was calculated and returned.
		Eval(*ExprContext) Value

		// CodeStr will return the code representation of the given expression.
		CodeStr() string

		// InspectStr returns a printable version of the expression.
		InspectStr() string
	}

	// ExprContext is the context on evaluation. It contains a resolvable set of
	// identifiers->values that can be chained.
	ExprContext struct {
		parent *ExprContext
		vals   map[string]Value
	}

	// Value represents any arbitrary value within the lisp interpreting
	// environment. While it just extends "expr", the implicit contract is that no
	// work is actually performed at eval time; it just returns itself.
	//
	// note (bs): not sure this is valuable; also not sure
	Value interface {
		Expr
	}

	// IdentValue is a representation of an identifier in the interpreted
	// environment, whose value is resolved by the context it is evaluated in.
	IdentValue struct {
		// note (bs): I'd like to eventually make it so that identifiers could be
		// "compound lookups"; e.g. "Foo.Bar.A"; in which case I think this should
		// not just be a string. Arguably, that should have it's own datatype
		// anyway.
		Val string
	}

	// NumberValue is a representation of a number value within the interpreted
	// environment.
	NumberValue struct {
		Val float64
	}

	// NilValue is a representation of an null value within the interpreted
	// environment.
	NilValue struct {
	}

	// StringValue is a representation of a string within the interpreted
	// environment.
	StringValue struct {
		Val string
	}

	// BoolValue is a representation of a boolean within the interpreted
	// environment.
	BoolValue struct {
		Val bool
	}

	// FuncValue is a representation of a basic function within the interpreted
	// environment.
	FuncValue struct {
		// Name is the function identifier as it appears in the code.
		Name string

		// Fn is the function body the function value references.
		Fn func(*ExprContext, ...Expr) (Value, error)
	}

	// CellValue is a representation of a pair of values within the interpreted
	// environment. This can be composed to represent lists with standard car/cdr
	// operators.
	CellValue struct {
		Left, Right Value
	}

	// CallExpr is a function call. The first expression is treated as a function,
	// with the remaining elements passed to it.
	CallExpr struct {
		Exprs []Expr
	}

	// IfExpr is an if expression. The condition is evaluated: if true, case1 is
	// evaluated and returned; if false
	IfExpr struct {
		Cond         Expr
		Case1, Case2 Expr
	}

	// FnExpr is a function definition expression. It has a set of arguments and a
	// body, and will evaluate the body with the given arguments when called.
	FnExpr struct {
		Args []Arg
		Body []Expr
	}

	// Arg is a single element in a function list.
	Arg struct {
		Ident string
	}

	// LetExpr represents an assignment of a value to an identifier. When
	// evaluated, adds the value to the evaluation context.
	LetExpr struct {
		Ident *IdentValue
		Value Expr
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
		vals: vals,
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
		Val: ident,
	}
}

// InspectStr will output the name of the identifier.
func (iv *IdentValue) InspectStr() string {
	return fmt.Sprintf("'%s'", iv.Val)
}

// Eval will traverse the context for the identifier and return nil if the value
// is not defined.
//
// todo (bs): consider making failed resolution an error. In this case, it
// should be a "severe error" that bubbles back and most likely halts execution.
// It's *possible* the right way to handle that is by creating a modified value
// interface that can directly support the notion of error.
func (iv *IdentValue) Eval(ec *ExprContext) Value {
	v, ok := ec.Resolve(iv.Val)
	if !ok {
		return NewNilValue()
	}
	return v
}

// CodeStr will return the code representation of the ident value.
func (iv *IdentValue) CodeStr() string {
	return iv.Val
}

// NewNumberValue instantiates a new number with the given value.
func NewNumberValue(v float64) *NumberValue {
	return &NumberValue{
		Val: v,
	}
}

// InspectStr prints the number.
func (nv *NumberValue) InspectStr() string {
	return fmt.Sprintf("%f", nv.Val)
}

// Eval just returns itself.
func (nv *NumberValue) Eval(*ExprContext) Value {
	return nv
}

// CodeStr will return the code representation of the number value.
func (nv *NumberValue) CodeStr() string {
	// todo (bs): this isn't wrong, exactly, but consider printing integers as
	// integers. Of course, that starts getting into the deeper issue of how just
	// having floats is too primitive and there really need to be integers.
	return fmt.Sprintf("%f", nv.Val)
}

// NewNilValue creates a new nil value.
//
// todo (bs): this should return a singleton; no need for duplicates given that
// it's unmodifiable.
func NewNilValue() *NilValue {
	return &NilValue{}
}

// InspectStr outputs "nil".
func (nv *NilValue) InspectStr() string {
	return "nil"
}

// Eval returns the nil value.
func (nv *NilValue) Eval(*ExprContext) Value {
	// note (bs): not sure about this. In general, I feel like eval needs to be
	// more intelligent
	return nv
}

// CodeStr will return the code representation of the nil value.
func (nv *NilValue) CodeStr() string {
	return fmt.Sprintf("nil")
}

// NewStringValue creates a new string value from the given string.
func NewStringValue(str string) *StringValue {
	return &StringValue{
		Val: str,
	}
}

// InspectStr prints the string.
func (sv *StringValue) InspectStr() string {
	return fmt.Sprintf("\"%s\"", sv.Val)
}

// Eval returns the string value.
func (sv *StringValue) Eval(*ExprContext) Value {
	return sv
}

// CodeStr will return the code representation of the string value.
func (sv *StringValue) CodeStr() string {
	// note (bs): this doesn't matter now as it's not supported, but just note
	// that this doesn't work with multiline strings
	return fmt.Sprintf("\"%s\"", sv.Val)
}

// NewBoolValue creates a bool with the given value.
//
// todo (bs): this probably should return singletons for true/false
func NewBoolValue(v bool) *BoolValue {
	return &BoolValue{
		Val: v,
	}
}

// InspectStr prints "true"/"false" based on the value.
func (bv *BoolValue) InspectStr() string {
	return fmt.Sprintf("%t", bv.Val)
}

// Eval returns the bool value.
func (bv *BoolValue) Eval(*ExprContext) Value {
	return bv
}

// CodeStr will return the code representation of the boolean value.
func (bv *BoolValue) CodeStr() string {
	if bv.Val {
		return "true"
	}
	return "false"
}

// NewFuncValue creates a function with the given value.
func NewFuncValue(
	name string,
	fn func(*ExprContext, ...Expr) (Value, error),
) *FuncValue {
	return &FuncValue{
		Fn: fn,
	}
}

// InspectStr outputs some information about the function.
func (fv *FuncValue) InspectStr() string {
	// note (bs): probably want to customize this to print some details about the
	// function itself. That will involve (optionally) retaining the declaration
	// name of the function.
	return fmt.Sprintf("<func>")
}

// Eval evaluates the function using the provided context.
func (fv *FuncValue) Eval(ec *ExprContext) Value {
	return fv
}

// Exec executes the underlying function with the given context and arguments.
func (fv *FuncValue) Exec(ec *ExprContext, exprs ...Expr) (Value, error) {
	return fv.Fn(ec, exprs...)
}

// CodeStr will return the code representation of the function value.
func (fv *FuncValue) CodeStr() string {
	return fv.Name
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
		Left:  left,
		Right: right,
	}
}

// Eval returns the cell.
func (cv *CellValue) Eval(*ExprContext) Value {
	return cv
}

// InspectStr outputs the contents of all the cells.
func (cv *CellValue) InspectStr() string {
	// todo (bs): if second cell is a node, treat this different
	return fmt.Sprintf("(%s . %s)", cv.Left.InspectStr(), cv.Right.InspectStr())
}

// CodeStr will return the code representation of the cell value.
func (cv *CellValue) CodeStr() string {
	return fmt.Sprintf("(cons %s %s)\n", cv.Left.CodeStr(), cv.Right.CodeStr())
}

// NewCallExpr creates a new CallExpr out of the given sub-expressions. Will
// treat the first argument as the function, and the remaining arguments as the
// arguments.
func NewCallExpr(exprs ...Expr) *CallExpr {
	return &CallExpr{
		Exprs: exprs,
	}
}

// Eval will evaluate the expression and return its results.
func (ce *CallExpr) Eval(ec *ExprContext) Value {
	if len(ce.Exprs) == 0 {
		return NewNilValue()
	}

	v1 := ce.Exprs[0].Eval(ec)
	asFn, isFn := v1.(*FuncValue)
	if !isFn {
		if len(ce.Exprs) == 1 {
			return v1
		}
		// fixme (bs): this needs to return an error. Again, either eval needs to be
		// modified to explicitly return (value, error), or Value's
		// signature/behavior needs to be modified to support the notion of an error.
		return nil
	}

	value, err := asFn.Exec(ec, ce.Exprs[1:]...)
	var _ = err // fixme (bs): again, need to handle error passback
	return value
}

// InspectStr returns a user-readable representation of the call expression.
func (ce *CallExpr) InspectStr() string {
	if len(ce.Exprs) == 0 {
		return "<call nil>"
	}
	return fmt.Sprintf("<call \"%s\">", ce.Exprs[0].InspectStr())
}

// CodeStr will return the code representation of the call expression.
func (ce *CallExpr) CodeStr() string {
	var sb strings.Builder
	sb.WriteString("(")
	for i, e := range ce.Exprs {
		if i > 0 {
			sb.WriteString(" ")
		}
		sb.WriteString(e.CodeStr())
	}
	sb.WriteString(")\n")
	return sb.String()
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
		Cond:  cond,
		Case1: case1,
		Case2: case2,
	}
}

// Eval evaluates the if and returns the evaluated contents of the according
// case.
func (ie *IfExpr) Eval(ec *ExprContext) Value {
	condV := ie.Cond.Eval(ec)
	asBool, isBool := condV.(*BoolValue)
	if !isBool {
		// fixme (bs): this should return an error
		return NewNilValue()
	}
	if asBool.Val {
		return ie.Case1.Eval(ec)
	}
	return ie.Case2.Eval(ec)
}

// CodeStr will return the code representation of the if expression.
func (ie *IfExpr) CodeStr() string {
	var sb strings.Builder
	sb.WriteString("(if ")
	sb.WriteString(ie.Cond.CodeStr())
	sb.WriteString("\n")
	sb.WriteString(ie.Case1.CodeStr())
	sb.WriteString("\n")
	sb.WriteString(ie.Case2.CodeStr())
	sb.WriteString(")\n")
	return sb.String()
}

// InspectStr returns a user-readable representation of the if expression.
func (ie *IfExpr) InspectStr() string {
	return fmt.Sprintf("(todo)")
}

// NewFnExpr builds a new function expression with the given arguments and body.
func NewFnExpr(args []Arg, body []Expr) *FnExpr {
	return &FnExpr{
		Args: args,
		Body: body,
	}
}

// Eval returns an evaluate-able function value. Note that this does *not*
// execute the function; it must be evaluated within a call to be actually
// executed.
func (fe *FnExpr) Eval(parentEc *ExprContext) Value {
	// fixme (bs): I don't think this should be returning a func value per se.
	// This is a good case where perhaps having some plain functions in place of
	// the strict AST would make sense; but I'm not sure yet.

	return NewFuncValue("", func(
		callEc *ExprContext,
		callExprs ...Expr,
	) (Value, error) {

		if len(fe.Args) != len(callExprs) {
			return nil, fmt.Errorf("expected %d arguments in call; got %d",
				len(fe.Args), len(callExprs))
		}
		evalEc := parentEc.SubContext(nil)
		for i, arg := range fe.Args {
			evalEc.Add(arg.Ident, callExprs[i].Eval(callEc))
		}

		var evalV Value
		for _, e := range fe.Body {
			evalV = e.Eval(evalEc)
		}
		if evalV == nil {
			evalV = NewNilValue()
		}
		return evalV, nil
	})
}

// CodeStr will return the code representation of the fn expression.
func (fe *FnExpr) CodeStr() string {
	// fixme (bs): implement

	var sb strings.Builder
	sb.WriteString("(fn (")
	for i, a := range fe.Args {
		if i > 0 {
			sb.WriteString(" ")
		}
		sb.WriteString(a.Ident)
	}
	sb.WriteString(")\n")

	for _, e := range fe.Body {
		sb.WriteString(e.CodeStr())
	}
	sb.WriteString(")\n")
	return sb.String()
}

// InspectStr returns a user-readable representation of the function expression.
func (fe *FnExpr) InspectStr() string {
	// ques (bs): what should this be? I'd say the name and the arg list. The
	// names not strictly known here though; so maybe just the arg list?
	var sb strings.Builder
	sb.WriteString("fn (")
	for i, a := range fe.Args {
		if i > 0 {
			sb.WriteString(" ")
		}
		sb.WriteString(a.Ident)
	}
	sb.WriteString(")")
	return sb.String()
}

// Eval will assign the underlying value to the ident on the context, and return
// the value.
func (le *LetExpr) Eval(ec *ExprContext) Value {
	identStr := le.Ident.Val
	v := le.Value.Eval(ec)
	ec.Add(identStr, v)
	return v
}

// CodeStr will return the code representation of the let expression.
func (le *LetExpr) CodeStr() string {
	return fmt.Sprintf("(let %s %s)", le.Ident.Val, le.Value.CodeStr())
}

// InspectStr returns a user-readable representation of the let expression.
func (le *LetExpr) InspectStr() string {
	return fmt.Sprintf("<assign \"%s\">", le.Ident.Val)
}
