package main

import (
	"blip/blipUtil"
	"blip/examples/template/helloWorld"
	"context"
	"os"
)

//go:generate blip -dir=../.. --supportBranch=blip/blipUtil --rebuild
func main() {
	// Note run buildTemplates to generate whenever the template is changed.

	// This template expects an array of numbers passed as an argument
	// An optional "name" can be added to the context

	// Show runetime of the template
	blipUtil.Instance().SetMonitor(&blipUtil.DebugBlipMonitor{})

	var ctx = context.Background()

	err := helloWorld.HelloWorldProcess([]int{1, 2, 3, 4, 5}, ctx, os.Stdout)
	if err != nil {
		panic(err)
	}

	ctx = context.WithValue(ctx, "name", "Blip Programmer")

	err = helloWorld.HelloWorldProcess([]int{1, 2, 3, 4, 5}, ctx, os.Stdout)
	if err != nil {
		panic(err)
	}

	ctx = context.WithValue(ctx, "name", "")

	err = helloWorld.HelloWorldProcess([]int{1, 2, 3, 4, 5}, ctx, os.Stdout)
	if err != nil {
		panic(err)
	}
}
