package internal

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLexer1(t *testing.T) {
	sample := `This is a simple string with
no replacements!!
`
	lex := NewLexer(sample, "TestLexer1")
	tkn := lex.NextToken()
	assert.Equal(t, LITERAL, string(tkn.Type))
	assert.Equal(t, sample, tkn.Literal)
	tkn = lex.NextToken()
	assert.Equal(t, EOF, string(tkn.Type))
	// assert.Equal(t, "t", tkn.Literal)
}

func TestLexer2(t *testing.T) {
	sample := `yikes @@, its good`

	lex := NewLexer(sample, "TestLexer1")
	tkn := lex.NextToken()
	assert.Equal(t, LITERAL, string(tkn.Type))
	assert.Equal(t, "yikes ", tkn.Literal)

	assert.Equal(t, "@", lex.NextToken().Literal)
	assert.Equal(t, LITERAL, string(lex.PriorToken().Type))

	assert.Equal(t, ", its good", lex.NextToken().Literal)
	assert.Equal(t, LITERAL, string(lex.PriorToken().Type))

	tkn = lex.NextToken()
	assert.Equal(t, EOF, string(tkn.Type))

}

func TestHello(t *testing.T) {
	sample := `@arg name string, amount int64
@context user models.User
@badCommand
<h1>@= name@</h1>
@== unsafeExp("45") @`
	lex := NewLexer(string(sample), "TestLexer1")
	tkn := lex.NextToken()
	assert.Equal(t, ARG, string(tkn.Type))
	assert.Equal(t, "name string, amount int64", string(tkn.Literal))

	tkn = lex.NextToken()
	assert.Equal(t, CONTEXT, string(tkn.Type))
	assert.Equal(t, "user models.User", string(tkn.Literal))

	tkn = lex.NextToken()
	assert.Equal(t, ILLEGAL, string(tkn.Type))
	assert.Equal(t, "Invalid command found: @badCommand", string(tkn.Literal))

	tkn = lex.NextToken()
	assert.Equal(t, LITERAL, string(tkn.Type))
	assert.Equal(t, "\n<h1>", string(tkn.Literal))

	tkn = lex.NextToken()
	assert.Equal(t, ATDisplay, string(tkn.Type))
	assert.Equal(t, "name", string(tkn.Literal))

	tkn = lex.NextToken()
	assert.Equal(t, LITERAL, string(tkn.Type))
	assert.Equal(t, "</h1>\n", string(tkn.Literal))

	tkn = lex.NextToken()
	assert.Equal(t, ATDisplayUnsafe, string(tkn.Type))
	assert.Equal(t, "unsafeExp(\"45\") ", string(tkn.Literal))

	tkn = lex.NextToken()
	assert.Equal(t, EOF, string(tkn.Type))

}

func TestComment1(t *testing.T) {
	sample := `@// This is a comment
After comment   @// another comment
After2
@// Event at end
@// And multiple
Last
`
	lex := NewLexer(string(sample), "TestLexer1")
	tkn := lex.NextToken()
	assert.Equal(t, LITERAL, string(tkn.Type))
	assert.Equal(t, "After comment   ", string(tkn.Literal))

	tkn = lex.NextToken()
	assert.Equal(t, LITERAL, string(tkn.Type))
	assert.Equal(t, "After2\n", string(tkn.Literal))

	tkn = lex.NextToken()
	assert.Equal(t, LITERAL, string(tkn.Type))
	assert.Equal(t, "Last\n", string(tkn.Literal))

	tkn = lex.NextToken()
	assert.Equal(t, EOF, string(tkn.Type))

}

func TestComment2(t *testing.T) {
	sample := `@* This is a comment
that's multi line! 
*@After comment   @// another comment
After2
@* Event at end
And multiple
*@
Last
`
	lex := NewLexer(string(sample), "TestLexer1")
	tkn := lex.NextToken()
	assert.Equal(t, LITERAL, string(tkn.Type))
	assert.Equal(t, "After comment   ", string(tkn.Literal))

	tkn = lex.NextToken()
	assert.Equal(t, LITERAL, string(tkn.Type))
	assert.Equal(t, "After2\n", string(tkn.Literal))

	tkn = lex.NextToken()
	assert.Equal(t, LITERAL, string(tkn.Type))
	assert.Equal(t, "\nLast\n", string(tkn.Literal))

	tkn = lex.NextToken()
	assert.Equal(t, EOF, string(tkn.Type))

}

func TestIncludeSimple(t *testing.T) {
	sample := `@include head 
After2
@include tail `
	lex := NewLexer(string(sample), "TestLexer1")
	tkn := lex.NextToken()
	assert.Equal(t, INCLUDE, string(tkn.Type))
	assert.Equal(t, "head ", string(tkn.Literal))

	tkn = lex.NextToken()
	assert.Equal(t, LITERAL, string(tkn.Type))
	assert.Equal(t, "After2\n", string(tkn.Literal))

	tkn = lex.NextToken()
	assert.Equal(t, INCLUDE, string(tkn.Type))
	assert.Equal(t, "tail ", string(tkn.Literal))

	tkn = lex.NextToken()
	assert.Equal(t, EOF, string(tkn.Type))

}

func TestInclude(t *testing.T) {
	sample := `@extend head args
@content myContent {@
After2
@extend tail `
	lex := NewLexer(string(sample), "TestLexer1")
	tkn := lex.NextToken()
	assert.Equal(t, EXTEND, string(tkn.Type))
	assert.Equal(t, "head args", string(tkn.Literal))

	//tkn = lex.NextToken()
	//assert.Equal(t, LITERAL, string(tkn.Type))
	//assert.Equal(t, "\n", string(tkn.Literal))

	// assert.Equal(t, "\n", lex.NextToken().Literal)

	tkn = lex.NextToken()
	assert.Equal(t, CONTENT, string(tkn.Type))
	assert.Equal(t, "myContent ", string(tkn.Literal))

	tkn = lex.NextToken()
	assert.Equal(t, LITERAL, string(tkn.Type))
	assert.Equal(t, "\nAfter2\n", string(tkn.Literal))

	tkn = lex.NextToken()
	assert.Equal(t, EXTEND, string(tkn.Type))
	assert.Equal(t, "tail ", string(tkn.Literal))

	tkn = lex.NextToken()
	assert.Equal(t, EOF, string(tkn.Type))

}

//
//func TestIncludeError(t *testing.T) {
//	sample := `@include head {
//After2
//@extend tail {@`
//	lex := NewLexer(string(sample), "TestLexer1")
//	tkn := lex.NextToken()
//	assert.Equal(t, ILLEGAL, string(tkn.Type))
//	assert.Equal(t, "Expected @ after the {", string(tkn.Literal))
//
//	tkn = lex.NextToken()
//	assert.Equal(t, LITERAL, string(tkn.Type))
//	assert.Equal(t, "\nAfter2\n", string(tkn.Literal))
//
//	tkn = lex.NextToken()
//	assert.Equal(t, EXTEND, string(tkn.Type))
//	assert.Equal(t, "tail ", string(tkn.Literal))
//
//	tkn = lex.NextToken()
//	assert.Equal(t, EOF, string(tkn.Type))
//
//}

//func TestIncludeError2(t *testing.T) {
//	sample := `@include head
//`
//	lex := NewLexer(string(sample), "TestLexer1")
//	tkn := lex.NextToken()
//	assert.Equal(t, ILLEGAL, string(tkn.Type))
//	assert.Equal(t, "Expected @", string(tkn.Literal))
//
//	tkn = lex.NextToken()
//	assert.Equal(t, EOF, string(tkn.Type))
//
//}

func TestYield(t *testing.T) {
	sample := `@yield head 
After2
@}`
	lex := NewLexer(string(sample), "TestLexer1")
	tkn := lex.NextToken()
	assert.Equal(t, YIELD, string(tkn.Type))
	assert.Equal(t, "head ", string(tkn.Literal))

	tkn = lex.NextToken()
	assert.Equal(t, LITERAL, string(tkn.Type))
	assert.Equal(t, "After2\n", string(tkn.Literal))

	tkn = lex.NextToken()
	assert.Equal(t, ENDBLOCK, string(tkn.Type))
	assert.Equal(t, "@}", string(tkn.Literal))

	tkn = lex.NextToken()
	assert.Equal(t, EOF, string(tkn.Type))

}

func TestCodeBlock(t *testing.T) {
	sample := `@{ if l.State == "test" { @}
    @{ } @}
`
	lex := NewLexer(string(sample), "TestLexer1")
	tkn := lex.NextToken()
	assert.Equal(t, STARTBLOCK, string(tkn.Type))
	assert.Equal(t, "@{", string(tkn.Literal))

	tkn = lex.NextToken()
	assert.Equal(t, LITERAL, string(tkn.Type))
	assert.Equal(t, "if l.State == \"test\" { ", string(tkn.Literal))

	tkn = lex.NextToken()
	assert.Equal(t, ENDBLOCK, string(tkn.Type))
	assert.Equal(t, "@}", string(tkn.Literal))

	tkn = lex.NextToken()
	assert.Equal(t, LITERAL, string(tkn.Type))
	assert.Equal(t, "\n    ", string(tkn.Literal))

	tkn = lex.NextToken()
	assert.Equal(t, STARTBLOCK, string(tkn.Type))
	assert.Equal(t, "@{", string(tkn.Literal))

	tkn = lex.NextToken()
	assert.Equal(t, LITERAL, string(tkn.Type))
	assert.Equal(t, "} ", string(tkn.Literal))

	tkn = lex.NextToken()
	assert.Equal(t, ENDBLOCK, string(tkn.Type))
	assert.Equal(t, "@}", string(tkn.Literal))

	assert.Equal(t, "\n", lex.NextToken().Literal)

	tkn = lex.NextToken()
	assert.Equal(t, EOF, string(tkn.Type))

}

func TestEndCodeBlock(t *testing.T) {
	sample := `@}`
	lex := NewLexer(string(sample), "TestLexer1")
	tkn := lex.NextToken()
	assert.Equal(t, ENDBLOCK, string(tkn.Type))
	assert.Equal(t, "@}", string(tkn.Literal))

	tkn = lex.NextToken()
	assert.Equal(t, EOF, string(tkn.Type))

}

func TestCodeBlock2(t *testing.T) {
	sample := `@{ if l.State == "test" { @}`
	lex := NewLexer(string(sample), "TestLexer1")
	tkn := lex.NextToken()
	assert.Equal(t, STARTBLOCK, string(tkn.Type))
	assert.Equal(t, "@{", string(tkn.Literal))

	tkn = lex.NextToken()
	assert.Equal(t, LITERAL, string(tkn.Type))
	assert.Equal(t, "if l.State == \"test\" { ", string(tkn.Literal))

	tkn = lex.NextToken()
	assert.Equal(t, ENDBLOCK, string(tkn.Type))
	assert.Equal(t, "@}", string(tkn.Literal))

	tkn = lex.NextToken()
	assert.Equal(t, EOF, string(tkn.Type))

}

func TestCodeImports(t *testing.T) {
	sample := `@import go imports to be injected
@import go imports to be injected2
@import go imports to be injected3
@import go imports to be injected4
`
	lex := NewLexer(string(sample), "TestLexer1")
	tkn := lex.NextToken()
	assert.Equal(t, IMPORT, string(tkn.Type))
	assert.Equal(t, "go imports to be injected", string(tkn.Literal))
	tkn = lex.NextToken()
	assert.Equal(t, IMPORT, string(tkn.Type))
	assert.Equal(t, "go imports to be injected2", string(tkn.Literal))
	tkn = lex.NextToken()
	assert.Equal(t, IMPORT, string(tkn.Type))
	assert.Equal(t, "go imports to be injected3", string(tkn.Literal))
	tkn = lex.NextToken()
	assert.Equal(t, IMPORT, string(tkn.Type))
	assert.Equal(t, "go imports to be injected4", string(tkn.Literal))

	tkn = lex.NextToken()
	assert.Equal(t, EOF, string(tkn.Type))

}

func TestCodeFunction(t *testing.T) {
	sample := `@func

// A group of functions


func (l *Lexer) readTils(chars []rune) string {
	pos := l.position
	for !l.isAnyChar(chars) {
		l.readChar()
		if l.isEOF() {
			break
		}
	}
	l.readChar()
	return string(l.runes[pos:l.position])
}

@}`
	lex := NewLexer(string(sample), "TestLexer1")
	tkn := lex.NextToken()
	assert.Equal(t, FUNCTS, string(tkn.Type))
	assert.Equal(t, "@func", string(tkn.Literal))

	tkn = lex.NextToken()
	assert.Equal(t, LITERAL, string(tkn.Type))
	//
	//	assert.Equal(t, `    go imports to be injected
	//    go imports to be injected
	//    go imports to be injected
	//`, string(tkn.Literal))

	tkn = lex.NextToken()
	assert.Equal(t, ENDBLOCK, string(tkn.Type))
	assert.Equal(t, "@}", string(tkn.Literal))

	tkn = lex.NextToken()
	assert.Equal(t, EOF, string(tkn.Type))

}

//func TestCodeBaseExtendInvalid(t *testing.T) {
//	sample := `@extend root "Index Page" @{
//@}
//@import "templates"`
//	lex := NewLexer(string(sample), "TestLexer1")
//	var tk = lex.NextToken()
//
//	// Illegal since {@ was not ound
//	fmt.Printf("M:%s at %d\n", tk.Literal, tk.Line)
//	assert.Equal(t, ILLEGAL, string(tk.Type))
//
//}

func TestCodeBaseIf(t *testing.T) {
	sample := `@if errors != nil 
<div>test1</div>
@else
<div>test2</div>
@endif
`
	lex := NewLexer(string(sample), "TestLexer1")
	var tk = lex.NextToken()

	// Illegal since {@ was not ound

	assert.Equal(t, IF, string(tk.Type))
	fmt.Printf("M:%s at %d\n", tk.Literal, tk.Line)

	tk = lex.NextToken()
	assert.Equal(t, LITERAL, string(tk.Type))
	fmt.Printf("M:%s at %d\n", tk.Literal, tk.Line)

	tk = lex.NextToken()
	assert.Equal(t, ELSE, string(tk.Type))
	fmt.Printf("M:%s at %d\n", tk.Literal, tk.Line)

	tk = lex.NextToken()
	assert.Equal(t, LITERAL, string(tk.Type))
	fmt.Printf("M:%s at %d\n", tk.Literal, tk.Line)

	tk = lex.NextToken()
	assert.Equal(t, ENDIF, string(tk.Type))
	fmt.Printf("M:%s at %d\n", tk.Literal, tk.Line)

}
