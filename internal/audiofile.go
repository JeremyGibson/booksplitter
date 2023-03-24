package internal

import (
	"bytes"
	"fmt"
	"github.com/frolovo22/tag"
	"github.com/schollz/progressbar/v3"
	"image"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

const baseOutputPathAudio = "/data/audiofiles/extracted"
const baseImageDir = "/data/audiofiles/images"

var coverImageBuffer = new(image.YCbCr)

type AudioFileMeta struct {
	FileInfo struct {
		FilePath  string `json:"file_path"`
		Source    string `json:"source"`
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
	AudioMeta AudioFileMeta
	videoMeta FFMpegVideoMeta
}

func executeCmd(cmd string) {
	command := exec.Command("bash", "-c", cmd)
	var out bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &out
	command.Stderr = &stderr
	err := command.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
	}

}
func ReadFrameAsJpeg(inFileName string, am AudioFileMeta) {
	fmt.Printf("Extracting a cover image for: %s\n\n", inFileName)
	tempImageFile := filepath.Join(baseImageDir, normalizeFileName(am.FileInfo.Album)+".jpg")
	firstCmd := "ffmpeg -ss " + am.FileInfo.ImageTime + " -i " + inFileName + " -frames:v 1 -vf scale=800:-1 " + tempImageFile + " -y"
	executeCmd(firstCmd)
}

func (a *AudioFileMeta) setOutput() string {
	audioFileDir := filepath.Join(
		baseOutputPathAudio,
		normalizeFileName(a.FileInfo.Artist),
		normalizeFileName(a.FileInfo.Album),
	)
	if _, err := os.Stat(audioFileDir); os.IsNotExist(err) {
		err := os.MkdirAll(audioFileDir, os.ModePerm)
		fmt.Printf("%s", err)
	}
	return audioFileDir
}

func (a *AudioFileMeta) setTrackMetadata(trackPath string, title string, trackNum int) {
	tags, err := tag.ReadFile(trackPath)
	if err != nil {
		fmt.Println(err)
	}
	year, err := strconv.Atoi(a.FileInfo.Year)
	time, err := time.Parse("2006-01-02", a.FileInfo.Date)
	tags.SetTitle(title)
	tags.SetTrackNumber(trackNum, len(a.Tracks))
	tags.SetArtist(a.FileInfo.Artist)
	tags.SetYear(year)
	tags.SetDate(time)
	tags.SetDescription(a.FileInfo.Source)
	tags.SetAlbumArtist(a.FileInfo.Artist)
	tags.SetAlbum(a.FileInfo.Album)
	tags.SetGenre(a.FileInfo.Genre)
	tags.SaveFile(trackPath)
}

func (ae *AudioExtractor) ProcessAudioFile() {
	ReadFrameAsJpeg(ae.AudioMeta.FileInfo.FilePath, ae.AudioMeta)
	extractTo := ae.AudioMeta.setOutput()
	pb := progressbar.Default(int64(len(ae.AudioMeta.Tracks)))
	for num, track := range ae.AudioMeta.Tracks {
		fileName := fmt.Sprintf("%03d_%s.flac", num+1, normalizeFileName(track.Title))
		pb.Describe(fmt.Sprintf("Extracting: %s", fileName))
		fmt.Printf("Extracting: %s\n", fileName)
		outName := filepath.Join(extractTo, fileName)
		cmd := "ffmpeg -i " + ae.AudioMeta.FileInfo.FilePath + " -ss " + track.FromSecs + " -to " + track.ToSecs + " -sample_fmt s16 -q:a 0 -map a " + outName + " -y"
		executeCmd(cmd)
		ae.AudioMeta.setTrackMetadata(outName, track.Title, num+1)
		pb.Add(1)
	}
}
