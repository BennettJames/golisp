package golisp2

import "fmt"

type (
	// ArgMapper is a utility that makes it easier to map lists of values to
	ArgMapper struct {
		iter valueIterator
		err  error
	}

	// valueIterator is a generic way to traverse/process a set of value-like
	// objects.
	valueIterator interface {
		// Next returns the next value in the iterator. If none are left, (nil, nil)
		// will be returned.
		Next() (Value, error)
	}

	// valueSet implements valueIterator through simply iterating through the set.
	valueSet struct {
		vals     []Value
		argIndex int
	}

	// exprSet implements valueIterator by evaluating expressions on demand.
	exprSet struct {
		ec       *EvalContext
		exprs    []Expr
		argIndex int
	}
)

// ArgMapperValues creates an argument mapper for the provided values.
func ArgMapperValues(vals ...Value) *ArgMapper {
	return &ArgMapper{
		iter: &valueSet{
			vals:     vals,
			argIndex: 0,
		},
		err: nil,
	}
}

// ArgMapperExprs creates an argument mapper for the provided context/expr set.
func ArgMapperExprs(ec *EvalContext, exprs []Expr) *ArgMapper {
	return &ArgMapper{
		iter: &exprSet{
			ec:       ec,
			exprs:    exprs,
			argIndex: 0,
		},
		err: nil,
	}
}

// ReadString will try to read the next argument as a string value, or report an
// error.
func (am *ArgMapper) ReadString(v **StringValue) *ArgMapper {
	switch tV := am.next().(type) {
	case *StringValue:
		*v = tV
	}
	return am
}

// ReadBool will try to read the next argument as a bool value, or report an
// error.
func (am *ArgMapper) ReadBool(v **BoolValue) *ArgMapper {
	switch tV := am.next().(type) {
	case *BoolValue:
		*v = tV
	}
	return am
}

// ReadFunc will try to read the next function as a list value, or report an
// error.
func (am *ArgMapper) ReadFunc(v **FuncValue) *ArgMapper {
	switch tV := am.next().(type) {
	case *FuncValue:
		*v = tV
	}
	return am
}

// ReadNumber will try to read the next argument as a number value, or report an
// error.
func (am *ArgMapper) ReadNumber(v **NumberValue) *ArgMapper {
	switch tV := am.next().(type) {
	case *NumberValue:
		*v = tV
	}
	return am
}

// ReadCell will try to read the next argument as a cell value, or report an
// error.
func (am *ArgMapper) ReadCell(v **CellValue) *ArgMapper {
	switch tV := am.next().(type) {
	case *CellValue:
		*v = tV
	}
	return am
}

// ReadList will try to read the next argument as a list value, or report an
// error.
func (am *ArgMapper) ReadList(v **ListValue) *ArgMapper {
	switch tV := am.next().(type) {
	case *ListValue:
		*v = tV
	}
	return am
}

// ReadMap will try to read the next argument as a map value, or report an
// error.
func (am *ArgMapper) ReadMap(v **MapValue) *ArgMapper {
	switch tV := am.next().(type) {
	case *MapValue:
		*v = tV
	}
	return am
}

// ReadValue will try to read the next argument as any value, or report an
// error.
func (am *ArgMapper) ReadValue(v *Value) *ArgMapper {
	if nextV := am.next(); nextV != nil {
		*v = nextV
	}
	return am
}

// Err returns any encountered errors during the mapping.
func (am *ArgMapper) Err() error {
	return am.err
}

func (am *ArgMapper) next() Value {
	nextV := am.maybeNext()
	if nextV == nil {
		// note (bs): this is a little imprecise; may wish to make it possible to
		// better label the source of errors. That's a broader problem than just
		// this; really.
		am.err = fmt.Errorf("ArgMapper: not enough arguments")
	}
	return nextV
}

func (am *ArgMapper) maybeNext() Value {
	if am.err != nil {
		return nil
	}
	nextV, nextVErr := am.iter.Next()
	if nextVErr != nil {
		am.err = nextVErr
		return nil
	}
	return nextV
}

func (vs *valueSet) Next() (Value, error) {
	if vs.argIndex >= len(vs.vals) {
		return nil, nil
	}
	v := vs.vals[vs.argIndex]
	vs.argIndex++
	return v, nil
}

func (es *exprSet) Next() (Value, error) {
	if es.argIndex >= len(es.exprs) {
		return nil, nil
	}
	v, err := es.exprs[es.argIndex].Eval(es.ec)
	es.argIndex++
	return v, err
}
