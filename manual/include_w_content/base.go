package include_w_content

import (
	"blip/manual"
	"blip/support"
	"context"
	"fmt"
	"html"
	"io"
)

// Example Base template
var base1 = []byte("<body>")
var base2 = []byte("</body>\n")
var base3 = []byte("\n")

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
	// if err = si.Write(w, []byte(support.Instance().GetCtxStr(c, "title"))); err != nil {
	si.WriteStr(w, html.EscapeString(title))

	si.Write(w, base3)

	// @yield myContent @
	si.CallCtxFunc(c, "myContent")

	// @yield myJavascript @
	si.CallCtxFunc(c, "myJavascript")

	si.Write(w, base2)

	// @= user.Name@
	si.WriteStr(w, user.Name)

	return

}
