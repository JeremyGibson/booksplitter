package internal

import (
	"github.com/frolovo22/tag"
	"log"
)

func Tag(filename string) tag.Metadata {
	tags, err := tag.ReadFile(filename)
	if err != nil {
		log.Fatal("Error while opening mp3 file: ", err)
	}
	return tags
}
