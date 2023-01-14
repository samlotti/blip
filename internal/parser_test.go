package internal

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

func TestParser1(t *testing.T) {
	sample := `@import  "fmt"
@import  "blip/manual"
@arg name string
@arg game Game
@context user User 
@func


func SimpleProcess(n string, p1 string, p2 string, c context.Context, w io.Writer) (terror error) {
	w.Write([]byte(fmt.Sprintf("SimpleProcess:%s %s %s\n", n, p1, p2)))
	return
}
func InnerProcess(c context.Context, w io.Writer) (terror error) {
	var si = blipUtil.Instance()
	w.Write([]byte(fmt.Sprintf("InnerProcess\n")))
	si.CallCtxFunc(c, "title")
	return
}
func TableProcess(n string, c context.Context, w io.Writer) (terror error) {
	var si = blipUtil.Instance()
	w.Write([]byte(fmt.Sprintf("TableProcess:%s \n", n)))
	si.CallCtxFunc(c, "game")
	return
}

type User struct {
	Name string
	Id   string
}

type Game struct {
	Id		string
	P1 		string
	P2 		string
}



func doIt() {
	var u manual.User
	u.Name = "test"
	fmt.Println(u.Name)

	fmt.Println("A test line!")
	fmt.Println("A test line!2")
}
@end

@func
// More embedded functions and var
var name string
func doIt2() {
	fmt.Println("A test line!")
	fmt.Println("A test line!2")
}

func doIt3() {
	fmt.Println("A test line!")
	fmt.Println("A test line!2")
}

@end
Hello @= user.Name@!
@include simple user.Name, game.P1, game.P2 
@yield cheese 
@extend table user.Name 
	@content game 
		This is a game against @= game.P1 @ and @= game.P2 @
		@code
			// Start block
			//Some code goes in here
			si.WriteStr(w, "dd")
			// End Block
		@end

		@extend inner 
			@content title 
				<tr></tr>
			@end
		@end

		@for   error in errors 
			@= error @
		@end

		@if error!=nil 
			content
		@else
			content2
		@end
	
	@end
@end
Goodbye!!
`
	lex := NewLexer(sample, "TestLexer1")
	parser := New(lex)
	parser.Parse()

	// parser.Dump()

	assert.Equal(t, 2, len(parser.imports))
	assert.Equal(t, 2, len(parser.args))
	assert.Equal(t, 1, len(parser.context))

	for idx, err := range parser.errors {
		fmt.Printf("Err:%d   %v\n", idx, err)
	}
	assert.Equal(t, 0, len(parser.errors))

	// assert.Equal(t, "t", tkn.Literal)
	fmt.Print("\n\n========================================\n")
	NewRender(parser).RenderOutput(os.Stdout, "template", "index", "html", "test", &BlipOptions{
		SupportBranch: "",
	})
	fmt.Print("\n========================================\n\n")

}

func TestParserCode(t *testing.T) {
	sample := `@code
    // Sample code. Note the @ is not adjusted
	fmt.printf("This is a @code: %s", test)
@end
@arg test string
`
	lex := NewLexer(sample, "TestLexer1")
	parser := New(lex)
	parser.Parse()

	assert.Equal(t, 0, len(parser.imports))
	assert.Equal(t, 1, len(parser.args))

	for idx, err := range parser.errors {
		fmt.Printf("Err:%d   %v\n", idx, err)
	}
	assert.Equal(t, 0, len(parser.errors))

	// assert.Equal(t, "t", tkn.Literal)
	fmt.Print("\n\n========================================\n")

	fmt.Print("\n\n========================================\n")

	var bresult bytes.Buffer
	NewRender(parser).RenderOutput(&bresult, "template", "index", "html", "test", &BlipOptions{
		SupportBranch: "",
	})
	result := bresult.String()
	fmt.Print(result)
	fmt.Print("\n========================================\n\n")
	assert.True(t, strings.Contains(result, "fmt.printf(\"This is a @code: %s\", test)"))
}

func TestParserText(t *testing.T) {
	sample := `@text
<tr><td>@code</td></tr>
<tr><td>"abc"</td></tr>
@end
@arg test string
@bool= true @
`
	lex := NewLexer(sample, "TestLexer1")
	parser := New(lex)
	parser.Parse()

	assert.Equal(t, 0, len(parser.imports))
	assert.Equal(t, 1, len(parser.args))

	for idx, err := range parser.errors {
		fmt.Printf("Err:%d   %v\n", idx, err)
	}
	assert.Equal(t, 0, len(parser.errors))

	// assert.Equal(t, "t", tkn.Literal)
	fmt.Print("\n\n========================================\n")

	var bresult bytes.Buffer
	NewRender(parser).RenderOutput(&bresult, "template", "index", "html", "test", &BlipOptions{
		SupportBranch: "",
	})
	result := bresult.String()
	fmt.Print(result)
	fmt.Print("\n========================================\n\n")

	// @code is not touched
	assert.True(t, strings.Contains(result, "\\n<tr><td>@code</td></tr>\\n"))
	assert.True(t, strings.Contains(result, "<tr><td>\\\"abc\\\"</td></tr>\\n"))
	assert.True(t, strings.Contains(result, "si.WriteBool(w, true )"))
}
