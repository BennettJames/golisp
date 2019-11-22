package golisp2

import "fmt"

type (
	// CellValue is a representation of a pair of values within the interpreted
	// environment. This can be composed to represent lists with standard car/cdr
	// operators.
	CellValue struct {
		Left, Right Value
	}

	// NumberValue is a representation of a number within the interpreted
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
		// Fn is the function body the function value references.
		Fn func(*EvalContext, ...Value) (Value, error)
	}
)

// NewCellValue creates a cell with the given left/right values. Either can be
// 'nil'.
func NewCellValue(left, right Value) *CellValue {
	if left == nil {
		left = &NilValue{}
	}
	if right == nil {
		right = &NilValue{}
	}
	return &CellValue{
		Left:  left,
		Right: right,
	}
}

// Eval returns the cell.
func (cv *CellValue) Eval(*EvalContext) (Value, error) {
	return cv, nil
}

// InspectStr outputs the contents of all the cells.
func (cv *CellValue) InspectStr() string {
	// todo (bs): if second cell is a node, treat this different
	return fmt.Sprintf("(%s . %s)", cv.Left.InspectStr(), cv.Right.InspectStr())
}

// InspectStr prints the number.
func (nv *NumberValue) InspectStr() string {
	return fmt.Sprintf("%f", nv.Val)
}

// InspectStr outputs "nil".
func (nv *NilValue) InspectStr() string {
	return "nil"
}

// InspectStr prints the string.
func (sv *StringValue) InspectStr() string {
	return fmt.Sprintf("\"%s\"", sv.Val)
}

// InspectStr prints "true"/"false" based on the value.
func (bv *BoolValue) InspectStr() string {
	return fmt.Sprintf("%t", bv.Val)
}

// InspectStr outputs some information about the function.
func (fv *FuncValue) InspectStr() string {
	// note (bs): probably want to customize this to print some details about the
	// function itself. That will involve (optionally) retaining the declaration
	// name of the function.
	return fmt.Sprintf("<func>")
}
