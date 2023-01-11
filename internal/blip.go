package internal

import (
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

var Version = "0.8.8"
var Name = "Blip Template Compiler"

type BlipOptions struct {
	Sdir              string
	Rebuild           bool
	Watch             bool
	SupportBranch     string
	RenderLineNumbers bool
}

func GteProcess(opt *BlipOptions) {

	fmt.Printf(" __          __     ___  ___        __            ___  ___  __  \n|__) |    | |__)     |  |__   |\\/| |__) |     /\\   |  |__  /__` \n|__) |___ | |        |  |___  |  | |    |___ /~~\\  |  |___ .__/ \n")

	fmt.Printf("Blip Processing: Version: %s\n", Version)
	fmt.Printf("Rebuild All: %v\n", opt.Rebuild)
	fmt.Printf("Render: LineNumbers %v\n", opt.Rebuild)
	fmt.Printf("Source folder: %s\n", opt.Sdir)

	processDir(opt.Sdir, opt)
	fmt.Printf("\n")

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
	fileType := "text"

	if strings.HasSuffix(file.Name(), ".go") {
		return
	}

	if !strings.Contains(file.Name(), ".blip") {
		return
	}

	trimmedName := file.Name()
	if !strings.HasSuffix(file.Name(), ".blip") {
		sections := strings.Split(file.Name(), ".")
		fileType = sections[len(sections)-1]
		trimmedName = strings.TrimSuffix(trimmedName, "."+fileType)
	}

	// destDir
	destDir := findGoMod() + "/blipped/" + path.Base(sdir)
	err := os.MkdirAll(destDir, 0755)
	if err != nil {
		panic(fmt.Sprintf("Error creating directory: %s: %s", destDir, err))
		return
	}

	destFName := destDir + "/" + trimmedName + ".go"
	sourceFName := sdir + "/" + file.Name()

	fmt.Printf("\nProcess blip: %s --> %s", sourceFName, destDir)

	inBytes, err := ioutil.ReadFile(sourceFName)
	if err != nil {
		panic(fmt.Sprintf("Error reading file: %s: %s", sourceFName, err))
		return
	}

	sfi, err := os.Stat(sourceFName)
	if err != nil {
		fmt.Printf(" -- Unable to stat: %s : %s\n", sourceFName, err)
		return
	}
	dfi, err := os.Stat(destFName)
	if err == nil {
		if sfi.ModTime().Before(dfi.ModTime()) {
			if !opt.Rebuild {
				fmt.Printf("-- Not modified \n")
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
		fmt.Printf(" -- Error: %s\n", err)
		return
	}
	NewRender(parser).RenderOutput(dfile, dirSects[len(dirSects)-1], fileSects[0], fileType, sourceFName, opt)
	err = dfile.Close()
	if err != nil {
		fmt.Printf(" -- Error: %s\n", err)
		return
	}

	if parser.hasErrors() {
		fmt.Printf("\n\n*** Errors found in : %s\n", file.Name())
		for _, err := range parser.errors {
			fmt.Printf("Error: %d:%d %s\n", err.lineNum, err.linePos, err.msg)
		}
		fmt.Printf("\n\n")
	}
	// fmt.Printf("\n")

}

// findGoMod --
// Goes up the directories until it find the go.mod file
func findGoMod() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("Cannot get directory: %s", err))
	}
	for {
		_, err := os.Stat(dir + "/go.mod")
		if err == nil {
			// This is the bae directory!
			return dir
		}
		if errors.Is(err, os.ErrNotExist) {
			dir, _ = path.Split(strings.TrimSuffix(dir, "/"))
			continue
		} else {
			panic(fmt.Sprintf("Cannot get base directory, searching for directory with go.mod: %s", err))
		}
	}
}
