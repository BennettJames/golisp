package golisp2

import (
	"fmt"
	"strconv"
)

// ParseTokens reads in the tokens, and converts them to a set of expressions.
// Returns the set, and any parse errors that are encountered in the process.
func ParseTokens(ts *TokenScanner) ([]Expr, error) {
	exprs := []Expr{}
	for !ts.Done() {
		maybeNext := ts.Next()
		if maybeNext == nil {
			if ts.Done() {
				break
			} else {
				return nil, fmt.Errorf("unexpected nil token")
			}
		}
		next := *maybeNext

		switch next.Typ {
		case CommentTT:
			continue
		case OpenParenTT:
			expr, exprErr := parseCallExpr(ts)
			if exprErr != nil {
				return nil, exprErr
			}
			exprs = append(exprs, expr)
		default:
			return nil, NewParseError("unexpected top level token", next)
		}
	}
	return exprs, nil
}

func parseCallExpr(ts *TokenScanner) (Expr, error) {

	exprs := []Expr{}
	for !ts.Done() {
		maybeNext := ts.Next()
		if maybeNext == nil {
			break
		}
		next := *maybeNext

		switch next.Typ {
		case CloseParenTT:
			if len(exprs) > 0 {
				// note (bs): this is awfully clumsy. I think it'd be better to actually
				// send out a subparser to be responsible for handling built-ins/macros,
				// but I'll worry about that later.

				// fixme (bs): if there is an error in any of these conversions, it
				// (most likely) should be counted as a parse error. But: as they're
				// pre-converted to expressions and expressions do not contain any
				// information about source. Technically, I feel like this would be best
				// solved with a "future feature" - expressions should contain a
				// reference to their source. Then, when the error is encountered here
				// or at runtime; it's easy to unwind the stack.
				if asIdent, isIdent := exprs[0].(*IdentValue); isIdent {
					switch asIdent.Val {
					case "if":
						return convertToIfExpr(exprs[1:])
					case "fn":
						return convertToFnExpr(exprs[1:])
					case "let":
						return convertToLetExpr(exprs[1:])
					}
				}
			}
			return NewCallExpr(exprs...), nil

		case OpenParenTT:
			subCall, subCallErr := parseCallExpr(ts)
			if subCallErr != nil {
				return nil, subCallErr
			}
			exprs = append(exprs, subCall)

		case IdentTT:
			identV, identErr := parseIdentValue(next)
			if identErr != nil {
				return nil, identErr
			}
			exprs = append(exprs, identV)

		case OpTT:
			opFn, opFnErr := parseOpValue(next)
			if opFnErr != nil {
				return nil, opFnErr
			}
			exprs = append(exprs, opFn)

		case NumberTT:
			numV, numErr := parseNumberValue(next)
			if numErr != nil {
				return nil, numErr
			}
			exprs = append(exprs, numV)

		case StringTT:
			strV, strErr := parseStringValue(next)
			if strErr != nil {
				return nil, strErr
			}
			exprs = append(exprs, strV)

		case CommentTT:
			// do nothing

		default:
			return nil, NewParseError("invalid token", next)
		}
	}

	// ques (bs): should this be considered a parse error? Right now, the
	// type is fairly inflexible: it's strictly bound to tokens. Should it
	// be expanded to include the possibility of *absence* of a token as the
	// source of an error?
	return nil, fmt.Errorf("encountered end mid-expression")
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

func convertToIfExpr(exprs []Expr) (*IfExpr, error) {
	if len(exprs) == 0 {
		return nil, fmt.Errorf("if requires a condition")
	}
	cond := exprs[0]
	var left, right Expr
	if len(exprs) > 1 {
		left = exprs[1]
	}
	if len(exprs) > 2 {
		right = exprs[2]
	}
	return NewIfExpr(cond, left, right), nil
}

func convertToFnExpr(exprs []Expr) (*FnExpr, error) {
	if len(exprs) == 0 {
		return nil, fmt.Errorf("fn requires an argument list")
	}
	args := []Arg{}
	asCall, isCall := exprs[0].(*CallExpr)
	if !isCall {
		return nil, fmt.Errorf("fn requires an argument list as the first expr")
	}
	for _, e := range asCall.Exprs {
		asIdent, isIdent := e.(*IdentValue)
		if !isIdent {
			return nil, fmt.Errorf("fn argument list must be all idents")
		}
		args = append(args, Arg{
			Ident: asIdent.Val,
		})
	}
	return NewFnExpr(args, exprs[1:]), nil
}

func convertToLetExpr(exprs []Expr) (Expr, error) {
	if len(exprs) != 2 {
		return nil, fmt.Errorf("let requires exactly two arguments")
	}
	asIdent, isIdent := exprs[0].(*IdentValue)
	if !isIdent {
		return nil, fmt.Errorf("let requires an ident as the first values")
	}
	valExpr := exprs[1]

	return &LetExpr{
		Ident: asIdent,
		Value: valExpr,
	}, nil
}
