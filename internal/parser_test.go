package internal

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
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
@}

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

@}
Hello @= user.Name@!
@include simple user.Name, game.P1, game.P2 @
@yield cheese @
@extend table user.Name {@
	@content game {@
		This is a game against @= game.P1 @ and @= game.P2 @
		@{
			// Start block
			//Some code goes in here
			si.WriteStr(w, "dd")
			// End Block
		@}

		@extend inner {@
			@content title {@
				<tr></tr>
			@}
		@}

		@for   error in errors {@
			@= error @
		@endfor

		@if error!=nil {@
			content
		@else
			content2
		@endif
	
	@}
@}
Goodbye!!
`
	lex := NewLexer(sample, "TestLexer1")
	parser := New(lex)
	parser.Parse()

	parser.Dump()

	assert.Equal(t, 2, len(parser.imports))
	assert.Equal(t, 2, len(parser.args))
	assert.Equal(t, 1, len(parser.context))

	for idx, err := range parser.errors {
		fmt.Printf("Err:%d   %v\n", idx, err)
	}
	assert.Equal(t, 0, len(parser.errors))

	// assert.Equal(t, "t", tkn.Literal)
	fmt.Print("\n\n========================================\n")
	parser.renderOutput(os.Stdout, "template", "index", &BlipOptions{
		SupportBranch: "",
	})
	fmt.Print("\n========================================\n\n")

}
