package main

import (
	"blip/manual"
	"blip/manual/include_w_content"
	"context"
	"fmt"
	"os"
)

func main() {
	ctx := context.TODO()

	// The caller must add the context values that would be required.
	ctx = context.WithValue(ctx, "title", "GoTemplateEngine <b>With</b> Content")
	ctx = context.WithValue(ctx, "user", manual.User{
		Name: "Gte",
		Id:   "45",
	})
	// manual.Base_Process(ctx, os.Stdout)

	games := []manual.Game{
		manual.Game{
			Id:       "1",
			Opponent: "Mike",
		},
		manual.Game{
			Id:       "2",
			Opponent: "Steve",
		},
	}
	err := include_w_content.IndexPageContentProcess(games, ctx, os.Stdout)
	if err != nil {
		fmt.Println(err)
	}

}
