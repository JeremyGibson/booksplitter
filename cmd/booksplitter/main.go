package main

import (
	"flag"
	"fmt"
	"github.com/JeremyGibson/booksplitter/internal"
	"io/fs"
	"os"
	"path/filepath"
)

func processSingleBook(path string) {
	ab := internal.AudioBook{File: path}
	ab.SetChapters()
	ab.SetFormat()
	ab.ExtractChapters()
	fmt.Printf("Finished: %s\n\n\n", path)
}

func processMultipleBooks(rootDir string) {
	err := filepath.Walk(rootDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if filepath.Ext(info.Name()) == ".m4b" {
			processSingleBook(path)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", rootDir, err)
		return
	}
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
