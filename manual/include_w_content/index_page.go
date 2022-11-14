package include_w_content

import (
	"blip/blipUtil"
	"blip/manual"
	"context"
	"fmt"
	"io"
)

var indexpage1 []byte = []byte("<td>")
var indexpage2 []byte = []byte("</td>\n")

// IndexPageContentProcess
// @arg games []Game
func IndexPageContentProcess(games []manual.Game, c context.Context, w io.Writer) (terror error) {
	var si = blipUtil.Instance()
	defer func() {
		if err := recover(); err != nil {
			terror = fmt.Errorf("%v", err)
		}
	}()

	// @content myContent {@
	var c1 context.Context
	var f1 = func() {
		// @{ for _, game := range games { @}
		for _, game := range games {

			// <tr>
			si.Write(w, indexpage1)

			// @= game.Opponent @
			si.WriteStr(w, game.Opponent)

			var i2 = 500
			si.WriteStr(w, fmt.Sprintf("%d", i2))

			// </tr>\n
			si.Write(w, indexpage2)
		}
	}
	c1 = context.WithValue(c, "myContent", f1)

	// @include Base {@
	terror = BaseProcess(c1, w)
	if terror != nil {
		return
	}

	// @{ } @}

	return
}
