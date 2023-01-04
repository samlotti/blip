package main

import (
	"context"
	"fmt"
	"github.com/samlotti/blip/blipUtil"
	"github.com/samlotti/blip/blipped/component"
	"github.com/samlotti/blip/blipped/html"
	"github.com/samlotti/blip/examples/template/blipWebServer/model"
	"log"
	"net/http"
	"strconv"
	"strings"
)

//go:generate blip -dir=../ --supportBranch=blip/blipUtil
func main() {
	fmt.Printf("Running the blip server example:  http://localhost:8181\n")

	// Show timings of the individual template renders
	blipUtil.Instance().SetMonitor(&blipUtil.DebugBlipMonitor{})

	http.HandleFunc("/", Index)
	http.HandleFunc("/favicon.ico", http.NotFound)
	http.HandleFunc("/users/listAll", UListAll)
	http.HandleFunc("/users/listActive", UListActive)
	http.HandleFunc("/users/view/", UView)

	if err := http.ListenAndServe(":8181", nil); err != nil {
		log.Fatal(err)
	}

}

func BaseCtx() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "title", "Blip Example")
	return ctx
}

func Index(w http.ResponseWriter, r *http.Request) {
	ctx := BaseCtx()
	html.IndexRender(ctx, w)
}

func UListAll(w http.ResponseWriter, r *http.Request) {
	ctx := BaseCtx()
	users := model.GetUsers()
	html.ListUsersRender(users, ctx, w)
}

func UListActive(w http.ResponseWriter, r *http.Request) {
	ctx := BaseCtx()
	users := model.GetUsers()
	html.ListActiveUsersRender(users, ctx, w)
}

func UView(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/users/view/"))
	ctx := BaseCtx()
	user := model.GetUser(id)
	component.UserDetailRender(user, ctx, w)
}
