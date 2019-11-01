package golisp2

import (
	"fmt"
	"strings"
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

	CallExpr struct {
		exprs []Expr
	}

	IfExpr struct {
		cond         Expr
		case1, case2 Expr
	}

	FnExpr struct {
		args []Arg
		body []Expr
	}

	Arg struct {
		ident string
	}
)

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
func (nv *IdentValue) PrintStr() string {
	return fmt.Sprintf("'%s'", nv.ident)
}

// Eval will traverse the context for the identifier and return nil if the value
// is not defined.
//
// todo (bs): consider making failed resolution an error. In this case, it
// should be a "severe error" that bubbles back and most likely halts execution.
// It's *possible* the right way to handle that is by creating a modified value
// interface that can directly support the notion of error.
func (nv *IdentValue) Eval(ec *ExprContext) Value {
	v, ok := ec.Resolve(nv.ident)
	if !ok {
		return NewNilValue()
	}
	return v
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
func (nv *StringValue) PrintStr() string {
	return fmt.Sprintf("\"%s\"", nv.val)
}

// Eval returns the string value.
func (nv *StringValue) Eval(*ExprContext) Value {
	return nv
}

// Get returns the raw string value.
func (nv *StringValue) Get() string {
	return nv.val
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
func (nv *BoolValue) PrintStr() string {
	return fmt.Sprintf("%t", nv.val)
}

// Eval returns the bool value.
func (nv *BoolValue) Eval(*ExprContext) Value {
	return nv
}

// Get returns the raw bool value.
func (nv *BoolValue) Get() bool {
	return nv.val
}

// NewFuncValue creates a function with the given value.
func NewFuncValue(fn func(*ExprContext, ...Expr) (Value, error)) *FuncValue {
	return &FuncValue{
		fn: fn,
	}
}

// PrintStr outputs some information about the function.
func (nv *FuncValue) PrintStr() string {
	// note (bs): probably want to customize this to print some details about the
	// function itself
	return fmt.Sprintf("<func>")
}

// Eval evaluates the function using the provided context.
func (nv *FuncValue) Eval(ec *ExprContext) Value {
	return nv
}

// Get returns the function value.
func (nv *FuncValue) Get() func(*ExprContext, ...Expr) (Value, error) {
	return nv.fn
}

func (nv *FuncValue) Exec(ec *ExprContext, exprs ...Expr) (Value, error) {
	return nv.fn(ec, exprs...)
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

func NewCallExpr(exprs ...Expr) *CallExpr {
	return &CallExpr{
		exprs: exprs,
	}
}

func (pe *CallExpr) Eval(ec *ExprContext) Value {
	if len(pe.exprs) == 0 {
		return NewNilValue()
	}

	v1 := pe.exprs[0].Eval(ec)
	asFn, isFn := v1.(*FuncValue)
	if !isFn {
		if len(pe.exprs) == 1 {
			return v1
		}
		// fixme (bs): this needs to return an error. Again, either eval needs to be
		// modified to explicitly return (value, error), or Value's
		// signature/behavior needs to be modified to support the notion of an error.
		return nil
	}

	value, err := asFn.Exec(ec, pe.exprs[1:]...)
	var _ = err // fixme (bs): again, need to handle error passback
	return value
}

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

func NewFnExpr(args []Arg, body []Expr) *FnExpr {
	return &FnExpr{
		args: args,
		body: body,
	}
}

func (fe *FnExpr) Eval(parentEc *ExprContext) Value {
	return NewFuncValue(func(
		callEc *ExprContext,
		callExprs ...Expr,
	) (Value, error) {
		evalEc := &ExprContext{
			parent: parentEc,
			vals:   map[string]Value{},
		}

		// note (bs): this is very strict; will eventually likely want the ability
		// to specify things like optional args, varargs, and defaults. For now, the
		// arity of a user-defined function is always exact.
		if len(fe.args) != len(callExprs) {
			return nil, fmt.Errorf("expected %d arguments in call; got %d",
				len(fe.args), len(callExprs))
		}
		for i, arg := range fe.args {
			v := callExprs[i].Eval(callEc)
			evalEc.vals[arg.ident] = v
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

// note (bs): I'm o.k. with this for the immediate future, but functions should
// likely exist in a separate file. Let's worry about that later: I think
// packages should be broken up and all this should be moved out.

func addFn(c *ExprContext, exprs ...Expr) (Value, error) {
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

func subFn(c *ExprContext, exprs ...Expr) (Value, error) {
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

func multFn(c *ExprContext, exprs ...Expr) (Value, error) {
	total := float64(1)
	for _, e := range exprs {
		v := e.Eval(c)
		asNum, isNum := v.(*NumberValue)
		if !isNum {
			return nil, fmt.Errorf("non-number value in add: %v", asNum.PrintStr())
		}
		total *= asNum.Get()
	}
	return &NumberValue{
		val: total,
	}, nil
}

func divFn(c *ExprContext, exprs ...Expr) (Value, error) {
	total := float64(1)
	for i, e := range exprs {
		v := e.Eval(c)
		asNum, isNum := v.(*NumberValue)
		if !isNum {
			return nil, fmt.Errorf("non-number value in add: %v", asNum.PrintStr())
		}
		if i == 0 {
			total = asNum.Get()
		} else {
			total /= asNum.Get()
		}
	}
	return &NumberValue{
		val: total,
	}, nil
}

func concatFn(c *ExprContext, exprs ...Expr) (Value, error) {
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

func consFn(c *ExprContext, exprs ...Expr) (Value, error) {
	if len(exprs) > 2 {
		return nil, fmt.Errorf("cons expects 0-2 argument; got %d", len(exprs))
	}
	v1 := exprs[0].Eval(c)
	v2 := exprs[1].Eval(c)
	return NewCellValue(v1, v2), nil
}

func carFn(c *ExprContext, exprs ...Expr) (Value, error) {
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

func cdrFn(c *ExprContext, exprs ...Expr) (Value, error) {
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

func andFn(c *ExprContext, exprs ...Expr) (Value, error) {
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

func orFn(c *ExprContext, exprs ...Expr) (Value, error) {
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

func notFn(c *ExprContext, exprs ...Expr) (Value, error) {
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

func eqNumFn(ec *ExprContext, exprs ...Expr) (Value, error) {
	if len(exprs) != 2 {
		return nil, fmt.Errorf("eq expects 2 arguments; got %d", len(exprs))
	}
	v1 := exprs[0].Eval(ec)
	v2 := exprs[1].Eval(ec)
	v1AsNum, v1IsNum := v1.(*NumberValue)
	v2AsNum, v2IsNum := v2.(*NumberValue)
	if !v1IsNum {
		return nil, fmt.Errorf("eq expects number arguments")
	}
	if !v2IsNum {
		return nil, fmt.Errorf("eq expects number arguments")
	}
	return NewBoolValue(v1AsNum.Get() == v2AsNum.Get()), nil
}

func gtNumFn(ec *ExprContext, exprs ...Expr) (Value, error) {
	if len(exprs) != 2 {
		return nil, fmt.Errorf("gt expects 2 arguments; got %d", len(exprs))
	}
	v1 := exprs[0].Eval(ec)
	v2 := exprs[1].Eval(ec)
	v1AsNum, v1IsNum := v1.(*NumberValue)
	v2AsNum, v2IsNum := v2.(*NumberValue)
	if !v1IsNum {
		return nil, fmt.Errorf("gt expects number arguments")
	}
	if !v2IsNum {
		return nil, fmt.Errorf("gt expects number arguments")
	}
	return NewBoolValue(v1AsNum.Get() > v2AsNum.Get()), nil
}

func ltNumFn(ec *ExprContext, exprs ...Expr) (Value, error) {
	if len(exprs) != 2 {
		return nil, fmt.Errorf("lt expects 2 arguments; got %d", len(exprs))
	}
	v1 := exprs[0].Eval(ec)
	v2 := exprs[1].Eval(ec)
	v1AsNum, v1IsNum := v1.(*NumberValue)
	v2AsNum, v2IsNum := v2.(*NumberValue)
	if !v1IsNum {
		return nil, fmt.Errorf("lt expects number arguments")
	}
	if !v2IsNum {
		return nil, fmt.Errorf("lt expects number arguments")
	}
	return NewBoolValue(v1AsNum.Get() < v2AsNum.Get()), nil
}

func gteNumFn(ec *ExprContext, exprs ...Expr) (Value, error) {
	if len(exprs) != 2 {
		return nil, fmt.Errorf("gte expects 2 arguments; got %d", len(exprs))
	}
	v1 := exprs[0].Eval(ec)
	v2 := exprs[1].Eval(ec)
	v1AsNum, v1IsNum := v1.(*NumberValue)
	v2AsNum, v2IsNum := v2.(*NumberValue)
	if !v1IsNum {
		return nil, fmt.Errorf("gte expects number arguments")
	}
	if !v2IsNum {
		return nil, fmt.Errorf("gte expects number arguments")
	}
	return NewBoolValue(v1AsNum.Get() >= v2AsNum.Get()), nil
}

func lteNumFn(ec *ExprContext, exprs ...Expr) (Value, error) {
	if len(exprs) != 2 {
		return nil, fmt.Errorf("lte expects 2 arguments; got %d", len(exprs))
	}
	v1 := exprs[0].Eval(ec)
	v2 := exprs[1].Eval(ec)
	v1AsNum, v1IsNum := v1.(*NumberValue)
	v2AsNum, v2IsNum := v2.(*NumberValue)
	if !v1IsNum {
		return nil, fmt.Errorf("lte expects number arguments")
	}
	if !v2IsNum {
		return nil, fmt.Errorf("lte expects number arguments")
	}
	return NewBoolValue(v1AsNum.Get() <= v2AsNum.Get()), nil
}
