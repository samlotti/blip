package blipUtil

import (
	"context"
	"fmt"
	"html"
	"io"
	"log"
)

// BlipUtil From a code block
type BlipUtil struct {
	verbose bool
}

var inst = BlipUtil{
	verbose: false,
}

func Instance() *BlipUtil {
	return &inst
}

func (t *BlipUtil) GetCtxStr(c context.Context, key string) string {
	str, ok := (c.Value(key)).(string)
	if !ok {
		return ""
	}
	return str
}

func (t *BlipUtil) CallCtxFunc(c context.Context, key string) (terror error) {
	f, ok := (c.Value(key)).(func() error)
	if ok {
		terror = f()
		if t.verbose {
			if terror != nil {
				log.Println(fmt.Sprintf("blip had error from content include: %s\n", terror))
			}
		}
		return
	}
	return
}

func (t *BlipUtil) Write(w io.Writer, bytes []byte) {
	_, err := w.Write(bytes)
	if err != nil {
		if err != nil {
			log.Println(fmt.Sprintf("blip had error writing: %s\n", err))
		}
		panic(err)
	}
}
func (t *BlipUtil) WriteStr(w io.Writer, bytes string) {
	_, err := w.Write([]byte(bytes))
	if err != nil {
		if err != nil {
			log.Println(fmt.Sprintf("blip had error writing: %s\n", err))
		}
		panic(err)
	}

}
func (t *BlipUtil) WriteStrSafe(w io.Writer, bytes string) {
	_, err := w.Write([]byte(html.EscapeString(bytes)))
	if err != nil {
		if err != nil {
			log.Println(fmt.Sprintf("blip had error writing: %s\n", err))
		}
		panic(err)
	}
}
func (t *BlipUtil) IncProcess() {

}

// AddCtxError
// Will add an 'error' context variable that is of  []string.
// If it already exists then the new message is added to the end.
// To support
//  @context errors []string = make([]string,0)
func (t *BlipUtil) AddCtxError(ctx context.Context, message string) context.Context {
	var errors []string = make([]string, 0)
	if ctx.Value("errors") != nil {
		errors = ctx.Value("errors").([]string)
	}
	errors = append(errors, message)
	ctx = context.WithValue(ctx, "errors", errors)
	return ctx
}

// HasError
// Return true if there were errors added to the context
// To support
//  @context errors []string = make([]string,0)
func (t *BlipUtil) HasError(ctx context.Context) bool {
	return ctx.Value("errors") != nil
}
