package support

import (
	"context"
	"html"
	"io"
)

// THelper From a code block
type THelper struct {
}

var inst = THelper{}

func Instance() *THelper {
	return &inst
}

func (t *THelper) GetCtxStr(c context.Context, key string) string {
	str, ok := (c.Value(key)).(string)
	if !ok {
		return ""
	}
	return str
}

func (t *THelper) CallCtxFunc(c context.Context, key string) {
	f, ok := (c.Value(key)).(func())
	if ok {
		f()
	}
}

func (t *THelper) Write(w io.Writer, bytes []byte) {
	_, err := w.Write(bytes)
	if err != nil {
		panic(err)
	}
}
func (t *THelper) WriteStr(w io.Writer, bytes string) {
	_, err := w.Write([]byte(bytes))
	if err != nil {
		panic(err)
	}

}
func (t *THelper) WriteStrSafe(w io.Writer, bytes string) {
	_, err := w.Write([]byte(html.EscapeString(bytes)))
	if err != nil {
		panic(err)
	}
}
func (t *THelper) IncProcess() {

}
