package golisp2

import (
	"fmt"
	"io"
	"strconv"
)

// ParseTokens reads in the tokens, and converts them to a set of expressions.
// Returns the set, and any parse errors that are encountered in the process.
func ParseTokens(ts *TokenScanner) ([]Expr, error) {
	ts.Advance() // initialize scan
	exprs, exprsErr := maybeParseExprs(ts)
	if exprsErr != nil {
		return nil, exprsErr
	}
	if !ts.Done() {
		// todo (bs): improve this error
		return nil, fmt.Errorf("scan completed before end of input")
	}
	return exprs, nil
}

// maybeParseExprs will read as many expressions as it can, until it hits EOF or
// a close boundary character.
func maybeParseExprs(ts *TokenScanner) ([]Expr, error) {
	exprs := []Expr{}
	for {
		maybeExpr, maybeExprErr := maybeParseExpr(ts)
		if maybeExprErr != nil {
			return nil, maybeExprErr
		}
		if maybeExpr == nil {
			return exprs, nil
		}
		exprs = append(exprs, maybeExpr)
	}
}

// maybeParseExpr will attempt to read a complete expression from the scanner.
// Will return an error if any problems are encountered. Will return (nil, nil)
// if the scanner has no more expressions, or the end of a block is reached.
func maybeParseExpr(ts *TokenScanner) (Expr, error) {
	maybeNextToken := ts.Token()
	if maybeNextToken == nil {
		if ts.Err() != nil && ts.Err() != io.EOF {
			// todo (bs): wrap error using go 1.13
			return nil, ts.Err()
		}
		return nil, nil
	}
	nextToken := *maybeNextToken

	switch nextToken.Typ {
	case CloseParenTT:
		return nil, nil
	case OpenParenTT:
		return tryParseCall(ts)
	case IdentTT:
		ts.Advance()
		return parseIdentValue(nextToken)
	case OpTT:
		ts.Advance()
		return parseOpValue(nextToken)
	case NumberTT:
		ts.Advance()
		return parseNumberValue(nextToken)
	case StringTT:
		ts.Advance()
		return parseStringValue(nextToken)
	default:
		return nil, NewParseError("invalid token", nextToken)
	}
}

func tryParseCall(ts *TokenScanner) (Expr, error) {
	maybeStartToken := ts.Token()
	if maybeStartToken == nil {
		return nil, fmt.Errorf("parse call called on empty scanner")
	}
	startToken := *maybeStartToken
	if startToken.Typ != OpenParenTT {
		return nil, NewParseError(
			"call expression must start with open paren", startToken)
	}

	ts.Advance()
	maybeNextToken := ts.Token()
	if maybeNextToken == nil {
		return nil, NewParseError("parse ended inside of call", startToken)
	}
	nextToken := *maybeNextToken
	if nextToken.Typ == IdentTT {
		// note (bs): strongly consider making this a data structure; will make
		// rejecting usages of reserved words as idents much easier
		switch nextToken.Value {
		case "if":
			return tryParseIfTail(ts)
		case "fn":
			return tryParseFnTail(ts)
		case "let":
			return tryParseLetTail(ts)
		case "defun":
			panic("defun not implemented")
		case "import":
			panic("import not implemented")
		}
	}

	return tryParseCallTail(ts)
}

func tryParseCallTail(ts *TokenScanner) (Expr, error) {
	bodyExprs, bodyExprsErr := maybeParseExprs(ts)
	if bodyExprsErr != nil {
		return nil, bodyExprsErr
	}
	if err := expectCallClose(ts); err != nil {
		return nil, err
	}
	return &CallExpr{
		Exprs: bodyExprs,
	}, nil
}

func parseStringValue(token ScannedToken) (*StringValue, error) {
	v := token.Value
	if len(v) == 0 {
		return NewStringValue(""), nil
	}
	leadI, tailI := 0, len(v)
	if v[0] == '"' {
		leadI = 1
	}
	if len(v) > 1 && v[len(v)-1] == '"' {
		tailI = len(v) - 1
	}
	return NewStringValue(v[leadI:tailI]), nil
}

func parseIdentValue(token ScannedToken) (Value, error) {
	// todo (bs): this should search for certain reserved words, and reject them.
	// e.g. any of the "structural builtins" like if, defun, let, etc.

	switch token.Value {
	case "nil":
		return NewNilValue(), nil
	case "true":
		return NewBoolValue(true), nil
	case "false":
		return NewBoolValue(false), nil
	default:
		return NewIdentValue(token.Value), nil
	}
}

func parseNumberValue(token ScannedToken) (*NumberValue, error) {
	// todo (bs): given that this is, you know, a *parser*, it's awfully clumsy to
	// outsource the final number parsing to Go. The manual parse should be able
	// to correctly map this to a number.
	f, e := strconv.ParseFloat(token.Value, 64)
	if e != nil {
		// todo (bs): this should wrap the error
		return nil, e
	}
	return NewNumberValue(f), nil
}

func parseOpValue(token ScannedToken) (*FuncValue, error) {
	// todo (bs): strongly consider moving this to a map rather than a case
	// statement
	switch token.Value {
	case "+":
		return NewFuncValue("+", addFn), nil
	case "-":
		return NewFuncValue("-", subFn), nil
	case "*":
		return NewFuncValue("*", multFn), nil
	case "/":
		return NewFuncValue("/", divFn), nil
	case "==":
		return NewFuncValue("==", eqNumFn), nil
	case "<":
		return NewFuncValue("<", ltNumFn), nil
	case ">":
		return NewFuncValue(">", gtNumFn), nil
	case "<=":
		return NewFuncValue("<=", lteNumFn), nil
	case ">=":
		return NewFuncValue(">=", gteNumFn), nil
	default:
		return nil, NewParseError("unrecognized operator", token)
	}
}

func tryParseIfTail(ts *TokenScanner) (Expr, error) {
	maybeStartToken := ts.Token()
	if maybeStartToken == nil {
		return nil, fmt.Errorf("tryParseIfTail called on empty scanner")
	}
	startToken := *maybeStartToken
	if startToken.Typ != IdentTT || startToken.Value != "if" {
		return nil, NewParseError("tryParseIfTail called on non-if", startToken)
	}
	ts.Advance()

	ifBody, ifBodyErr := maybeParseExprs(ts)
	if ifBodyErr != nil {
		return nil, ifBodyErr
	}
	var cond, case1, case2 Expr
	if len(ifBody) == 0 {
		return nil, NewParseError("if statement must have condition", startToken)
	}
	cond = ifBody[0]
	if len(ifBody) > 1 {
		case1 = ifBody[1]
	}
	if len(ifBody) > 2 {
		case2 = ifBody[2]
	}
	if len(ifBody) > 3 {
		return nil, NewParseError(
			"if statement can have no more than 3 expressions", startToken)
	}
	if err := expectCallClose(ts); err != nil {
		return nil, err
	}

	return &IfExpr{
		Cond:  wrapNilExpr(cond),
		Case1: wrapNilExpr(case1),
		Case2: wrapNilExpr(case2),
	}, nil
}

func tryParseFnTail(ts *TokenScanner) (Expr, error) {
	maybeStartToken := ts.Token()
	if maybeStartToken == nil {
		return nil, fmt.Errorf("tryParseFnTail called on empty scanner")
	}
	startToken := *maybeStartToken
	if startToken.Typ != IdentTT || startToken.Value != "fn" {
		return nil, NewParseError("tryParseFnTail called on non-fn", startToken)
	}
	ts.Advance()

	args, argsErr := tryParseFnArgs(ts)
	if argsErr != nil {
		return nil, argsErr
	}
	bodyExprs, bodyExprsErr := maybeParseExprs(ts)
	if bodyExprsErr != nil {
		return nil, bodyExprsErr
	}
	if err := expectCallClose(ts); err != nil {
		return nil, err
	}

	return &FnExpr{
		Args: args,
		Body: bodyExprs,
	}, nil
}

func tryParseFnArgs(ts *TokenScanner) ([]Arg, error) {
	if err := expectCallOpen(ts); err != nil {
		return nil, err
	}
	args := []Arg{}
	for {
		maybeNextToken := ts.Token()
		if maybeNextToken == nil {
			// todo (bs): add proper parse error info here
			return nil, fmt.Errorf("unexpected end of input")
		}
		nextToken := *maybeNextToken
		ts.Advance()
		switch nextToken.Typ {
		case IdentTT:
			args = append(args, Arg{
				Ident: nextToken.Value,
			})
		case CloseParenTT:
			return args, nil
		default:
			return nil, NewParseError("args can only contain idents", nextToken)
		}
	}
}

func tryParseLetTail(ts *TokenScanner) (Expr, error) {
	maybeStartToken := ts.Token()
	if maybeStartToken == nil {
		return nil, fmt.Errorf("tryParseLetTail called on empty scanner")
	}
	startToken := *maybeStartToken
	if startToken.Typ != IdentTT || startToken.Value != "let" {
		return nil, NewParseError("tryParseLetTail called on non-let", startToken)
	}
	ts.Advance()

	letExprs, letExprsErr := maybeParseExprs(ts)
	if letExprsErr != nil {
		return nil, letExprsErr
	}
	if len(letExprs) != 2 {
		return nil, NewParseError(
			fmt.Sprintf("let expects 2 arguments, got %d",
				len(letExprs)), startToken)
	}
	asIdent, isIdent := letExprs[0].(*IdentValue)
	if !isIdent {
		return nil, NewParseError(
			"let expects an ident as first argument", startToken)
	}
	val := letExprs[1]
	if err := expectCallClose(ts); err != nil {
		return nil, err
	}

	return &LetExpr{
		Ident: asIdent,
		Value: val,
	}, nil
}

// expectCallOpen will read a open paren from the scanner and advance, or
// return an error.
func expectCallOpen(ts *TokenScanner) error {
	maybeNext := ts.Token()
	if maybeNext == nil {
		// todo (bs): turn into parse error w/ location info
		return fmt.Errorf("unexpected end of input")
	}
	next := *maybeNext
	if next.Typ != OpenParenTT {
		return NewParseError("expected open paren", next)
	}
	ts.Advance()
	return nil
}

// expectCallClose will read a close paren from the scanner and advance, or
// return an error.
func expectCallClose(ts *TokenScanner) error {
	maybeNext := ts.Token()
	if maybeNext == nil {
		// todo (bs): turn into parse error w/ location info
		return fmt.Errorf("unexpected end of input")
	}
	next := *maybeNext
	if next.Typ != CloseParenTT {
		return NewParseError("expected close paren", next)
	}
	ts.Advance()
	return nil
}

// wrapNilExpr will return a nil expr if e is nil.
func wrapNilExpr(e Expr) Expr {
	if e == nil {
		return NewNilValue()
	}
	return e
}
