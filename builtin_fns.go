package golisp2

import (
	"fmt"
	"math"
	"strings"
)

// BuiltinContext returns a context that contains the full set of builtin
// functions. Note this just includes built-in plain functions; not operators.
func BuiltinContext() *EvalContext {
	return NewContext(map[string]Value{
		"concat": &FuncValue{Fn: concatFn},
		"cons":   &FuncValue{Fn: consFn},
		"car":    &FuncValue{Fn: carFn},
		"cdr":    &FuncValue{Fn: cdrFn},
		"and":    &FuncValue{Fn: andFn},
		"or":     &FuncValue{Fn: orFn},
		"not":    &FuncValue{Fn: notFn},

		"strEq": &FuncValue{Fn: strEqFn},

		"list":       &FuncValue{Fn: listCreateFn},
		"listGet":    &FuncValue{Fn: listGetFn},
		"listFilter": &FuncValue{Fn: listFilterFn},
		"listMap":    &FuncValue{Fn: listMapFn},
		"listReduce": &FuncValue{Fn: listReduceFn},
		"len":        &FuncValue{Fn: lenFn},

		"map":       &FuncValue{Fn: mapCreateFn},
		"mapGet":    &FuncValue{Fn: mapGetFn},
		"mapFilter": &FuncValue{Fn: mapFilterFn},
		"mapMap":    &FuncValue{Fn: mapMapFn},
		"mapReduce": &FuncValue{Fn: mapReduceFn},
		"mapKeys":   &FuncValue{Fn: mapKeysFn},
		"mapValues": &FuncValue{Fn: mapValuesFn},

		"print": &FuncValue{Fn: printFn},
	})
}

//
// Explicit, named built-ins
//

func concatFn(c *EvalContext, vals ...Value) (Value, error) {
	var strVals []*StringValue
	err := ArgMapperValues(vals...).
		ReadStrings(&strVals).
		Complete()
	if err != nil {
		return nil, err
	}
	var sb strings.Builder
	for _, v := range strVals {
		sb.WriteString(v.Val)
	}
	return &StringValue{
		Val: sb.String(),
	}, nil
}

func strEqFn(c *EvalContext, vals ...Value) (Value, error) {
	var v1, v2 *StringValue
	err := ArgMapperValues(vals...).
		ReadString(&v1).
		ReadString(&v2).
		Complete()
	if err != nil {
		return nil, err
	}
	return &BoolValue{
		Val: v1.Val == v2.Val,
	}, nil
}

func consFn(c *EvalContext, vals ...Value) (Value, error) {
	var v1, v2 Value
	err := ArgMapperValues(vals...).
		MaybeReadValue(&v1).
		MaybeReadValue(&v2).
		Complete()
	if err != nil {
		return nil, err
	}
	return NewCellValue(v1, v2), nil
}

func carFn(c *EvalContext, vals ...Value) (Value, error) {
	var v1 *CellValue
	err := ArgMapperValues(vals...).
		ReadCell(&v1).
		Complete()
	if err != nil {
		return nil, err
	}
	return v1.Left, nil
}

func cdrFn(c *EvalContext, vals ...Value) (Value, error) {
	var v1 *CellValue
	err := ArgMapperValues(vals...).
		ReadCell(&v1).
		Complete()
	if err != nil {
		return nil, err
	}
	return v1.Right, nil
}

func andFn(c *EvalContext, vals ...Value) (Value, error) {
	var firstV *BoolValue
	var remainingVals []*BoolValue
	err := ArgMapperValues(vals...).
		ReadBool(&firstV).
		ReadBools(&remainingVals).
		Complete()
	if err != nil {
		return nil, err
	}
	if !firstV.Val {
		return &BoolValue{Val: false}, nil
	}
	for _, v := range remainingVals {
		if !v.Val {
			return &BoolValue{Val: false}, nil
		}
	}
	return &BoolValue{Val: true}, nil
}

func orFn(c *EvalContext, vals ...Value) (Value, error) {
	var firstV *BoolValue
	var remainingVals []*BoolValue
	err := ArgMapperValues(vals...).
		ReadBool(&firstV).
		ReadBools(&remainingVals).
		Complete()
	if err != nil {
		return nil, err
	}
	if firstV.Val {
		return &BoolValue{Val: true}, nil
	}
	for _, v := range remainingVals {
		if v.Val {
			return &BoolValue{Val: true}, nil
		}
	}
	return &BoolValue{Val: false}, nil
}

func notFn(c *EvalContext, vals ...Value) (Value, error) {
	var v1 *BoolValue
	err := ArgMapperValues(vals...).
		ReadBool(&v1).
		Complete()
	if err != nil {
		return nil, err
	}
	return &BoolValue{
		Val: !v1.Val,
	}, nil
}

//
// Mathematical operator built-ins
//

func addFn(c *EvalContext, vals ...Value) (Value, error) {
	var firstVal *NumberValue
	var remainingVals []*NumberValue
	err := ArgMapperValues(vals...).
		ReadNumber(&firstVal).
		ReadNumbers(&remainingVals).
		Complete()
	if err != nil {
		return nil, err
	}
	total := firstVal.Val
	for _, v := range remainingVals {
		total += v.Val
	}
	return &NumberValue{
		Val: total,
	}, nil
}

func subFn(c *EvalContext, vals ...Value) (Value, error) {
	var firstVal *NumberValue
	var remainingVals []*NumberValue
	err := ArgMapperValues(vals...).
		ReadNumber(&firstVal).
		ReadNumbers(&remainingVals).
		Complete()
	if err != nil {
		return nil, err
	}
	if len(remainingVals) == 0 {
		return &NumberValue{
			Val: -firstVal.Val,
		}, nil
	}
	total := firstVal.Val
	for _, v := range remainingVals {
		total -= v.Val
	}
	return &NumberValue{
		Val: total,
	}, nil
}

func multFn(c *EvalContext, vals ...Value) (Value, error) {
	var firstVal *NumberValue
	var remainingVals []*NumberValue
	err := ArgMapperValues(vals...).
		ReadNumber(&firstVal).
		ReadNumbers(&remainingVals).
		Complete()
	if err != nil {
		return nil, err
	}
	total := firstVal.Val
	for _, v := range remainingVals {
		total *= v.Val
	}
	return &NumberValue{
		Val: total,
	}, nil
}

func divFn(c *EvalContext, vals ...Value) (Value, error) {
	var firstVal *NumberValue
	var remainingVals []*NumberValue
	err := ArgMapperValues(vals...).
		ReadNumber(&firstVal).
		ReadNumbers(&remainingVals).
		Complete()
	if err != nil {
		return nil, err
	}
	total := firstVal.Val
	for _, v := range remainingVals {
		total /= v.Val
	}
	return &NumberValue{
		Val: total,
	}, nil
}

//
// Comparison operator built-in
//

func eqNumFn(ec *EvalContext, vals ...Value) (Value, error) {
	var v1, v2 *NumberValue
	err := ArgMapperValues(vals...).
		ReadNumber(&v1).
		ReadNumber(&v2).
		Complete()
	if err != nil {
		return nil, err
	}
	return &BoolValue{
		Val: v1.Val == v2.Val,
	}, nil
}

func gtNumFn(ec *EvalContext, vals ...Value) (Value, error) {
	var v1, v2 *NumberValue
	err := ArgMapperValues(vals...).
		ReadNumber(&v1).
		ReadNumber(&v2).
		Complete()
	if err != nil {
		return nil, err
	}
	return &BoolValue{
		Val: v1.Val > v2.Val,
	}, nil
}

func ltNumFn(ec *EvalContext, vals ...Value) (Value, error) {
	var v1, v2 *NumberValue
	err := ArgMapperValues(vals...).
		ReadNumber(&v1).
		ReadNumber(&v2).
		Complete()
	if err != nil {
		return nil, err
	}
	return &BoolValue{
		Val: v1.Val < v2.Val,
	}, nil
}

func gteNumFn(ec *EvalContext, vals ...Value) (Value, error) {
	var v1, v2 *NumberValue
	err := ArgMapperValues(vals...).
		ReadNumber(&v1).
		ReadNumber(&v2).
		Complete()
	if err != nil {
		return nil, err
	}
	return &BoolValue{
		Val: v1.Val >= v2.Val,
	}, nil
}

func lteNumFn(ec *EvalContext, vals ...Value) (Value, error) {
	var v1, v2 *NumberValue
	err := ArgMapperValues(vals...).
		ReadNumber(&v1).
		ReadNumber(&v2).
		Complete()
	if err != nil {
		return nil, err
	}
	return &BoolValue{
		Val: v1.Val <= v2.Val,
	}, nil
}

//
// List functions
//

// listCreateFn creates a new list out of the given arguments.
func listCreateFn(ec *EvalContext, vals ...Value) (Value, error) {
	return &ListValue{
		Vals: vals,
	}, nil
}

// listGetFn gets and returns the given index from the list. If it doesn't exit;
// returns zero.
func listGetFn(ec *EvalContext, vals ...Value) (Value, error) {
	var asList *ListValue
	var asNum *NumberValue
	err := ArgMapperValues(vals...).
		ReadList(&asList).
		ReadNumber(&asNum).
		Complete()
	if err != nil {
		return nil, err
	}

	index := int(math.Floor(asNum.Val))
	if index < 0 || index >= len(asList.Vals) {
		return nil, fmt.Errorf("listGet out of bounds")
	}
	return asList.Vals[index], nil
}

// listFilterFn expects a list and a function argument. The function will take an
// element, and return either true or false. It will be called on each element
// of the list, and all values that are marked true will be collected and
// returned in a new list.
func listFilterFn(ec *EvalContext, vals ...Value) (Value, error) {
	var asList *ListValue
	var asFn *FuncValue
	err := ArgMapperValues(vals...).
		ReadList(&asList).
		ReadFunc(&asFn).
		Complete()
	if err != nil {
		return nil, err
	}

	filteredVals := []Value{}
	for _, v := range asList.Vals {
		// todo (bs): double check that this couldn't contaminate the scope
		filterVal, filterErr := asFn.Fn(ec, v)
		if filterErr != nil {
			return nil, fmt.Errorf("listFilter encountered an error: %w", filterErr)
		}
		switch tV := filterVal.(type) {
		case *NilValue:
			continue
		case *BoolValue:
			if tV.Val {
				filteredVals = append(filteredVals, v)
			}
		default:
			return nil, fmt.Errorf("listFilter fn must return boolean")
		}
	}

	return &ListValue{
		Vals: filteredVals,
	}, nil
}

// listMapFn expects a list and a function argument. The function will take an
// element and return an element. It will be called on each element on the list;
// and the returned values will be returned in a new list.
func listMapFn(ec *EvalContext, vals ...Value) (Value, error) {
	var asList *ListValue
	var asFn *FuncValue
	err := ArgMapperValues(vals...).
		ReadList(&asList).
		ReadFunc(&asFn).
		Complete()
	if err != nil {
		return nil, err
	}

	mappedVals := []Value{}
	for _, v := range asList.Vals {
		mapVal, mapErr := asFn.Fn(ec, v)
		if mapErr != nil {
			return nil, fmt.Errorf("listMap encountered an error: %w", mapErr)
		}
		mappedVals = append(mappedVals, mapVal)
	}

	return &ListValue{
		Vals: mappedVals,
	}, nil
}

// listReduceFn expects a value, list, and a function argument. The value is the
// "initial value" of the reduction. The function take two arguments; the
// "reduced value" and an element from the list. It will be called with the
// initial value, and iteratively called with the results of the past map and
// the next element in the list.
func listReduceFn(ec *EvalContext, vals ...Value) (Value, error) {
	var initVal Value
	var asList *ListValue
	var asFn *FuncValue
	err := ArgMapperValues(vals...).
		ReadValue(&initVal).
		ReadList(&asList).
		ReadFunc(&asFn).
		Complete()
	if err != nil {
		return nil, err
	}

	reducedVal := initVal
	for _, v := range asList.Vals {
		innerRVal, err := asFn.Fn(ec, reducedVal, v)
		if err != nil {
			return nil, fmt.Errorf("listReduce encountered an error: %w", err)
		}
		reducedVal = innerRVal
	}

	return reducedVal, nil
}

//
// Map functions
//

// mapCreateFn creates a new map out of the given arguments.
func mapCreateFn(ec *EvalContext, vals ...Value) (Value, error) {
	if len(vals)%2 != 0 {
		return nil, fmt.Errorf("map expects even number of arguments; got %d", len(vals))
	}

	mapVals := map[string]Value{}
	for i := 0; i+1 < len(vals); i += 2 {
		k, v := vals[i], vals[i+1]
		asStr, isStr := k.(*StringValue)
		if !isStr {
			return nil, fmt.Errorf("map expects hashable keys")
		}
		mapVals[asStr.Val] = v
	}

	return &MapValue{
		Vals: mapVals,
	}, nil
}

// mapGetFn gets and returns the given key from the map. If it doesn't exist;
// returns nil.
func mapGetFn(ec *EvalContext, vals ...Value) (Value, error) {
	var asMap *MapValue
	var asStr *StringValue
	err := ArgMapperValues(vals...).
		ReadMap(&asMap).
		ReadString(&asStr).
		Complete()
	if err != nil {
		return nil, err
	}

	val, hasVal := asMap.Vals[asStr.Val]
	if !hasVal {
		return &NilValue{}, nil
	}
	return val, nil
}

// mapFilterFn expects a map and a function argument. The function will take a
// key/value pair, and return either true or false. It will be called on each
// element of the list, and all values that are marked true will be collected
// and returned in a new list.
func mapFilterFn(ec *EvalContext, vals ...Value) (Value, error) {
	var asMap *MapValue
	var asFn *FuncValue
	err := ArgMapperValues(vals...).
		ReadMap(&asMap).
		ReadFunc(&asFn).
		Complete()
	if err != nil {
		return nil, err
	}

	filteredVals := map[string]Value{}
	for k, v := range asMap.Vals {
		filterVal, filterErr := asFn.Fn(ec, &StringValue{Val: k}, v)
		if filterErr != nil {
			return nil, fmt.Errorf("mapFilter encountered an error: %w", filterErr)
		}
		switch tV := filterVal.(type) {
		case *NilValue:
			continue
		case *BoolValue:
			if tV.Val {
				filteredVals[k] = v
			}
		default:
			return nil, fmt.Errorf("mapFilter fn must return boolean")
		}
	}

	return &MapValue{
		Vals: filteredVals,
	}, nil
}

// mapMapFn expects a map and a function argument. The function will take an
// key/value pair and return an updated value. It will be called on each element
// on the map; and the returned values will be returned in a new map.
func mapMapFn(ec *EvalContext, vals ...Value) (Value, error) {
	var asMap *MapValue
	var asFn *FuncValue
	err := ArgMapperValues(vals...).
		ReadMap(&asMap).
		ReadFunc(&asFn).
		Complete()
	if err != nil {
		return nil, err
	}

	mappedVals := map[string]Value{}
	for k, v := range asMap.Vals {
		mappedVal, mapErr := asFn.Fn(ec, &StringValue{Val: k}, v)
		if mapErr != nil {
			return nil, fmt.Errorf("mapMap encountered an error: %w", mapErr)
		}
		mappedVals[k] = mappedVal
	}

	return &MapValue{
		Vals: mappedVals,
	}, nil
}

// mapReduceFn expects a value, map, and a function argument. The value is the
// "initial value" of the reduction. The function take three arguments; the
// "reduced value" and a key/value pair from the map. It will be called with the
// initial value, and iteratively called with the results of the past map and
// the next element in the map.
func mapReduceFn(ec *EvalContext, vals ...Value) (Value, error) {
	var initVal Value
	var asMap *MapValue
	var asFn *FuncValue
	err := ArgMapperValues(vals...).
		ReadValue(&initVal).
		ReadMap(&asMap).
		ReadFunc(&asFn).
		Complete()
	if err != nil {
		return nil, err
	}

	reducedVal := initVal
	for k, v := range asMap.Vals {
		innerRVal, err := asFn.Fn(ec, reducedVal, &StringValue{Val: k}, v)
		if err != nil {
			return nil, fmt.Errorf("mapReduce encountered an error: %w", err)
		}
		reducedVal = innerRVal
	}

	return reducedVal, nil
}

// mapKeysFn takes a map and returns it's keys as a list.
func mapKeysFn(ec *EvalContext, vals ...Value) (Value, error) {
	var asMap *MapValue
	err := ArgMapperValues(vals...).
		ReadMap(&asMap).
		Complete()
	if err != nil {
		return nil, err
	}

	keys := make([]Value, 0, len(asMap.Vals))
	for k := range asMap.Vals {
		keys = append(keys, &StringValue{Val: k})
	}

	return &ListValue{
		Vals: keys,
	}, nil
}

// mapValuesFn takes a map and returns it's values as a list.
func mapValuesFn(ec *EvalContext, vals ...Value) (Value, error) {
	var asMap *MapValue
	err := ArgMapperValues(vals...).
		ReadMap(&asMap).
		Complete()
	if err != nil {
		return nil, err
	}

	values := make([]Value, 0, len(asMap.Vals))
	for _, v := range asMap.Vals {
		values = append(values, v)
	}

	return &ListValue{
		Vals: values,
	}, nil
}

//
// Misc values
//

// printFn outputs the values in stdout.
func printFn(ec *EvalContext, vals ...Value) (Value, error) {
	for i, v := range vals {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(v.InspectStr())
	}
	fmt.Println()
	return &NilValue{}, nil
}

// lenFn will return the length of maps, lists, and strings.
func lenFn(ec *EvalContext, vals ...Value) (Value, error) {
	var val Value
	err := ArgMapperValues(vals...).
		ReadValue(&val).
		Complete()
	if err != nil {
		return nil, err
	}

	// ques (bs): should this be solved via subtyping?
	switch tV := val.(type) {
	case *ListValue:
		return &NumberValue{
			Val: float64(len(tV.Vals)),
		}, nil
	case *StringValue:
		return &NumberValue{
			Val: float64(len(tV.Val)),
		}, nil
	case *MapValue:
		return &NumberValue{
			Val: float64(len(tV.Vals)),
		}, nil
	default:
		return nil, fmt.Errorf("Cannot get length of type %T", tV)
	}
}
