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
		Eval(*EvalContext) (Value, error)

		// CodeStr will return the code representation of the given expression.
		CodeStr() string

		// SourcePos returns the location that the expression started in the source.
		SourcePos() ScannerPosition
	}

	// CallExpr is a function call. The first expression is treated as a function,
	// with the remaining elements passed to it.
	CallExpr struct {
		Exprs []Expr
		Pos   ScannerPosition
	}

	// IfExpr is an if expression. Cond is evaluated: if true, Case1 is
	// evaluated and returned; if false Case2 will be.
	IfExpr struct {
		Cond         Expr
		Case1, Case2 Expr
		Pos          ScannerPosition
	}

	// FnExpr is a function definition expression. It has a set of arguments and a
	// body, and will evaluate the body with the given arguments when called.
	FnExpr struct {
		Args []Arg
		Body []Expr
		Pos  ScannerPosition
	}

	// Arg is a single element in a function list.
	Arg struct {
		Ident string
	}

	// LetExpr represents an assignment of a value to an identifier. When
	// evaluated, adds the value to the evaluation context.
	LetExpr struct {
		Ident *IdentLiteral
		Value Expr
		Pos   ScannerPosition
	}
)

// NewCallExpr creates a new CallExpr out of the given sub-expressions. Will
// treat the first argument as the function, and the remaining arguments as the
// arguments.
func NewCallExpr(exprs ...Expr) *CallExpr {
	return &CallExpr{
		Exprs: exprs,
	}
}

// Eval will evaluate the expression and return its results.
func (ce *CallExpr) Eval(ec *EvalContext) (Value, error) {
	if len(ce.Exprs) == 0 {
		return &NilValue{}, nil
	}

	fn, fnErr := evalToFunc(ec, ce.Exprs[0])
	if fnErr != nil {
		return nil, fnErr
	}

	vals := []Value{}
	for _, expr := range ce.Exprs[1:] {
		v, err := expr.Eval(ec)
		if err != nil {
			// todo (bs): augment with trace
			return nil, err
		}
		vals = append(vals, v)
	}
	callVal, callValErr := fn.Fn(ec, vals...)
	return callVal, callValErr
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

// SourcePos is the location in source this expression came from.
func (ce *CallExpr) SourcePos() ScannerPosition {
	return ce.Pos
}

// NewIfExpr builds a new if statement with the given condition and cases. The
// cases may be left nil.
func NewIfExpr(cond Expr, case1, case2 Expr) *IfExpr {
	if case1 == nil {
		case1 = NewNilLiteral()
	}
	if case2 == nil {
		case2 = NewNilLiteral()
	}
	return &IfExpr{
		Cond:  cond,
		Case1: case1,
		Case2: case2,
	}
}

// Eval evaluates the if and returns the evaluated contents of the according
// case.
func (ie *IfExpr) Eval(ec *EvalContext) (Value, error) {
	condV, condVErr := ie.Cond.Eval(ec)
	if condVErr != nil {
		return nil, condVErr
	}
	asBool, isBool := condV.(*BoolValue)
	if !isBool {
		return nil, &TypeError{
			Actual:   fmt.Sprintf("%T", condV),
			Expected: fmt.Sprintf("%T", (*BoolValue)(nil)),
			Pos:      ie.Cond.SourcePos(),
		}
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

// SourcePos is the location in source this expression came from.
func (ie *IfExpr) SourcePos() ScannerPosition {
	return ie.Pos
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
func (fe *FnExpr) Eval(parentEc *EvalContext) (Value, error) {

	// ques (bs): how should stack traces work here? At this point, for full
	// traces (rather than just "origination errors")

	fn := func(_ *EvalContext, vals ...Value) (Value, error) {
		if len(fe.Args) != len(vals) {

			// todo (bs): add pos information.
			return nil, fmt.Errorf("expected %d arguments in call; got %d",
				len(fe.Args), len(vals))
		}

		evalEc := parentEc.SubContext(nil)
		for i, arg := range fe.Args {
			evalEc.Add(arg.Ident, vals[i])
		}

		var evalV Value
		for _, e := range fe.Body {
			v, err := e.Eval(evalEc)
			if err != nil {
				// todo (bs): add pos information
				return nil, err
			}
			evalV = v
		}
		if evalV == nil {
			evalV = &NilValue{}
		}
		return evalV, nil
	}

	return &FuncValue{
		Fn: fn,
	}, nil
}

// CodeStr will return the code representation of the fn expression.
func (fe *FnExpr) CodeStr() string {
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

// SourcePos is the location in source this expression came from.
func (fe *FnExpr) SourcePos() ScannerPosition {
	return fe.Pos
}

// Eval will assign the underlying value to the ident on the context, and return
// the value.
func (le *LetExpr) Eval(ec *EvalContext) (Value, error) {
	identStr := le.Ident.Val
	v, err := le.Value.Eval(ec)
	if err != nil {
		// todo (bs): maybe add pos information
		return nil, err
	}
	ec.Add(identStr, v)
	return v, nil
}

// CodeStr will return the code representation of the let expression.
func (le *LetExpr) CodeStr() string {
	return fmt.Sprintf("(let %s %s)", le.Ident.Val, le.Value.CodeStr())
}

// SourcePos is the location in source this expression came from.
func (le *LetExpr) SourcePos() ScannerPosition {
	return le.Pos
}

// evalToFunc will evaluate the given expression, expecting a function. Will
// return a well-formed error i
func evalToFunc(evalCtx *EvalContext, expr Expr) (*FuncValue, error) {
	var val Value
	switch v := expr.(type) {
	case *IdentLiteral:
		// In the case of idents, manually inspect to see if it's nil. This is to
		// make errors more obvious in the case of a function simply being an
		// undefined name.
		identVal, hasIdent := evalCtx.Resolve(v.Val)
		if !hasIdent {
			return nil, &EvalError{
				Msg: fmt.Sprintf(
					"undefined identifier '%s' cannot be used as function", v.Val),
				Pos: v.SourcePos(),
			}
		}
		val = identVal
	default:
		var v1Err error
		val, v1Err = expr.Eval(evalCtx)
		if v1Err != nil {
			// note (bs): for stack errors; this would still need to be wrapped
			return nil, v1Err
		}
	}
	asFn, isFn := val.(*FuncValue)
	if !isFn {
		return nil, &TypeError{
			Actual:   fmt.Sprintf("%T", val),
			Expected: fmt.Sprintf("%T", (*FuncValue)(nil)),
			Pos:      expr.SourcePos(),
		}
	}
	return asFn, nil
}
