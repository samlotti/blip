package internal

import (
	"fmt"
)

// TokenType == May be better as an integer!
type TokenType string

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Pos     int
}

const (
	ILLEGAL = "Illegal"
	EOF     = "Eof"
	EOL     = '\n'

	LITERAL         = "LITERAL"
	ARG             = "@arg"     // Literal will be the remainder of the line
	CONTEXT         = "@context" // Context variable expected
	ATDisplayInt    = "@int="    // write integer
	ATDisplayInt64  = "@int64="  // write integer
	ATDisplay       = "@="       // Literal will be up to the eol/eof or next @   @= name @
	ATDisplayUnsafe = "@=="      // Literal will be up to the eol/eof or next @   @= name @
	IMPORT          = "@import"  // Placed at the begging for go imports
	INCLUDE         = "@include" // includes another template but no embedded content
	EXTEND          = "@extend"  // includes another template
	CONTENT         = "@content" // The content to embed
	YIELD           = "@yield"   // provide content to the included template
	ENDBLOCK        = "@}"       // The end of block style commands
	STARTBLOCK      = "@{"       // Start of a code block
	FUNCTS          = "@func"    // functions
	IF              = "@if"      // The if statement convert to if <content> {
	ELSE            = "@else"    // converts to } else {
	ENDIF           = "@endif"   // converts to }
	END             = "@end"     // converts to }
	FOR             = "@for"     // convert for for range loop
	ENDFOR          = "@endfor"  // converts to }

)

type Lexer struct {
	input        string // The complete input
	runes        []rune
	FName        string
	lineNum      int  // The line number
	lPos         int  // Position of token on the line
	position     int  // the current character
	readPosition int  // The next position
	ch           rune // current character
	priorToken   Token
}

func NewLexer(input string, fname string) *Lexer {
	// input = input + "\n"
	l := &Lexer{input: input, FName: fname, lineNum: 1, lPos: 0, runes: []rune(input)}
	// log.Printf("Lexer: %s\n", input)
	l.readChar() // Prime the first character
	return l
}

// readChar
// places current character in l.ch
// advances the pointer to the next character
func (l *Lexer) readChar() {

	if l.isEOF() {
		return
	}

	if l.readPosition >= len(l.runes) {
		l.ch = 0
	} else {
		l.ch = l.runes[l.readPosition]
		if l.ch == '\n' {
			l.lineNum += 1
			l.lPos = 0
		}
	}
	l.position = l.readPosition
	l.lPos += 1
	l.readPosition += 1
}
func (l *Lexer) isEOL() bool {
	return l.ch == '\n'
}
func (l *Lexer) isEOF() bool {
	return l.readPosition > len(l.runes)
}
func (l *Lexer) peekCharAt(idx int) rune {
	if l.readPosition >= len(l.runes)-idx {
		return 0
	} else {
		return l.runes[l.readPosition+idx]
	}
}
func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.runes) {
		return 0
	} else {
		return l.runes[l.readPosition]
	}
}

//func (l *Lexer) peekCharNext() rune {
//	if l.readPosition >= len(l.runes)-1 {
//		return 0
//	} else {
//		return l.runes[l.readPosition+1]
//	}
//}

//func (l *Lexer) priorChar() rune {
//	if l.readPosition < 2 {
//		return 0
//	}
//	return l.runes[l.readPosition-2]
//}

func (l *Lexer) newToken(tokenType TokenType, ch rune) Token {
	var r = Token{Type: tokenType, Literal: string(ch), Line: l.lineNum, Pos: l.lPos}
	l.priorToken = r
	return r
}

//func (l *Lexer) newToken2(tokenType TokenType, ch rune, ch2 rune) Token {
//	var r = Token{Type: tokenType, Literal: string(ch) + string(ch2), Line: l.lineNum, Pos: l.lPos}
//	l.priorToken = r
//	return r
//}
func (l *Lexer) newTokenStr(tokenType TokenType, ch string) Token {
	var r = Token{Type: tokenType, Literal: ch, Line: l.lineNum, Pos: l.lPos}
	l.priorToken = r
	return r
}

func (l *Lexer) PriorToken() *Token {
	return &l.priorToken
}

// NextToken
// On input  l.ch is the current character
// On output should be at next character after token required characters.
// peekChar will be looking at the next character
func (l *Lexer) NextToken() (tk *Token) {

	defer func() {
		if err := recover(); err != nil {
			tk1 := l.newTokenStr(ILLEGAL, fmt.Sprintf("Error parsing: %s", err))
			tk = &tk1
		}
	}()

	var tok Token

	// if l.skipWhitespace() {
	//	// Add a semi
	//	return l.newToken(SEMI, ';')
	//}

	if l.isEOF() {
		tok := l.newToken(EOF, '0')
		return &tok
	}

	switch l.ch {
	case '@':
		// @//
		if l.peekChar() == '/' && l.peekCharAt(1) == '/' {
			l.readTil('\n')
			l.readChar()
			return l.NextToken()
		}

		if l.peekChar() == '*' {
			l.bypassMultilineComment()
			l.readChar()
			return l.NextToken()
		}

		if l.peekChar() == '@' {
			tok = l.newToken(LITERAL, '@')
			l.readChar()
			l.readChar()
		} else {
			tok = l.pickCommand()
		}
		// This is the special character
		break
	default:
		tok = l.newTokenStr(LITERAL, l.readLiteral())
	}
	// Should always be on the next character

	return &tok

}

func (l *Lexer) readLiteral() string {
	pos := l.position
	for !l.isAt() && !l.isEOF() {
		l.readChar()
	}
	l.readChar()
	return string(l.runes[pos:l.position])
}

func (l *Lexer) isChar(char rune) bool {
	if l.peekChar() == char {
		return true
	}
	return false
}
func (l *Lexer) isAnyChar(chars []rune) bool {
	pc := l.peekChar()
	for _, ch := range chars {
		if ch == pc {
			return true
		}
	}
	return false
}

func (l *Lexer) isStr(chars []rune) bool {
	for idx, ch := range chars {
		if l.peekCharAt(idx) != ch {
			return false
		}
	}
	return true
}

func (l *Lexer) isAt() bool {
	return l.isChar('@')
}

func (l *Lexer) readTilStr(str []rune) string {
	pos := l.position
	for !l.isStr(str) && !l.isEOF() {
		l.readChar()
	}

	if l.isEOF() {
		panic(fmt.Sprintf("Expected to find '%s' not found.", string(str)))
	}

	posend := l.position
	for range str {
		l.readChar()
	}
	return string(l.runes[pos : posend+1])
}

func (l *Lexer) readTilStrSingleLine(str []rune) string {
	pos := l.position
	for !l.isStr(str) && !l.isEOF() && !l.isEOL() {
		l.readChar()
	}

	if l.isEOL() {
		panic(fmt.Sprintf("Expected to find '%s' on same line, not found.", string(str)))
	}

	if l.isEOF() {
		panic(fmt.Sprintf("Expected to find '%s' not found.", string(str)))
	}

	posend := l.position
	for range str {
		l.readChar()
	}
	return string(l.runes[pos : posend+1])
}

func (l *Lexer) readTil(char rune) string {
	pos := l.position
	for !l.isChar(char) && !l.isEOF() {
		l.readChar()
	}
	l.readChar()
	return string(l.runes[pos:l.position])
}

func (l *Lexer) readTils(chars []rune) string {
	pos := l.position
	for !l.isAnyChar(chars) {
		l.readChar()
		if l.isEOF() {
			break
		}
	}

	if l.isEOF() {
		panic(fmt.Sprintf("Expected any of '%s' not found.", string(chars)))
	}

	l.readChar()
	return string(l.runes[pos:l.position])
}

// On input, we have '@'
func (l *Lexer) pickCommand() Token {

	if l.peekChar() == '}' {
		// This is an end block ... not really a command! and doesn't need a space
		tk := l.newTokenStr(ENDBLOCK, "@}")
		l.readChar()
		l.readChar()
		return tk
	}

	cmd := l.readTils([]rune{' ', '\n'})

	if l.ch != '\n' {
		l.readChar()
	}

	advance := false

	var tkn Token
	switch cmd {
	case "@yield":
		tkn = l.newTokenStr(YIELD, l.readTil(EOL))
		advance = true
	case "@if":
		tkn = l.newTokenStr(IF, l.readTil(EOL))
		advance = true
	case "@for":
		tkn = l.newTokenStr(FOR, l.readTil(EOL))
		advance = true
	case "@content":
		tkn = l.newTokenStr(CONTENT, l.readTil(EOL))
		advance = true
	case "@include":
		// tkn = l.getIncludeToken()
		tkn = l.newTokenStr(INCLUDE, l.readTil(EOL))
		advance = true
	case "@extend":
		//tkn = l.getExtendToken()
		tkn = l.newTokenStr(EXTEND, l.readTil(EOL))
		//if strings.Contains(tkn.Literal, "@") {
		//	tkn = l.newTokenStr(ILLEGAL, fmt.Sprintf("@ invalid in the : %s", tkn.Type))
		//}
		advance = true
	case "@int=":
		tkn = l.newTokenStr(ATDisplayInt, l.readTilStrSingleLine([]rune{'@'}))
		advance = true
	case "@int64=":
		tkn = l.newTokenStr(ATDisplayInt64, l.readTilStrSingleLine([]rune{'@'}))
		advance = true
	case "@==":
		// tkn = l.newTokenStr(ATDisplayUnsafe, l.readTils([]rune{EOL, '@'}))
		tkn = l.newTokenStr(ATDisplayUnsafe, l.readTilStrSingleLine([]rune{'@'}))
		advance = true
	case "@=":
		// tkn = l.newTokenStr(ATDisplay, l.readTils([]rune{EOL, '@'}))
		tkn = l.newTokenStr(ATDisplay, l.readTilStrSingleLine([]rune{'@'}))
		advance = true
		// advance because we consumed more
	case "@arg":
		tkn = l.newTokenStr(ARG, l.readTil(EOL))
		advance = true
		// consume the eol
	case "@context":
		tkn = l.newTokenStr(CONTEXT, l.readTil(EOL))
		// consume the eol
		advance = true
	case "@else":
		tkn = l.newTokenStr(ELSE, "")
		advance = false
	case "@endif":
		tkn = l.newTokenStr(ENDIF, "")
		advance = false
	case "@end":
		tkn = l.newTokenStr(END, "")
		advance = false
	case "@endfor":
		tkn = l.newTokenStr(ENDFOR, "")
		advance = false
	case "@import":
		tkn = l.newTokenStr(IMPORT, l.readTil(EOL))
		// consume the eol
		advance = true
	case "@func":
		tkn = l.newTokenStr(FUNCTS, cmd)
		advance = false
	case "@{":
		tkn = l.newTokenStr(STARTBLOCK, cmd)
		advance = false
	case "@}":
		tkn = l.newTokenStr(ENDBLOCK, cmd)
		advance = false
	default:
		tkn = l.newTokenStr(ILLEGAL, fmt.Sprintf("Invalid command found: %s", cmd))
		advance = false
	}

	if advance {
		l.readChar()
	}

	return tkn
}

// @* ... *@
func (l *Lexer) bypassMultilineComment() {
	l.readChar() // Get past the *
	for !l.isEOF() {
		if l.ch == '*' {
			if l.peekChar() == '@' {
				l.readChar()
				return
			}
		}
		l.readChar()
	}
}

//
//func (l *Lexer) getExtendToken() Token {
//	content := l.readTilStr([]rune("{@"))
//	return l.newTokenStr(EXTEND, content)
//}
//
//func (l *Lexer) getIncludeToken() Token {
//	content := l.readTils([]rune{'@', '{'})
//	t := INCLUDE
//	if l.ch == '{' {
//		if l.peekChar() != '@' {
//			return l.newTokenStr(ILLEGAL, "Expected @ after the {")
//		}
//		l.readChar()
//		t = EXTEND
//	} else {
//		if l.isEOF() {
//			return l.newTokenStr(ILLEGAL, "Expected @")
//		}
//	}
//
//	l.readChar()
//
//	return l.newTokenStr(TokenType(t), content)
//}
