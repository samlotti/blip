package main

import (
	"flag"
	"fmt"
	"github.com/samlotti/blip/internal"
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
	flag.BoolVar(&goptions.RenderLineNumbers, "renderLineNumbers", false, "Render template line numbers in the generated Go code.  defaults false for easier diffing in source control. ex: adding one line will not show all next line numbers as differences ")

	flag.Parse()

	if help {
		fmt.Println(internal.Name)
		fmt.Printf("Blip Processing: Version: %s\n", internal.Version)
		flag.PrintDefaults()
		return
	}

	internal.GteProcess(&goptions)

}
