package golisp2

import (
	"fmt"
	"strconv"
	"strings"
)

// ExecString executes the lisp program contained in str, and returns the
// output.
func ExecString(str string) (string, error) {
	ts := NewTokenScanner(NewRuneScanner(strings.NewReader(str)))
	exprs, exprsErr := ParseTokens(ts)
	if exprsErr != nil {
		return "", exprsErr
	}

	c := BuiltinContext()

	var sb strings.Builder
	for _, e := range exprs {
		v := e.Eval(c)
		sb.WriteString(v.InspectStr())
		sb.WriteByte('\n')
	}

	return sb.String(), nil
}

// ParseTokens reads in the tokens, and converts them to a set of expressions.
// Returns the set, and any parse errors that are encountered in the process.
func ParseTokens(ts *TokenScanner) ([]Expr, error) {
	exprs := []Expr{}
	for !ts.Done() {
		// note (bs): token interchange here is super hacky; needs to be better
		maybeOpen := ts.Next()

		if maybeOpen == nil {
			if ts.Done() {
				break
			} else {
				return nil, fmt.Errorf("unexpected nil token")
			}
		}

		if maybeOpen.Typ != OpenParenTT {
			return nil, fmt.Errorf("unexpected top-level token: %+v", maybeOpen)
		}
		expr, exprErr := parseCallExpr(ts)
		if exprErr != nil {
			return nil, exprErr
		}
		exprs = append(exprs, expr)
	}
	return exprs, nil
}

func parseCallExpr(ts *TokenScanner) (Expr, error) {

	exprs := []Expr{}
	for !ts.Done() {
		// yeah, so this won't work as-is. I could do a hacky fix: hoist the search
		// for open parens to the parent. Then, this for block also can look for
		// open parens. Not great, but not terrible.
		next := ts.Next()

		switch next.Typ {
		case CloseParenTT:
			if len(exprs) > 0 {
				// note (bs): this is awfully clumsy. I think it'd be better to actually
				// send out a subparser to be responsible for handling built-ins/macros,
				// but I'll worry about that later.
				if asIdent, isIdent := exprs[0].(*IdentValue); isIdent {
					switch asIdent.Str {
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
			identV, identErr := parseIdentValue(next.Value)
			if identErr != nil {
				return nil, identErr
			}
			exprs = append(exprs, identV)

		case OpTT:
			opFn, opFnErr := parseOpValue(next.Value)
			if opFnErr != nil {
				return nil, opFnErr
			}
			exprs = append(exprs, opFn)

		case NumberTT:
			numV, numErr := parseNumberValue(next.Value)
			if numErr != nil {
				return nil, numErr
			}
			exprs = append(exprs, numV)

		case StringTT:
			strV, strErr := parseStringValue(next.Value)
			if strErr != nil {
				return nil, strErr
			}
			exprs = append(exprs, strV)

		case CommentTT:
			// do nothing

		default:
			return nil, fmt.Errorf("unparsable token found: %+v", next)
		}
	}
	return nil, fmt.Errorf("encountered end mid-expression")
}

func parseStringValue(buf string) (*StringValue, error) {
	if len(buf) == 0 {
		return NewStringValue(""), nil
	}
	leadI, tailI := 0, len(buf)
	if buf[0] == '"' {
		leadI = 1
	}
	if len(buf) > 1 && buf[len(buf)-1] == '"' {
		tailI = len(buf) - 1
	}
	return NewStringValue(buf[leadI:tailI]), nil
}

func parseIdentValue(buf string) (Value, error) {
	// todo (bs): this should search for certain reserved words, and reject them.

	switch buf {
	case "nil":
		return NewNilValue(), nil
	case "true":
		return NewBoolValue(true), nil
	case "false":
		return NewBoolValue(false), nil
	default:
		return NewIdentValue(buf), nil
	}
}

func parseNumberValue(buf string) (*NumberValue, error) {
	f, e := strconv.ParseFloat(buf, 64)
	if e != nil {
		return nil, e
	}
	return NewNumberValue(f), nil
}

func parseOpValue(rawV string) (*FuncValue, error) {
	// todo (bs): strongly consider moving this to a map rather than a case
	// statement
	switch rawV {
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
		return nil, fmt.Errorf("unrecognized operator: %s", rawV)
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
	for _, c := range asCall.Get() {
		asIdent, isIdent := c.(*IdentValue)
		if !isIdent {
			return nil, fmt.Errorf("fn argument list must be all idents")
		}
		args = append(args, Arg{
			Ident: asIdent.Get(),
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

	// So - I *think* this would work, but I disagree with it on principle. First:
	// let's see if it works; then I can navel gaze once again on the relative
	// sensibility of this.
	//
	// Alright, so it works. What next? I think, given the nature of the
	// language-style I am aiming for, I'd like to still convert this to a data
	// structure. But I'd like to go farther than that: I'd like to convert all
	// AST elements to raw-ish structs. That means no purposeless accessors, and
	// only to have constructors in cases where it makes it easier to reason
	// about.
	//
	// So - for the let case, I'd say make it just have an ident and an assignment
	// expression. Can have a try-constructor like this that takes a list of
	// arguments and assigns them appropriately.
	//
	// This does raise a separate question for me though: should I simplify my
	// core set of function types? The following feels a little unnecessarily
	// complex. Granted, I am sorta explicitly abusing the built-in's here, but
	// this doesn't feel like a hard case. Is it possible my built-in's are
	// groping and are not finding the right level of abstraction?
	//
	// So, first up: I think it's fair to say there are two different notions of
	// evaluation, of sorts. There's the external "compute and return the
	// expression", and the internal notion of "evaluate all arguments, the
	// compute the base function". Are those truly different though; or is that
	// just one concept with two sides? Plain eval'ing can kinda be crammed into
	// the exec category just by ignoring the expr list (as this case does), but
	// that doesn't feel satisfying.
	//
	// Options? One is just to preserve the different notions. I still could
	// smooth the paths somewhat if I have some more obvious, core API's; could
	// even do things like make boxing a plain eval as an exec easy (for better or
	// worse).
	//
	// Alternatively: would it make sense to re-write exec's as eval's, so to
	// speak? That is:
	//
	// That wouldn't really seem to be doing anything though. Maybe it just comes
	// down to a fundamental divide: blocks are evaluated; functions are called.
	// The latter by it's nature has it's
	//
	// The problem here then is just the lack of a primitive for self-contained
	// expression blocks. Arguably, function declarations themselves are also a
	// little more convoluted then they need to be. I don't think it's wrong per
	// se to have a wrapper for them; but the wrapper is underpowered and I don't
	// think it necessarily should be mandatory.
	//
	// So - should I change that? I could define types like this:
	//
	//  type PlainFn func(*ExprContext, []Expr) (Value, Error)
	//  type PlainExpr func(*ExprContext) (Value, Error)
	//
	// then define "Call" and "Eval" on them, respectively.

	return &LetExpr{
		Ident: asIdent,
		Value: valExpr,
	}, nil
}
