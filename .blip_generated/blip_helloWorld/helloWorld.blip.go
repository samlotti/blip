package blip_helloWorld
// Do Not Edit
// Generated by Blip
// source blip: examples/template/helloWorld/helloWorld.blip

import (
	"blip/blipUtil"
	"context"
	"fmt"
	"io"
	"time"
)



func HelloWorldProcess( numArray []int, c context.Context, w io.Writer ) (terror error) {
    start := time.Now()

	var si = blipUtil.Instance()
	var escaper = si.GetEscaperFor( "text") 
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Catch panic %s: %s\n", "HelloWorldProcess", err)
			terror = fmt.Errorf("%v", err)
		}
	    si.RenderComplete(escaper, "helloWorld", "text", time.Since(start), terror)
	}()
	var name string = "World"
	if c.Value("name") != nil {
		name = c.Value("name").(string)
	}
	// Line: 8
	si.Write(w, []byte("\nHello from Blip Templates.\n\n// Note run buildTemplates to generate whenever the template is changed.\n// This template expects an array of numbers passed as an argument\n// An optional \"name\" can be added to the context\n\n"))
	// Line: 9
	if name == "" {
		// Line: 10
		si.Write(w, []byte("    No name!!!\n"))
		// Line: 11
	} else {
		// Line: 11
		si.Write(w, []byte("\n    "))
		// Line: 12
		for idx, c := range numArray { _ = idx
			// Line: 12
			si.Write(w, []byte("      "))
			// Line: 12
			si.WriteInt(w, c)
			// Line: 12
			si.Write(w, []byte(". Hello "))
			// Line: 12
			si.WriteStrSafe(w, name, escaper)
			// Line: 13
			si.Write(w, []byte("!\n    "))
			// Line: 12
		} // end of @for@12
		// Line: 14
		si.Write(w, []byte("\n"))
		// Line: 15
	}
	// Line: 16
	si.Write(w, []byte("\n\n"))
	return
}