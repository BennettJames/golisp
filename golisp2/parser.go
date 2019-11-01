package golisp2

import (
	"fmt"
	"strconv"
	"strings"
)

// ExecString executes the lisp program contained in str, and returns the
// output.
func ExecString(str string) (string, error) {
	ts := NewTokenScanner(NewRuneScanner(strings.NewReader(
		`((fn (x) (+ x x)) 5)`)))
	exprs, exprsErr := ParseTokens(ts)
	if exprsErr != nil {
		return "", exprsErr
	}
	var sb strings.Builder

	// todo (bs): seed this with built-ins. Preferrably; there'd be a way to mark
	// contexts as non-extendable and could have a shared built-in, but that's
	// very hypothetical.
	c := &ExprContext{
		vals: map[string]Value{},
	}
	for _, e := range exprs {
		v := e.Eval(c)
		sb.WriteString(v.PrintStr())
		sb.WriteByte('\n')
	}

	return sb.String(), nil
}

// ParseTokens reads in the tokens, and converts them to a set of expressions.
// Returns the set, and any parse errors that are encountered in the process.
func ParseTokens(ts *TokenScanner) ([]Expr, error) {

	// note (bs): I think the token scanner will likely need to be modified so it
	// retains a value in the same way that the rune scanner does. Having that
	// retention just makes it much easier to iterate through the values.
	//
	// regardless: the idea here is that I just parse and get an open paren, then
	// call a subparser to finish. If anything else is encountered, then an error
	// is generated.
	//
	// Not sure I should worry about this quite yet, but one thing I'll mention:
	// it'd probably be possible to retain token location in source as part of
	// this, and generate error messages off of that. Again, not sure one way or
	// another if that's worth the effort yet. I'd at least give it a gander after
	// the main body of work here is done. It might not be too bad; and it
	// wouldn't hurt to at least start making half-hearted attempts at
	// end-usability here.
	//
	// One thing I would like to do is make it so I could properly print out the
	// code from the AST. It feels wrong to me to throw that info out; I feel like
	// it wouldn't be unreasonable to say that I should have some means to print
	// out the entirety of the tree to something that works. I think that's a
	// fundamental part of homoiconicity.

	exprs := []Expr{}
	for !ts.Done() {
		// note (bs): token interchange here is super hacky; needs to be better
		maybeOpen := ts.Next()
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
					switch asIdent.ident {
					case "if":
						return convertToIfExpr(exprs[1:])
					case "fn":
						return convertToFnExpr(exprs[1:])
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
			// there are a few special cases here. For starters at least, if/fn I
			// believe need special handling. *technically* I could probably get by
			// deferring that, but I think doing it here makes sense.
			//
			// They are admittedly a little weird. Particularly: they can only be the
			// *first* expression. After that, they should be considered illegal.
			//
			// That actually has some implications. You need to do an initial check to
			// see if the ident has special handling
			//
			// The alternative, it's worth reminding, would be to treat everything as
			// vanilla expression, then have special cases after the fact. Perhaps I
			// should give that more consideration here; but I just don't like it
			// much.
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
		return NewFuncValue(addFn), nil
	case "-":
		return NewFuncValue(subFn), nil
	case "*":
		return NewFuncValue(multFn), nil
	case "/":
		return NewFuncValue(divFn), nil
	case "==":
		return NewFuncValue(eqNumFn), nil
	case "<":
		return NewFuncValue(ltNumFn), nil
	case ">":
		return NewFuncValue(gtNumFn), nil
	case "<=":
		return NewFuncValue(lteNumFn), nil
	case ">=":
		return NewFuncValue(gteNumFn), nil
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
	// so - how will this work? This is more complicated than the if conversion.
	// Technically I sort
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
