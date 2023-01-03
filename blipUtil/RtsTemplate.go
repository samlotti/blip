package blipUtil

import (
	"context"
	"fmt"
	"html"
	"io"
	"log"
	"strconv"
	"time"
)

type DefaultBlipMonitor struct {
}

func (m *DefaultBlipMonitor) RenderComplete(escaper IBlipEscaper, name string, langType string, duration time.Duration, err error) {
	// No - op
}

type DebugBlipMonitor struct {
}

func (m *DebugBlipMonitor) RenderComplete(escaper IBlipEscaper, name string, langType string, duration time.Duration, err error) {
	if err != nil {
		fmt.Printf("RenderComplete: %s, %s, %s, Err: %s\n", name, langType, duration, err)
	} else {
		fmt.Printf("RenderComplete: %s, %s, %s\n", name, langType, duration)
	}
}

type TextEscaper struct {
}

var textEscaperInst = TextEscaper{}

func (h *TextEscaper) GetFileType() string {
	return "text"
}
func (h *TextEscaper) Escape(inStr string) string {
	return inStr
}
func TextEscaperInstance() IBlipEscaper {
	return &textEscaperInst
}

type HtmlEscaper struct {
}

var htmlEscaperInst = HtmlEscaper{}

func (h *HtmlEscaper) GetFileType() string {
	return "html"
}
func (h *HtmlEscaper) Escape(inStr string) string {
	return html.EscapeString(inStr)
}
func HtmlEscaperInstance() IBlipEscaper {
	return &htmlEscaperInst
}

// BlipUtil From a code block
type BlipUtil struct {
	verbose  bool
	escapers map[string]IBlipEscaper
	monitor  IBlipMonitor
}

var inst = BlipUtil{
	verbose:  false,
	escapers: map[string]IBlipEscaper{"html": HtmlEscaperInstance(), "text": TextEscaperInstance()},
	monitor:  &DefaultBlipMonitor{},
}

func Instance() *BlipUtil {
	return &inst
}

func (t *BlipUtil) AddEscaper(fileType string, esc IBlipEscaper) {
	t.escapers[fileType] = esc
}

func (t *BlipUtil) SetMonitor(monitor IBlipMonitor) {
	t.monitor = monitor
}

func (t *BlipUtil) GetEscaperFor(fileType string) IBlipEscaper {
	esk, ok := t.escapers[fileType]
	if !ok {
		panic(fmt.Sprintf("Unknown IBlipEscaper for file type: %s", fileType))
	}
	return esk.(IBlipEscaper)
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
func (t *BlipUtil) WriteStrSafe(w io.Writer, bytes string, escaper IBlipEscaper) {
	_, err := w.Write([]byte(escaper.Escape(bytes)))
	if err != nil {
		if err != nil {
			log.Println(fmt.Sprintf("blip had error writing: %s\n", err))
		}
		panic(err)
	}
}
func (t *BlipUtil) RenderComplete(escaper IBlipEscaper, templateName string, langType string, duration time.Duration, err error) {
	t.monitor.RenderComplete(escaper, templateName, langType, duration, err)
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

func (t *BlipUtil) WriteInt(w io.Writer, val int) {
	_, err := w.Write([]byte(strconv.Itoa(val)))
	if err != nil {
		if err != nil {
			log.Println(fmt.Sprintf("blip had error writing: %s\n", err))
		}
		panic(err)
	}
}

func (t *BlipUtil) WriteInt64(w io.Writer, val int64) {
	_, err := w.Write([]byte(strconv.FormatInt(val, 10)))
	if err != nil {
		if err != nil {
			log.Println(fmt.Sprintf("blip had error writing: %s\n", err))
		}
		panic(err)
	}
}
