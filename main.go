package main

import (
	"log"
	"os"
	"fmt"
	"strings"
	"github.com/fsnotify/fsnotify"
	"github.com/jean-bonilha/goprint"
	"github.com/jean-bonilha/win-toolkit/printer"
)

var folderRoot string = "C:/"

var folderListen string = "C:/PrinterLabel"

// creates a new file watcher
var watcher, errWatcher = fsnotify.NewWatcher()

var printerHandle uintptr

func main() {

	if errWatcher != nil {
		fmt.Println("ERROR CREATE WATCHER", errWatcher)
	}
	defer watcher.Close()

	done := make(chan bool)

	go func() {
		var printerHandle, errOpenPrinter = goprint.GoOpenPrinter(printer.PrinterName)
		if errOpenPrinter != nil {log.Fatalln("Failed to open printer")}
		defer goprint.GoClosePrinter(printerHandle)
		fmt.Println(printer.PrinterName)

		for {
			select {
				// watch for events
			case event := <-watcher.Events:
				if event.Op == 0x1 {
					filePath := event.Name

					fileInfo, err := os.Stat(filePath)
					if err != nil {log.Fatalln("Failed to open fileInfo")}

					if fileInfo.IsDir() {
						handleDir(event.Name)
					} else {
						_ = printer.PrintFile(filePath, printerHandle, true)
					}
				}

				// watch for errors
			case err := <-watcher.Errors:
				fmt.Println("ERROR WATCHER", err)
			}
		}
	}()

	if err := watcher.Add(folderListen); err != nil {
		fmt.Println("ERROR ADD LISTENER", err)
	}

	if err := watcher.Add(folderRoot); err != nil {
		fmt.Println("ERROR ADD LISTENER", err)
	}

	printer.RunPrinterService()

	<-done
}

func handleDir(eventName string) {
	filePath := strings.Replace(eventName, "\\", "/", -1)
	if filePath == folderListen {
		if err := watcher.Add(folderListen); err != nil {
			fmt.Println("ERROR ADD LISTENER", err)
		}
	} else {
		_ = watcher.Remove(filePath)
	}
}
