package blip_template
// Do Not Edit
// Generated by Blip
// source blip: examples/template/test.blip

import (
	"blip/blipUtil"
	"context"
	"fmt"
	"io"
	"time"
)



func TestProcess( c context.Context, w io.Writer ) (terror error) {
    start := time.Now()

	var si = blipUtil.Instance()
	var escaper = si.GetEscaperFor( "text") 
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Catch panic %s: %s\n", "TestProcess", err)
			terror = fmt.Errorf("%v", err)
		}
	    si.RenderComplete(escaper, "test", "text", time.Since(start), terror)
	}()
	title := c.Value("title").(string)
	// Line: 3
	si.Write(w, []byte("<head>\n    <title>"))
	// Line: 3
	si.WriteStrSafe(w, title, escaper)
	// Line: 10
	si.Write(w, []byte("</title>\n\n</head>\n\n\n\n<body>\n    "))
	// Line: 11
	terror = si.CallCtxFunc(c, "styles")
	if terror != nil { return }
	// Line: 11
	si.Write(w, []byte("    "))
	// Line: 12
	terror = si.CallCtxFunc(c, "body")
	if terror != nil { return }
	// Line: 13
	si.Write(w, []byte("\n    "))
	// Line: 14
	terror = si.CallCtxFunc(c, "javascript")
	if terror != nil { return }
	// Line: 18
	si.Write(w, []byte("\n    (c)Copyright 2022\n</body>\n\n"))
	return
}