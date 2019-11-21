package golisp2

type (
	// EvalContext is the context on evaluation. It contains a resolvable set of
	// identifiers->values that can be chained.
	EvalContext struct {
		parent *EvalContext
		vals   map[string]Value
	}
)

// NewContext returns a new context with no parent. initialVals contains any
// values that the context should be initialized with; it can be left nil.
func NewContext(initialVals map[string]Value) *EvalContext {
	vals := map[string]Value{}
	for k, v := range initialVals {
		vals[k] = v
	}
	return &EvalContext{
		vals: vals,
	}
}

// SubContext creates a new context with the current context as it's parent.
func (ec *EvalContext) SubContext(initialVals map[string]Value) *EvalContext {
	sub := NewContext(initialVals)
	sub.parent = ec
	return sub
}

// Add extends the current context with the provided value.
func (ec *EvalContext) Add(ident string, val Value) {
	ec.vals[ident] = val
}

// Resolve traverses the expr for the given ident. Will return it if found;
// otherwise a nil value and "false".
func (ec *EvalContext) Resolve(ident string) (Value, bool) {
	if ec == nil {
		return NewNilLiteral(), false
	}
	if v, ok := ec.vals[ident]; ok {
		return v, true
	}
	return ec.parent.Resolve(ident)
}
