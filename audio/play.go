package audio

import (
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/mitchellh/go-homedir"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const basePath = "~/.sayit/test/oxford/uk"
const baseURL = "https://www.oxfordlearnersdictionaries.com/media/english/uk_pron"

func Play(word string) {
	filePath := getFile(word) // get the file and if not exists download it

	f, _ := os.Open(filePath) // open the file
	defer f.Close()

	stream(f) // stream the saved file

}

func stream(f *os.File) {
	streamer, format, _ := mp3.Decode(f)
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	//speaker.Play(streamer)
	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() { done <- true })))
	<-done
}

func saveFile(filePath string, url string) error {
	basePath := filepath.Dir(filePath)
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		if err := os.MkdirAll(basePath, 0775); err != nil {
			return err
		}
		fmt.Println("path created!")
	}

	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, _ := os.Create(filePath)
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func getFile(word string) string {
	fileName := word + "__gb_1.mp3"
	filePath := makePath(fileName)
	fmt.Println(filePath)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		_ = saveFile(filePath, makeURL(fileName))
		fmt.Println("downloaded!")
	}

	return filePath
}

func makePartialPath(fileName string) string {
	return fmt.Sprintf("/%s/%s/%s/%s", fileName[:1], fileName[:3], fileName[:5], fileName)
}

func makePath(fileName string) string {
	filePath, _ := homedir.Expand(basePath + makePartialPath(fileName))
	return filePath
}

func makeURL(fileName string) string {
	return baseURL + makePartialPath(fileName)
}
