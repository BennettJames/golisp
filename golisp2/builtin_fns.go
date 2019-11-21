package golisp2

import (
	"fmt"
	"strings"
)

// BuiltinContext returns a context that contains the full set of builtin
// functions. Note this just includes built-in plain functions; not operators.
func BuiltinContext() *EvalContext {
	return NewContext(map[string]Value{
		"concat": NewFuncLiteral("concat", concatFn),
		"cons":   NewFuncLiteral("cons", consFn),
		"car":    NewFuncLiteral("car", carFn),
		"cdr":    NewFuncLiteral("cdr", cdrFn),
		"and":    NewFuncLiteral("and", andFn),
		"or":     NewFuncLiteral("or", orFn),
		"not":    NewFuncLiteral("not", notFn),
	})
}

//
// Explicit, named built-ins
//

func concatFn(c *EvalContext, vals ...Value) (Value, error) {
	var sb strings.Builder
	for _, v := range vals {
		asStr, isStr := v.(*StringValue)
		if !isStr {
			return nil, fmt.Errorf("non-string value in add: %v", v.InspectStr())
		}
		sb.WriteString(asStr.Val)
	}
	return &StringValue{
		Val: sb.String(),
	}, nil
}

func consFn(c *EvalContext, vals ...Value) (Value, error) {
	if len(vals) > 2 {
		return nil, fmt.Errorf("cons expects 0-2 argument; got %d", len(vals))
	}
	var v1, v2 Value
	if len(vals) > 0 {
		v1 = vals[0]
	}
	if len(vals) > 1 {
		v2 = vals[1]
	}
	return NewCellValue(v1, v2), nil
}

func carFn(c *EvalContext, vals ...Value) (Value, error) {
	if len(vals) != 1 {
		return nil, fmt.Errorf("car expects 1 argument; got %d", len(vals))
	}
	asNode, isNode := vals[0].(*CellValue)
	if !isNode {
		return nil, fmt.Errorf("car expects a cell type, got %v", asNode)
	}
	return asNode.Left, nil
}

func cdrFn(c *EvalContext, vals ...Value) (Value, error) {
	if len(vals) != 1 {
		return nil, fmt.Errorf("cdr expects 1 argument; got %d", len(vals))
	}
	asNode, isNode := vals[0].(*CellValue)
	if !isNode {
		return nil, fmt.Errorf("cdr expects a cell type, got %v", asNode)
	}
	return asNode.Right, nil
}

func andFn(c *EvalContext, vals ...Value) (Value, error) {
	if len(vals) == 0 {
		return nil, fmt.Errorf("and expects at least 1 argument; got %d", len(vals))
	}
	total := true
	for _, v := range vals {
		asBool, isBool := v.(*BoolValue)
		if !isBool {
			return nil, fmt.Errorf("and expects bool types, got %v", v)
		}
		// todo (bs): strongly consider short-circuiting this if false is returned
		total = total && asBool.Val
	}
	return &BoolValue{
		Val: total,
	}, nil
}

func orFn(c *EvalContext, vals ...Value) (Value, error) {
	if len(vals) == 0 {
		return nil, fmt.Errorf("or expects at least 1 argument; got %d", len(vals))
	}
	total := false
	for _, v := range vals {
		asBool, isBool := v.(*BoolValue)
		if !isBool {
			return nil, fmt.Errorf("or expects bool types, got %v", v)
		}
		// todo (bs): strongly consider short-circuiting this if true is returned
		total = total || asBool.Val
	}
	return &BoolValue{
		Val: total,
	}, nil
}

func notFn(c *EvalContext, vals ...Value) (Value, error) {
	if len(vals) != 1 {
		return nil, fmt.Errorf("not expects 1 argument; got %d", len(vals))
	}
	asBool, isBool := vals[0].(*BoolValue)
	if !isBool {
		return nil, fmt.Errorf("not expects bool argument, got %v", vals[0])
	}
	return &BoolValue{
		Val: !asBool.Val,
	}, nil
}

//
// Mathematical operator built-ins
//

func addFn(c *EvalContext, vals ...Value) (Value, error) {
	total := float64(0)
	for _, v := range vals {
		asNum, isNum := v.(*NumberValue)
		if !isNum {
			// note (bs): eventually, try to make a version of this error that's more
			// portable, obvious, and a little more resilient to nil values.
			return nil, fmt.Errorf("non-number value in add: %v", v.InspectStr())
		}
		total += asNum.Val
	}
	return &NumberValue{
		Val: total,
	}, nil
}

func subFn(c *EvalContext, vals ...Value) (Value, error) {
	// ques (bs): should I still enforce minimum airity requirements here? I'm
	// sorta inclined to say yes; but not sure how much I care about that right
	// now. Particularly: that seems to get into deeper questions of type
	// enforcement. Something like this could just be reduced to a set of number
	// values, and an error returned if
	//
	// That all seems like a "later" task. I'd like to just grind a bit on the
	// core language; some better limitations or even

	total := float64(0)
	for i, v := range vals {
		asNum, isNum := v.(*NumberValue)
		if !isNum {
			// note (bs): eventually, try to make a version of this error that's more
			// portable, obvious, and a little more resilient to nil values.
			return nil, fmt.Errorf("non-number value in add: %v", v.InspectStr())
		}
		if i == 0 {
			total = asNum.Val
		} else {
			total -= asNum.Val
		}
	}

	return &NumberValue{
		Val: total,
	}, nil
}

func multFn(c *EvalContext, vals ...Value) (Value, error) {
	total := float64(1)
	for _, v := range vals {
		asNum, isNum := v.(*NumberValue)
		if !isNum {
			return nil, fmt.Errorf("non-number value in add: %v", v.InspectStr())
		}
		total *= asNum.Val
	}
	return &NumberValue{
		Val: total,
	}, nil
}

func divFn(c *EvalContext, vals ...Value) (Value, error) {
	total := float64(1)
	for i, v := range vals {
		asNum, isNum := v.(*NumberValue)
		if !isNum {
			return nil, fmt.Errorf("non-number value in add: %v", v.InspectStr())
		}
		if i == 0 {
			total = asNum.Val
		} else {
			total /= asNum.Val
		}
	}
	return &NumberValue{
		Val: total,
	}, nil
}

//
// Comparison operator built-in
//

func eqNumFn(ec *EvalContext, vals ...Value) (Value, error) {
	if len(vals) != 2 {
		return nil, fmt.Errorf("eq expects 2 arguments; got %d", len(vals))
	}
	v1, v2 := vals[0], vals[1]
	v1AsNum, v1IsNum := v1.(*NumberValue)
	v2AsNum, v2IsNum := v2.(*NumberValue)
	if !v1IsNum {
		return nil, fmt.Errorf("eq expects number arguments")
	}
	if !v2IsNum {
		return nil, fmt.Errorf("eq expects number arguments")
	}
	return &BoolValue{
		Val: v1AsNum.Val == v2AsNum.Val,
	}, nil
}

func gtNumFn(ec *EvalContext, vals ...Value) (Value, error) {
	if len(vals) != 2 {
		return nil, fmt.Errorf("gt expects 2 arguments; got %d", len(vals))
	}
	v1, v2 := vals[0], vals[1]
	v1AsNum, v1IsNum := v1.(*NumberValue)
	v2AsNum, v2IsNum := v2.(*NumberValue)
	if !v1IsNum {
		return nil, fmt.Errorf("gt expects number arguments")
	}
	if !v2IsNum {
		return nil, fmt.Errorf("gt expects number arguments")
	}
	return &BoolValue{
		Val: v1AsNum.Val > v2AsNum.Val,
	}, nil
}

func ltNumFn(ec *EvalContext, vals ...Value) (Value, error) {
	if len(vals) != 2 {
		return nil, fmt.Errorf("lt expects 2 arguments; got %d", len(vals))
	}
	v1, v2 := vals[0], vals[1]
	v1AsNum, v1IsNum := v1.(*NumberValue)
	v2AsNum, v2IsNum := v2.(*NumberValue)
	if !v1IsNum {
		return nil, fmt.Errorf("lt expects number arguments")
	}
	if !v2IsNum {
		return nil, fmt.Errorf("lt expects number arguments")
	}
	return &BoolValue{
		Val: v1AsNum.Val < v2AsNum.Val,
	}, nil
}

func gteNumFn(ec *EvalContext, vals ...Value) (Value, error) {
	if len(vals) != 2 {
		return nil, fmt.Errorf("gte expects 2 arguments; got %d", len(vals))
	}
	v1, v2 := vals[0], vals[1]
	v1AsNum, v1IsNum := v1.(*NumberValue)
	v2AsNum, v2IsNum := v2.(*NumberValue)
	if !v1IsNum {
		return nil, fmt.Errorf("gte expects number arguments")
	}
	if !v2IsNum {
		return nil, fmt.Errorf("gte expects number arguments")
	}
	return &BoolValue{
		Val: v1AsNum.Val >= v2AsNum.Val,
	}, nil
}

func lteNumFn(ec *EvalContext, vals ...Value) (Value, error) {
	if len(vals) != 2 {
		return nil, fmt.Errorf("lte expects 2 arguments; got %d", len(vals))
	}
	v1, v2 := vals[0], vals[1]
	v1AsNum, v1IsNum := v1.(*NumberValue)
	v2AsNum, v2IsNum := v2.(*NumberValue)
	if !v1IsNum {
		return nil, fmt.Errorf("lte expects number arguments")
	}
	if !v2IsNum {
		return nil, fmt.Errorf("lte expects number arguments")
	}
	return &BoolValue{
		Val: v1AsNum.Val <= v2AsNum.Val,
	}, nil
}
