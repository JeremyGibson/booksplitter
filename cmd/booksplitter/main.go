package main

import (
	"flag"
	"fmt"
	"github.com/JeremyGibson/booksplitter/internal"
	"os"
)

func processSingleBook(path string) {
	ab := internal.AudioBook{File: path}
	ab.SetChapters()
	ab.SetFormat()
	ab.ExtractChapters()
	fmt.Printf("Finished: %s/n", path)
}

func processMultipleBooks(rootDir string) {
	fmt.Printf("Processing: %s/n", rootDir)
}

func main() {
	var bookpathVar string
	flag.StringVar(&bookpathVar, "p", ".", "Path to an m4b file.")
	flag.Parse()
	fileinfo, err := os.Stat(bookpathVar)
	if err != nil {
		panic(0)
	}
	if fileinfo.IsDir() {
		processMultipleBooks(bookpathVar)
	} else {
		processSingleBook(bookpathVar)
	}
}
