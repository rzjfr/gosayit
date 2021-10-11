package audio

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/mitchellh/go-homedir"
	"github.com/rzjfr/sayit/log"
)

const basePath = "~/.sayit/test/oxford/uk"
const baseURL = "https://www.oxfordlearnersdictionaries.com/media/english/uk_pron"

// to set LogLevel
func InitLogger(verbose bool) {
	log.InitLogger(verbose)
	defer log.Logger.Sync()
}

// TODO: error handling refactor
// play the audio of the word from disk, download it before playing if not present
func Play(word string) error {
	filePath, err := getFile(word) // get teh absolute path of the mp3 file
	if err != nil {
		log.Logger.Debug(err)
		return err
	}

	f, err := os.Open(filePath) // open the mp3 file
	if err != nil {
		log.Logger.Debug(err)
		return err
	}
	defer f.Close()

	stream(f) // stream the file from disk
	return err
}

// Audio streams the given file if it's in mp3 format
func stream(f *os.File) {
	streamer, format, _ := mp3.Decode(f)
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	//speaker.Play(streamer)
	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() { done <- true })))
	<-done
	log.Logger.Debugf("Success in playing from disk, file: %v", f)
}

// Downloads and saves the fils from the given url in the given local path
func saveFile(filePath string, url string) error {
	basePath := filepath.Dir(filePath)
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		if err := os.MkdirAll(basePath, 0775); err != nil {
			return err
		}
		log.Logger.Debugf("Created filePath: %v", filePath)
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("URL: %v, StatusCode: %v", url, resp.StatusCode)
	}
	defer resp.Body.Close()

	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	log.Logger.Debugf("Success in getting URL: %v", url)
	return err
}

// returns path of the file on the disk, if does not exist ties to get it first
func getFile(word string) (string, error) {
	var err error
	fileName := word + "__gb_1.mp3"
	filePath := makePath(fileName)
	if _, err = os.Stat(filePath); os.IsNotExist(err) {
		err = saveFile(filePath, makeURL(fileName))
	}

	return filePath, err
}

// Create a alphabetical sharded path from the given file name
func makePartialPath(fileName string) string {
	return fmt.Sprintf("/%s/%s/%s/%s", fileName[:1], fileName[:3], fileName[:5], fileName)
}

// Creates the local path of the file
func makePath(fileName string) string {
	filePath, _ := homedir.Expand(basePath + makePartialPath(fileName))
	log.Logger.Debugf("Local filePath: %v", filePath)
	return filePath
}

// Creates the URL of the file
func makeURL(fileName string) string {
	return baseURL + makePartialPath(fileName)
}
