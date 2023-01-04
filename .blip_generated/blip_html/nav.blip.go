package blip_html
// Do Not Edit
// Generated by Blip
// source blip: ./html/nav.blip.html

import (
	"blip/blipUtil"
	"context"
	"fmt"
	"io"
	"time"
)



func NavProcess( pageName string, c context.Context, w io.Writer ) (terror error) {
    start := time.Now()

	var si = blipUtil.Instance()
	var escaper = si.GetEscaperFor( "html") 
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Catch panic %s: %s\n", "NavProcess", err)
			terror = fmt.Errorf("%v", err)
		}
	    si.RenderComplete(escaper, "nav", "html", time.Since(start), terror)
	}()
	// Line: 13
	si.Write(w, []byte("\n<nav class=\"navbar navbar-expand-lg navbar-inverse\">\n    <div class=\"container-fluid nav-collapsable\" style=\"margin-top:16px;\">\n        <div class=\"navbar-header\">\n            <button type=\"button\" class=\"navbar-toggle collapsed\" data-toggle=\"collapse\" data-target=\"#navbar\" aria-expanded=\"false\" aria-controls=\"navbar\">\n                <span class=\"sr-only\">Toggle navigation</span>\n                <span class=\"icon-bar\"></span>\n                <span class=\"icon-bar\"></span>\n                <span class=\"icon-bar\"></span>\n            </button>\n            <a class=\"navbar-brand\" href=\"#\">"))
	// Line: 13
	si.WriteStrSafe(w, pageName , escaper)
	// Line: 25
	si.Write(w, []byte("</a>\n        </div>\n        <div class=\"collapse navbar-collapse\" id=\"navbar\">\n            <ul class=\"nav navbar-nav\">\n                <li><a href=\"/\">Home</a><li>\n                <li><a href=\"/users/listAll\">All Users</a><li>\n                <li><a href=\"/users/listActive\">Active Users</a></li>\n            </ul>\n        </div>\n\n    </div>\n</nav>\n"))
	return
}