package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	ffmpeg "github.com/JeremyGibson/ffmpeg-go"
	"github.com/disintegration/imaging"
	"image/jpeg"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const baseOutputPathAudio = "/data/audio_processing/processed"
const tempImageFile = "/tmp/image.jpg"

var coverImageBuffer = new(bytes.Buffer)

type AudioFileMeta struct {
	FileInfo struct {
		FilePath  string `json:"file_path"`
		Album     string `json:"album"`
		ImageTime string `json:"image_time"`
		Artist    string `json:"artist"`
		Year      string `json:"year"`
		Date      string `json:"date"`
		Genre     string `json:"genre"`
	} `json:"file_info"`
	Tracks []struct {
		Title    string `json:"title"`
		FromSecs string `json:"from_secs"`
		ToSecs   string `json:"to_secs"`
	} `json:"tracks"`
}

type AudioExtractor struct {
	audioMeta AudioFileMeta
	videoMeta FFMpegVideoMeta
}

func ReadFrameAsJpeg(inFileName string, frameNum int) io.Reader {
	fmt.Printf("Extracting a cover image for: %s\n\n", inFileName)
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input(inFileName).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg", "loglevel": "quiet"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		panic(err)
	}
	return buf
}

func (a *AudioFileMeta) getFPS() int64 {
	kwargs := ffmpeg.KwArgs{"v": "error", "select_streams": "v:0", "show_entries": "stream=codec_name,bit_rate,channels,sample_rate : format=duration : format_tags : stream_tags", "of": "json"}
	meta, err := ffmpeg.Probe(a.FileInfo.FilePath, kwargs)
	if err != nil {
		fmt.Printf("%s", err)
	}
	videoMeta := FFMpegVideoMeta{}
	err = json.Unmarshal([]byte(meta), &videoMeta)
	if err != nil {
		fmt.Printf("%s", err)
	}
	fps := strings.Split(videoMeta.Streams[0].RFrameRate, "/")[0]
	fpsi, err := strconv.ParseInt(fps, 10, 64)
	return fpsi
}

func (a *AudioFileMeta) getTimeInSeconds(time string) int64 {
	// A timestamp string formatted like "00:00:00"
	timeArray := strings.Split(time, ":")
	hours, err := strconv.ParseInt(timeArray[0], 10, 64)
	minutes, err := strconv.ParseInt(timeArray[1], 10, 64)
	seconds, err := strconv.ParseInt(timeArray[2], 10, 64)
	seconds += minutes * 60
	seconds += (hours * 60) * 60
	if err != nil {
		fmt.Printf("%s", err)
	}
	return seconds
}

func (a *AudioFileMeta) extractImageFromSource() {
	fps := a.getFPS()
	timeinsecs := a.getTimeInSeconds(a.FileInfo.ImageTime)
	frametograb := fps * timeinsecs
	reader := ReadFrameAsJpeg(a.FileInfo.FilePath, int(frametograb))
	img, err := imaging.Decode(reader)
	if err != nil {
		panic(err)
	}
	dstImage800 := imaging.Resize(img, 800, 0, imaging.Lanczos)
	err = jpeg.Encode(coverImageBuffer, dstImage800, nil)
	if err != nil {
		return
	}
}

func (a *AudioFileMeta) setOutput() string {
	audioFileDir := filepath.Join(
		baseOutputPathAudio,
		normalizeFileName(a.FileInfo.Artist),
		normalizeFileName(a.FileInfo.Album),
	)
	err := os.MkdirAll(audioFileDir, os.ModeDir)
	if err != nil {
		fmt.Printf("%s", err)
		panic(err)
	}
	return audioFileDir
}

func (a *AudioFileMeta) setTrackMetadata(track string) {

}

func (ae *AudioExtractor) ProcessAudioFile() {
	ae.audioMeta.extractImageFromSource()
	extractTo := ae.audioMeta.setOutput()
	for num, track := range ae.audioMeta.Tracks {
		makeQuiet := true
		fileName := fmt.Sprintf("%03d_%s.flac", num+1, normalizeFileName(track.Title))
		fmt.Printf("Extracting: %s\n", fileName)
		outName := filepath.Join(extractTo, fileName)
		ikwargs := ffmpeg.KwArgs{"ss": track.FromSecs}
		okwargs := ffmpeg.KwArgs{"to": track.ToSecs, "sample_fmt": "s16", "vn": ""}
		if makeQuiet == true {
			okwargs = ffmpeg.KwArgs{"to": track.ToSecs, "vn": "", "sample_fmt": "s16", "loglevel": "quiet"}
		}
		err := ffmpeg.Input(a.FileInfo.FilePath, ikwargs).
			Output(outName, okwargs).
			OverWriteOutput().ErrorToStdOut().Run()
		if err != nil {
			fmt.Printf("%s", err)
			panic(err)
		}
	}

}
