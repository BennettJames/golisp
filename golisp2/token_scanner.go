package golisp2

import (
	"unicode"
)

// note (bs): probably will eventually want to turn this into a token scanner,
// rather than a monolithic parser.

// TokenizeString converts the provided string to a list of tokens.
func TokenizeString(str string) []ScannedToken {
	tokens := []ScannedToken{}

	cs := NewCharScanner(str)
	for !cs.Done() {
		next := nextToken(cs)
		if next != nil {
			tokens = append(tokens, *next)
		}
	}
	return tokens
}

func nextToken(cs *CharScanner) *ScannedToken {
	for {
		for peekIsSpace(cs) {
			// note (bs): consider still supporting whitespace as an explicit token
			// type. For the most part it doesn't really matter, but arguably newlines
			// in particular are important for "reconstituting" a function.
			cs.Next()
		}
		peekC, peekOk := cs.Peek()
		if !peekOk {
			return nil
		}

		// note (bs): again; don't think this is good design; just a place to start.
		switch peekC {
		case '(':
			return parseOpenParen(cs)
		case ')':
			return parseClosedParen(cs)
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			return parseNumberToken(cs)
		case '"':
			return parseStringToken(cs)
		case '-', '+', '/', '*', '&', '^', '%', '!', '|', '<', '>', '=':
			return parseOpToken(cs)

		default:
			// let's just lazily scan in an ident if it's anything else. That's not
			// even strictly correct; just trying to hustle along for a bit
			return parseIdentToken(cs)
		}
	}
}

func parseOpenParen(cs *CharScanner) *ScannedToken {
	nextC, nextOk := cs.Next()
	if !nextOk || nextC != '(' {
		return nil
	}
	return &ScannedToken{
		Typ:   OpenParenTT,
		Value: "(",
	}
}

func parseClosedParen(cs *CharScanner) *ScannedToken {
	nextC, nextOk := cs.Next()
	if !nextOk || nextC != ')' {
		return nil
	}
	return &ScannedToken{
		Typ:   CloseParenTT,
		Value: ")",
	}
}

// a note on parsers here: they all can return invalid tokens or nil. Not sure
// about returning nil here, to be honest - seems like it only happens if the
// scan unexpectedly stops, but that in of itself would be... wrong. Oh well.

func parseNumberToken(cs *CharScanner) *ScannedToken {
	numStr := ""
	for {
		peekC, peekOk := cs.Peek()
		if !peekOk || runeIsBoundary(peekC) {
			break
		}
		if !unicode.IsDigit(peekC) && peekC != '.' {
			return &ScannedToken{
				Typ:   InvalidTT,
				Value: numStr,
			}
		}
		numStr += string(peekC)
		cs.Next()
	}

	// note (bs): this would allow some numbers like "123." through. I don't think
	// I want to allow that. Should it be forbidden here? Maybe. I'd say let it go
	// for now, but I'd like to review some samples and come back to this whole
	// section and likely make the underlying scanner a bit better and make it all
	// a little more state-machine-y. I think the right way to view that is by
	// having a more stepped approach to the parse logic, wherein certain
	// "terminal" characters will return invalid if they're the final parse rather
	// than
	//
	// yeah, let's try to do that. This parsing o.k. for getting started, but it's
	// pretty clumsy. Id' also rather not coast on a post-hoc check w/ golang for
	// this one; outsourcing my parsing feels like cheating.

	if len(numStr) == 0 {
		return nil
	}

	return &ScannedToken{
		Typ:   NumberTT,
		Value: numStr,
	}
}

func parseOpToken(cs *CharScanner) *ScannedToken {
	// note (bs): I sorta feel like "batching" the read and ok for the scanner
	// ended up not working in my favor. Might make more sense to have the more
	// traditional model wherein the scanner has a simple "done" value that I can
	// read, and I can call advance (no return), next (one head), and peek (two
	// ahead) individually. Oh well; I'd say still scratch together some basic
	// token parsing here even if it's trash; then go back and make "CharScanner2"
	// if you are so inclined.
	//
	// Yeah, let's do something here. The current design feels like a very awkward
	// imposition of the for loop; I don't think I'll have to struggle to long to
	// get something that fits a little more naturally into Go iteration.

	// note (bs): it's technically possible for this to return a number type in
	// the case of '-'
	opStr := "" // note (bs): would be marginally better if this was a string builder
	for {
		peekC, peekOk := cs.Peek()
		if !peekOk || runeIsBoundary(peekC) {
			break
		}

		if unicode.IsDigit(peekC) {
			if opStr == "-" {
				numToken := parseNumberToken(cs)
				if numToken != nil {
					return &ScannedToken{
						Typ:   numToken.Typ,
						Value: opStr + numToken.Value,
					}
				}
			}
			return &ScannedToken{
				Typ:   InvalidTT,
				Value: opStr + string(peekC),
			}
		}
		if !runeIsOp(peekC) {
			return &ScannedToken{
				Typ:   InvalidTT,
				Value: opStr + string(peekC),
			}
		}
		opStr += string(peekC)

		cs.Next()
	}

	if len(opStr) == 0 {
		// the existence of this, among other things, makes me kinda doubt whether
		// my design here is any good
		return nil
	}

	return &ScannedToken{
		Typ:   OpTT,
		Value: opStr,
	}
}

func parseStringToken(cs *CharScanner) *ScannedToken {
	// fixme
	nextC, nextOk := cs.Next()
	var _, _ = nextC, nextOk

	return &ScannedToken{
		Typ: InvalidTT,
	}
}

func parseIdentToken(cs *CharScanner) *ScannedToken {
	nextC, nextOk := cs.Next()
	var _, _ = nextC, nextOk

	return &ScannedToken{
		Typ: InvalidTT,
	}
}

func peekIsBoundary(cs *CharScanner) bool {
	peekC, peekOk := cs.Peek()
	return !peekOk || runeIsBoundary(peekC)
}

func peekIsSpace(cs *CharScanner) bool {
	peekC, peekOk := cs.Peek()
	return peekOk && unicode.IsSpace(peekC)
}

func runeIsBoundary(r rune) bool {
	return unicode.IsSpace(r) ||
		r == '(' ||
		r == ')'
}

func runeIsOp(r rune) bool {
	switch r {
	case '-', '+', '/', '*', '&', '^', '%', '!', '|', '<', '>', '=':
		return true
	default:
		return false
	}
}
