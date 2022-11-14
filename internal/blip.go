package internal

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var Version = "0.1.6"
var Name = "Blip Template Compiler"

type BlipOptions struct {
	Sdir          string
	Rebuild       bool
	Watch         bool
	SupportBranch string
}

func GteProcess(opt *BlipOptions) {

	fmt.Println(`
██████╗ ██╗     ██╗██████╗     ████████╗███████╗███╗   ███╗██████╗ ██╗      █████╗ ████████╗███████╗███████╗
██╔══██╗██║     ██║██╔══██╗    ╚══██╔══╝██╔════╝████╗ ████║██╔══██╗██║     ██╔══██╗╚══██╔══╝██╔════╝██╔════╝
██████╔╝██║     ██║██████╔╝       ██║   █████╗  ██╔████╔██║██████╔╝██║     ███████║   ██║   █████╗  ███████╗
██╔══██╗██║     ██║██╔═══╝        ██║   ██╔══╝  ██║╚██╔╝██║██╔═══╝ ██║     ██╔══██║   ██║   ██╔══╝  ╚════██║
██████╔╝███████╗██║██║            ██║   ███████╗██║ ╚═╝ ██║██║     ███████╗██║  ██║   ██║   ███████╗███████║
╚═════╝ ╚══════╝╚═╝╚═╝            ╚═╝   ╚══════╝╚═╝     ╚═╝╚═╝     ╚══════╝╚═╝  ╚═╝   ╚═╝   ╚══════╝╚══════╝`)
	fmt.Printf("Blip Processing: Version: %s\n", Version)
	fmt.Printf("Rebuild All: %v\n", opt.Rebuild)
	fmt.Printf("Source folder: %s\n", opt.Sdir)

	processDir(opt.Sdir, opt)

	if opt.Watch {
		watchFiles(opt)
	}
}

func watchFiles(opt *BlipOptions) {
	fmt.Printf("---Watching for file changes in %s\n", opt.Sdir)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				// log.Println("event:", event)
				if event.Has(fsnotify.Create) {
					fi, err := os.Stat(event.Name)
					if err != nil {
						// fmt.Printf("Unable to stat: %s", event.Name)
						return
					}
					if fi.IsDir() {
						go func() {
							fmt.Printf("**** Please restart the process after adding directories!!!\n")
							fmt.Printf("**** Please restart the process after adding directories!!!\n")
							fmt.Printf("**** Please restart the process after adding directories!!!\n")
							fmt.Printf("**** Please restart the process after adding directories!!!\n")
							//fmt.Printf("Watching dir: %s\n", event.Name)
							//watcher.Add(event.Name)
						}()
						return
					}
				}
				if event.Has(fsnotify.Write) {
					// log.Println("modified file:", event.Name)
					slashIdx := strings.LastIndex(event.Name, "/")
					fi, err := os.Stat(event.Name)
					if err != nil {
						return
					}
					sDir := event.Name[0:slashIdx]
					go func() {
						defer func() {
							if err := recover(); err != nil {
								fmt.Printf("Error: %v", err)
							}
						}()
						processFile(sDir, fi, opt)
					}()

				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	// Add a path.
	var dirs []string
	getAllSubDirectories(&dirs, opt.Sdir)
	for _, dir := range dirs {
		err = watcher.Add(dir)
		fmt.Printf("Watching dir: %s\n", dir)
		if err != nil {
			log.Fatal(err)
		}

	}

	// Block main goroutine forever.
	<-make(chan struct{})
}

func getAllSubDirectories(list *[]string, sdir string) []string {
	files, err := ioutil.ReadDir(sdir)
	if err != nil {
		log.Fatal(err)
	}

	*list = append(*list, sdir)
	// fmt.Printf("> %s\n", sdir)

	for _, file := range files {
		if file.IsDir() {
			getAllSubDirectories(list, sdir+"/"+file.Name())
		}
	}
	return *list
}

func processDir(sdir string, opt *BlipOptions) {
	files, err := ioutil.ReadDir(sdir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			processDir(sdir+"/"+file.Name(), opt)
		} else {
			processFile(sdir, file, opt)
		}
	}
}

func processFile(sdir string, file fs.FileInfo, opt *BlipOptions) {
	if !strings.HasSuffix(file.Name(), ".blip") {
		return
	}

	destFName := sdir + "/" + file.Name() + ".go"
	sourceFName := sdir + "/" + file.Name()

	fmt.Printf("Process file: %s\n", sourceFName)
	fmt.Printf("Dest file: %s/%s\n", sdir, destFName)

	inBytes, err := ioutil.ReadFile(sourceFName)
	if err != nil {
		panic(fmt.Sprintf("Error reading file: %s: %s", sourceFName, err))
		return
	}

	sfi, err := os.Stat(sourceFName)
	if err != nil {
		fmt.Printf("Unable to stat: %s : %s", sourceFName, err)
		return
	}
	dfi, err := os.Stat(destFName)
	if err != nil {
		// fmt.Printf("Unable to stat: %s : %s", destFName, err)
		// return
	}

	if err == nil {
		// fmt.Printf("S:%v   D:%v   %v", sfi.ModTime(), dfi.ModTime(), opt.Rebuild)
		if sfi.ModTime().Before(dfi.ModTime()) {
			if !opt.Rebuild {
				fmt.Printf("File %s Not modified since built\n", sourceFName)
				return
			}
		}
	}

	lex := NewLexer(string(inBytes), sourceFName)
	parser := New(lex)
	parser.Parse()

	dirSects := strings.Split(sdir, "/")
	fileSects := strings.Split(file.Name(), ".")

	// _ = os.Remove(destFName)
	dfile, err := os.Create(destFName)
	if err != nil {
		fmt.Printf("Error creating: %s : %s", destFName, err)
		return
	}
	parser.renderOutput(dfile, dirSects[len(dirSects)-1], fileSects[0], opt)
	err = dfile.Close()
	if err != nil {
		fmt.Printf("Error closing: %s : %s", destFName, err)
		return
	}
	fmt.Printf("Wrote to: %s\n", destFName)

}
