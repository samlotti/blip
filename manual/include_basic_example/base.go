package include_basic_example

import (
	"blip/manual"
	"blip/support"
	"context"
	"fmt"
	"html"
	"io"
)

// Example Base template
var base1 = []byte("<h1>")
var base2 = []byte("</h1>\n")

func BaseProcess(c context.Context, w io.Writer) (terror error) {
	var si = support.Instance()
	defer func() {
		if err := recover(); err != nil {
			terror = fmt.Errorf("%v", err)
		}
	}()

	// Context variables
	// @context title string
	title := c.Value("title").(string)

	// @context user User
	user := c.Value("user").(manual.User)

	si.Write(w, base1)

	// @= title@
	si.WriteStr(w, html.EscapeString(title))

	si.Write(w, base2)

	// @= user.Name@
	si.WriteStr(w, user.Name)

	return
}
