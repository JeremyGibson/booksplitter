package internal

import (
	"encoding/json"
	"fmt"
	ffmpeg "github.com/JeremyGibson/ffmpeg-go"
	"github.com/schollz/progressbar/v3"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const baseOutputPath = "/data/Audiobooks/extracted"

type ChapterMeta struct {
	Chapters []struct {
		ID        int    `json:"id"`
		TimeBase  string `json:"time_base"`
		Start     int    `json:"start"`
		StartTime string `json:"start_time"`
		End       int    `json:"end"`
		EndTime   string `json:"end_time"`
		Tags      struct {
			Title string `json:"title"`
		} `json:"tags"`
	} `json:"chapters"`
}

type FormatMeta struct {
	Programs []any `json:"programs"`
	Streams  []struct {
		CodecName  string `json:"codec_name"`
		SampleRate string `json:"sample_rate"`
		Channels   int    `json:"channels"`
		BitRate    string `json:"bit_rate"`
		Tags       struct {
			CreationTime time.Time `json:"creation_time"`
			Language     string    `json:"language"`
			HandlerName  string    `json:"handler_name"`
			VendorID     string    `json:"vendor_id"`
		} `json:"tags"`
	} `json:"streams"`
	Format struct {
		Duration string `json:"duration"`
		Tags     struct {
			MajorBrand       string    `json:"major_brand"`
			MinorVersion     string    `json:"minor_version"`
			CompatibleBrands string    `json:"compatible_brands"`
			CreationTime     time.Time `json:"creation_time"`
			Genre            string    `json:"genre"`
			Title            string    `json:"title"`
			Artist           string    `json:"artist"`
			AlbumArtist      string    `json:"album_artist"`
			Album            string    `json:"album"`
			Comment          string    `json:"comment"`
			Copyright        string    `json:"copyright"`
			Date             string    `json:"date"`
		} `json:"tags"`
	} `json:"format"`
}

type AudioBook struct {
	File        string
	outLocation string
	chapters    ChapterMeta
	format      FormatMeta
}

func (m *AudioBook) SetMeta(kwargs ffmpeg.KwArgs) string {
	meta, err := ffmpeg.Probe(m.File, kwargs)
	if err != nil {
		fmt.Printf("%s", err)
	}
	return meta
}

func (m *AudioBook) SetFormat() {
	kwargs := ffmpeg.KwArgs{"v": "error", "select_streams": "a:0", "show_entries": "stream=codec_name,bit_rate,channels,sample_rate : format=duration : format_tags : stream_tags", "of": "json"}
	meta := m.SetMeta(kwargs)
	format := FormatMeta{}
	err := json.Unmarshal([]byte(meta), &format)
	if err != nil {
		fmt.Printf("%s", err)
	}
	m.format = format
}

func (m *AudioBook) SetChapters() {
	kwargs := ffmpeg.KwArgs{"print_format": "json", "show_chapters": "", "sexagesimal": "", "loglevel": "error"}
	meta := m.SetMeta(kwargs)
	chapters := ChapterMeta{}
	err := json.Unmarshal([]byte(meta), &chapters)
	if err != nil {
		fmt.Printf("%s", err)
	}
	m.chapters = chapters
}

func (m *AudioBook) SetOutput() {
	abookDir := filepath.Join(baseOutputPath, normalizeFileName(m.format.Format.Tags.Title))
	err := os.MkdirAll(abookDir, os.ModePerm)
	if err != nil {
		fmt.Printf("%s", err)
		panic(err)
	}
	m.outLocation = abookDir
}

func (m *AudioBook) ExtractChapters() {
	m.SetOutput()
	fmt.Printf("Extracting Chapters\nFrom: %s\nTo: %s\n\n", m.File, m.outLocation)
	pb := progressbar.Default(int64(len(m.chapters.Chapters)))
	for num, chapter := range m.chapters.Chapters {
		pb.Describe(fmt.Sprintf("Extracting: %s ", chapter.Tags.Title))
		fileName := fmt.Sprintf("%03d_%s.m4a", num, normalizeFileName(chapter.Tags.Title))
		if strings.Contains(fileName, "opening_credits") || strings.Contains(fileName, "end_credits") {
			pb.Add(1)
			continue
		}
		outname := filepath.Join(m.outLocation, fileName)
		kwargs := ffmpeg.KwArgs{"ss": chapter.StartTime, "to": chapter.EndTime, "c": "copy", "vn": "", "loglevel": "quiet"}
		err := ffmpeg.Input(m.File).Output(outname, kwargs).Run()
		if err != nil {
			fmt.Printf("%s", err)
			panic(err)
		}
		pb.Add(1)
	}
}
