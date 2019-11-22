package golisp2

import "fmt"

type (
	// IdentLiteral is a representation of an identifier in the interpreted
	// environment, whose value is resolved by the context it is evaluated in.
	IdentLiteral struct {
		// note (bs): I'd like to eventually make it so that identifiers could be
		// "compound lookups"; e.g. "Foo.Bar.A"; in which case I think this should
		// not just be a string. Arguably, that should have it's own datatype
		// anyway.
		Val string
		Pos ScannerPosition
	}

	// NumberLiteral is a representation of a number literal within the
	// interpreted environment.
	NumberLiteral struct {
		Num float64
		Pos ScannerPosition
	}

	// NilLiteral is a representation of an null literal within the interpreted
	// environment.
	NilLiteral struct {
		Pos ScannerPosition
	}

	// StringLiteral is a representation of a string literal within the
	// interpreted environment.
	StringLiteral struct {
		Str string
		Pos ScannerPosition
	}

	// BoolLiteral is a representation of a boolean literal within the interpreted
	// environment.
	BoolLiteral struct {
		Bool bool
		Pos  ScannerPosition
	}

	// FuncLiteral is a representation of a basic function declaration/assignment
	// within the interpreted environment.
	FuncLiteral struct {
		// Name is the function identifier as it appears in the code.
		Name string

		// Fn is the function body the function value references.
		Fn func(*EvalContext, ...Value) (Value, error)

		Pos ScannerPosition
	}
)

// NewIdentLiteral instantiates a new identifier literal with the given
// identifier token.
func NewIdentLiteral(ident string) *IdentLiteral {
	return &IdentLiteral{
		Val: ident,
	}
}

// Eval will traverse the context for the identifier and return nil if the value
// is not defined.
//
// todo (bs): consider making failed resolution an error. In this case, it
// should be a "severe error" that bubbles back and most likely halts execution.
// It's *possible* the right way to handle that is by creating a modified value
// interface that can directly support the notion of error.
func (iv *IdentLiteral) Eval(ec *EvalContext) (Value, error) {
	v, ok := ec.Resolve(iv.Val)
	if !ok {
		return &NilValue{}, nil
	}
	return v, nil
}

// CodeStr will return the code representation of the ident value.
func (iv *IdentLiteral) CodeStr() string {
	return iv.Val
}

// SourcePos is the location in source this value came from.
func (iv *IdentLiteral) SourcePos() ScannerPosition {
	return iv.Pos
}

// NewNumberLiteral instantiates a new number literal with the given value.
func NewNumberLiteral(v float64) *NumberLiteral {
	return &NumberLiteral{
		Num: v,
	}
}

// Eval just returns itself.
func (nv *NumberLiteral) Eval(*EvalContext) (Value, error) {
	return &NumberValue{
		Val: nv.Num,
	}, nil
}

// CodeStr will return the code representation of the number value.
func (nv *NumberLiteral) CodeStr() string {
	// todo (bs): this isn't wrong, exactly, but consider printing integers as
	// integers. Of course, that starts getting into the deeper issue of how just
	// having floats is too primitive and there really need to be integers.
	return fmt.Sprintf("%f", nv.Num)
}

// SourcePos is the location in source this value came from.
func (nv *NumberLiteral) SourcePos() ScannerPosition {
	return nv.Pos
}

// NewNilLiteral creates a new nil value.
func NewNilLiteral() *NilLiteral {
	return &NilLiteral{}
}

// Eval returns the nil value.
func (nv *NilLiteral) Eval(*EvalContext) (Value, error) {
	// note (bs): not sure about this. In general, I feel like eval needs to be
	// more intelligent
	return &NilValue{}, nil
}

// CodeStr will return the code representation of the nil value.
func (nv *NilLiteral) CodeStr() string {
	return fmt.Sprintf("nil")
}

// SourcePos is the location in source this value came from.
func (nv *NilLiteral) SourcePos() ScannerPosition {
	return nv.Pos
}

// NewStringLiteral creates a new string literal from the given string.
func NewStringLiteral(str string) *StringLiteral {
	return &StringLiteral{
		Str: str,
	}
}

// Eval returns the string value.
func (sv *StringLiteral) Eval(*EvalContext) (Value, error) {
	return &StringValue{
		Val: sv.Str,
	}, nil
}

// CodeStr will return the code representation of the string value.
func (sv *StringLiteral) CodeStr() string {
	// note (bs): this doesn't matter now as it's not supported, but just note
	// that this doesn't work with multiline strings
	return fmt.Sprintf("\"%s\"", sv.Str)
}

// SourcePos is the location in source this value came from.
func (sv *StringLiteral) SourcePos() ScannerPosition {
	return sv.Pos
}

// NewBoolLiteral creates a bool literal with the given value.
func NewBoolLiteral(v bool) *BoolLiteral {
	return &BoolLiteral{
		Bool: v,
	}
}

// Eval returns the bool value.
func (bv *BoolLiteral) Eval(*EvalContext) (Value, error) {
	return &BoolValue{
		Val: bv.Bool,
	}, nil
}

// CodeStr will return the code representation of the boolean value.
func (bv *BoolLiteral) CodeStr() string {
	if bv.Bool {
		return "true"
	}
	return "false"
}

// SourcePos is the location in source this value came from.
func (bv *BoolLiteral) SourcePos() ScannerPosition {
	return bv.Pos
}

// NewFuncLiteral creates a function literal with the given value.
func NewFuncLiteral(
	name string,
	fn func(*EvalContext, ...Value) (Value, error),
) *FuncLiteral {
	return &FuncLiteral{
		Fn: fn,
	}
}

// Eval evaluates the function using the provided context.
func (fv *FuncLiteral) Eval(ec *EvalContext) (Value, error) {
	return &FuncValue{
		Fn: fv.Fn,
	}, nil
}

// CodeStr will return the code representation of the function value.
func (fv *FuncLiteral) CodeStr() string {
	return fv.Name
}

// SourcePos is the location in source this value came from.
func (fv *FuncLiteral) SourcePos() ScannerPosition {
	return fv.Pos
}
