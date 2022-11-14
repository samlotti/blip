package main

import (
	"blip/internal"
	"flag"
	"fmt"
)

// main
func main() {

	var goptions = internal.BlipOptions{}
	var help bool

	flag.BoolVar(&help, "help", false, "Print help message")
	flag.StringVar(&goptions.Sdir, "dir", "./template", "The source directory containing templates")
	flag.BoolVar(&goptions.Rebuild, "rebuild", false, "rebuild all files")
	flag.BoolVar(&goptions.Watch, "watch", false, "will watch the directory for file names/new files")
	flag.StringVar(&goptions.SupportBranch, "supportBranch", "github.com/samlotti/blip/blipUtil", "Support branch name for include.")

	flag.Parse()

	if help {
		fmt.Println(internal.Name)
		fmt.Printf("Blip Processing: Version: %s\n", internal.Version)
		flag.PrintDefaults()
		return
	}

	internal.GteProcess(&goptions)

}
