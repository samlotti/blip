package blip_component
// Do Not Edit
// Generated by Blip
// source blip: html/component/userDetail.blip.html

import (
	"blip/blipUtil"
	"blip/examples/template/blipWebServer/model"
	"context"
	"fmt"
	"io"
	"time"
	html "blip/.blip_generated/blip_html"
)



func UserDetailProcess( user *model.User, c context.Context, w io.Writer ) (terror error) {
    start := time.Now()

	var si = blipUtil.Instance()
	var escaper = si.GetEscaperFor( "html") 
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Catch panic %s: %s\n", "UserDetailProcess", err)
			terror = fmt.Errorf("%v", err)
		}
	    si.RenderComplete(escaper, "userDetail", "html", time.Since(start), terror)
	}()
	// Line: 2
	si.Write(w, []byte("\n"))
	// Line: 6
	si.Write(w, []byte("\n"))
	// Line: 8
	var ctxL1 = context.WithValue(c, "__Blip__", 1)
		// Line: 10
		var contentF1S1 = func() (terror error) {
			// Line: 12
			si.Write(w, []byte("<style>\n</style>\n"))
			// End of content block
			return
		}
		ctxL1 = context.WithValue(ctxL1, "styles", contentF1S1)		// Line: 15
		var contentF1S2 = func() (terror error) {
			// Line: 19
			si.Write(w, []byte("\n<div class=\"container\">\n\n    <div class=\"page-header\">\n        <h1>"))
			// Line: 19
			si.WriteStrSafe(w, user.Name , escaper)
			// Line: 22
			si.Write(w, []byte("</h1>\n    </div>\n\n    <label>"))
			// Line: 22
			si.WriteStrSafe(w, user.Title, escaper)
			// Line: 23
			si.Write(w, []byte("</label>\n    <label>"))
			// Line: 23
			si.WriteStrSafe(w, user.EMail, escaper)
			// Line: 26
			si.Write(w, []byte("</label>\n\n    <div class=\"card\" style=\"width: 18rem;\">\n        <img class=\"card-img-top\" src=\""))
			// Line: 26
			si.WriteStr(w, user.Profile)
			// Line: 29
			si.Write(w, []byte("\" alt=\"Card image cap\">\n        <small>Images generated by <a href=\"https://generated.photos/\" target=\"gf\">Generated Photos</a></small>\n        <div class=\"card-body\">\n            <h5 class=\"card-title\">"))
			// Line: 29
			si.WriteStrSafe(w, user.Title, escaper)
			// Line: 30
			si.Write(w, []byte("</h5>\n            <a href=\"#\" class=\"xbtn xbtn-primary\">"))
			// Line: 30
			si.WriteStrSafe(w, user.EMail, escaper)
			// Line: 32
			si.Write(w, []byte("</a>\n\n            "))
			// Line: 33
			if user.Active {
				// Line: 34
				si.Write(w, []byte("            <div class=\"alert alert-success\" role=\"alert\">Active</div>\n            "))
				// Line: 35
			} else {
				// Line: 36
				si.Write(w, []byte("\n            <div class=\"alert alert-danger\" role=\"alert\">This user is inactive.</div>\n            "))
				// Line: 37
			}
			// Line: 44
			si.Write(w, []byte("\n        </div>\n    </div>\n\n</div>\n\n\n\n"))
			// End of content block
			return
		}
		ctxL1 = context.WithValue(ctxL1, "mainContent", contentF1S2)
	terror = html.RootProcess(user.Name, ctxL1, w)
	if terror != nil { return }
	// Line: 46
	si.Write(w, []byte("\n"))
	return
}