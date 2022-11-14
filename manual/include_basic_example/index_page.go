package include_basic_example

import (
	"blip/blipUtil"
	"blip/manual"
	"context"
	"fmt"
	"io"
)

var indexpage1 []byte = []byte("<tr>")
var indexpage2 []byte = []byte("</tr>\n")

// IndexPageProcess
// @arg games []Game
func IndexPageProcess(games []manual.Game, c context.Context, w io.Writer) (terror error) {
	var si = blipUtil.Instance()
	defer func() {
		if err := recover(); err != nil {
			terror = fmt.Errorf("%v", err)
		}
	}()

	// @include Base @
	terror = BaseProcess(c, w)
	if terror != nil {
		return
	}

	// @{ for _, game := range games { @}
	for _, game := range games {

		// <tr>
		si.Write(w, indexpage1)

		// @= game.Opponent @
		si.WriteStr(w, game.Opponent)

		// </tr>\n
		si.Write(w, indexpage2)
	}
	// @{ } @}

	return
}
