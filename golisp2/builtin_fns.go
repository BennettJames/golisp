package golisp2

import (
	"fmt"
	"strings"
)

// BuiltinContext returns a context that contains the full set of builtin
// functions. Note this just includes built-in plain functions; not operators.
func BuiltinContext() *ExprContext {
	return NewContext(map[string]Value{
		"concat": NewFuncValue("concat", concatFn),
		"cons":   NewFuncValue("cons", consFn),
		"car":    NewFuncValue("car", carFn),
		"cdr":    NewFuncValue("cdr", cdrFn),
		"and":    NewFuncValue("and", andFn),
		"or":     NewFuncValue("or", orFn),
		"not":    NewFuncValue("not", notFn),
	})
}

//
// Explicit, named built-ins
//

func concatFn(c *ExprContext, exprs ...Expr) (Value, error) {
	var sb strings.Builder
	for _, e := range exprs {
		v := e.Eval(c)
		asStr, isStr := v.(*StringValue)
		if !isStr {
			return nil, fmt.Errorf("non-number value in add: %v", v.InspectStr())
		}
		sb.WriteString(asStr.Get())
	}
	return &StringValue{
		Val: sb.String(),
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
	return asNode.Left, nil
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
	return asNode.Right, nil
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
		// todo (bs): strongly consider short-circuiting this if false is returned
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
		// todo (bs): strongly consider short-circuiting this if true is returned
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

//
// Mathematical operator built-ins
//

func addFn(c *ExprContext, exprs ...Expr) (Value, error) {
	total := float64(0)
	for _, e := range exprs {
		v := e.Eval(c)
		asNum, isNum := v.(*NumberValue)
		if !isNum {
			// note (bs): eventually, try to make a version of this error that's more
			// portable, obvious, and a little more resilient to nil values.
			return nil, fmt.Errorf("non-number value in add: %v", asNum.InspectStr())
		}
		total += asNum.Get()
	}
	return &NumberValue{
		Val: total,
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
			return nil, fmt.Errorf("non-number value in add: %v", v.InspectStr())
		}
		if i == 0 {
			total = asNum.Get()
		} else {
			total -= asNum.Get()
		}
	}

	return &NumberValue{
		Val: total,
	}, nil
}

func multFn(c *ExprContext, exprs ...Expr) (Value, error) {
	total := float64(1)
	for _, e := range exprs {
		v := e.Eval(c)
		asNum, isNum := v.(*NumberValue)
		if !isNum {
			return nil, fmt.Errorf("non-number value in add: %v", asNum.InspectStr())
		}
		total *= asNum.Get()
	}
	return &NumberValue{
		Val: total,
	}, nil
}

func divFn(c *ExprContext, exprs ...Expr) (Value, error) {
	total := float64(1)
	for i, e := range exprs {
		v := e.Eval(c)
		asNum, isNum := v.(*NumberValue)
		if !isNum {
			return nil, fmt.Errorf("non-number value in add: %v", asNum.InspectStr())
		}
		if i == 0 {
			total = asNum.Get()
		} else {
			total /= asNum.Get()
		}
	}
	return &NumberValue{
		Val: total,
	}, nil
}

//
// Comparison operator built-in
//

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
