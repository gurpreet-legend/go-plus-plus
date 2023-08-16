package main

import "unicode"

type tokenKind uint

const (
	// eg. "(", ")"
	syntaxToken tokenKind = iota
	// eg. "3", "12"
	integerToken
	// eg. "+", "define"
	identifierToken
)

type token struct {
	value string
	kind  tokenKind
	location int
}

func (t token) debug(source []rune) {
	// implement where have we gone through the source and give 
	// a customized error based on location of token, for err debugging
}

func removeWhitespace(source []rune, cursor int) int {
	for cursor < len(source) {
		if unicode.IsSpace(source[cursor]) {
			cursor++
			continue
		}
		break
	}
	return cursor
}


func lexSyntaxToken(source []rune, cursor int) (int, *token) {
	if source[cursor] == '(' || source[cursor] == ')' {
		return cursor + 1, &token{
			value: string([]rune{source[cursor]}),
			kind: syntaxToken,
			location: cursor,
		}
	}

	return cursor, nil
}

// lexIntegerToken("foo 123", 4) => "123"
// lexIntegerToken("foo 12 3", 4) => "12"
// lexIntegerToken("foo 12a 3", 4) => "12" <--- Ignoring this situation for now
func lexIntegerToken(source []rune, cursor int) (int, *token) {
	originalCursor := cursor
	var value []rune
	
	for cursor < len(source) {
		r := source[cursor]
		if r >= '0' && r <= '9' {
			value = append(value, r)
			cursor++
			continue
		}	
		break
	}

	if len(value) == 0 {
		return originalCursor, nil
	}
	
	return cursor, &token{
		value: string(value),
		kind: integerToken,
		location: originalCursor,
	}
}

//lexIdentifier("123 ab + ", 4) => "ab"
//lexIdentifier("123 ab123 + ", 4) => "ab123"
func lexIdentifierToken(source []rune, cursor int) (int, *token) {
	originalCursor := cursor
	var value []rune
	
	for cursor <len(source){
		r := source[cursor]
		if !unicode.IsSpace(r) {
			value = append(value, r)
			cursor++
			continue
		}

		break
	}

	if len(value) == 0 {
		return originalCursor, nil
	}

	return cursor, &token{
		value: string(value),
		kind: identifierToken,
		location: originalCursor,
	}
}

// eg. "( + 12 8)"
// lex("( + 12 8)") should give : ["(", "+", "12", "8", ")"]
func lex(raw string) []token {
	source := []rune(raw)
	var tokens []token
	var t *token

	cursor := 0
	for {
		//remove whitespaces
		cursor = removeWhitespace(source, cursor)
		if cursor == len(source) {
			break
		}

		cursor, t = lexSyntaxToken(source, cursor)
		if t != nil {
			tokens = append(tokens, *t)
			continue
		}

		cursor, t = lexIntegerToken(source, cursor)
		if t != nil {
			tokens = append(tokens, *t)
			continue
		}

		cursor, t = lexIdentifierToken(source, cursor)
		if t != nil {
			tokens = append(tokens, *t)
			continue
		}

		//nothing to lex, panic!
		// fmt.Println(tokens[len(tokens)-1].debug()) // err on line of code
		panic("Could not lex")

		//check for syntaxToken: if yes: `continue`

		// or check for integerToken: if yes: `continue`
	}

	return tokens
}