package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/JeremyGibson/booksplitter/internal"
	"io"
	"os"
)

func processMultipleAudioFiles(metaFilesPath string) {

}

func processSingleAudioFile(metaFilePath string) {
	jsonFile, err := os.Open(metaFilePath)
	byteValue, _ := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Printf("%s", err)
	}
	af := internal.AudioFileMeta{}
	ae := internal.AudioExtractor{}
	err = json.Unmarshal(byteValue, &af)
	if err != nil {
		fmt.Printf("%s", err)
	}
	ae.AudioMeta = af
	ae.ProcessAudioFile()
}

func main() {
	var audioMetadataVar string
	flag.StringVar(&audioMetadataVar, "p", ".", "Path to a meta file or directory of files.")
	flag.Parse()
	fileinfo, err := os.Stat(audioMetadataVar)
	if err != nil {
		panic(0)
	}
	if fileinfo.IsDir() {
		processMultipleAudioFiles(audioMetadataVar)
	} else {
		processSingleAudioFile(audioMetadataVar)
	}
}
