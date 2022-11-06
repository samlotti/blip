package main

import (
	"blip/manual"
	"blip/manual/include_basic_example"
	"context"
	"os"
)

func main() {
	ctx := context.TODO()

	// The caller must add the context values that would be required.
	ctx = context.WithValue(ctx, "title", "GoTemplateEngine")
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
		manual.Game{
			Id:       "3",
			Opponent: "Joe",
		},
	}
	include_basic_example.IndexPageProcess(games, ctx, os.Stdout)

}
