package dict

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/corpix/uarand"
	"github.com/mitchellh/go-homedir"
	"github.com/rzjfr/sayit/log"
)

const basePath = "~/.sayit/test/oxford/uk"
const baseURL = "https://www.oxfordlearnersdictionaries.com/definition/english/"

//TODO: get quiet flag of the cli
func Define(word string) error {
	filePath, err := getFile(word) // get the absolute path of the html file
	if err != nil {
		log.Logger.Debug(err)
		return err
	}

	file, err := os.Open(filePath)
	if err != nil {
		log.Logger.Debug(err)
		return err
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		return err
	}

	_ = heading(doc)
	_ = phonetics(doc)
	_ = definitions(doc)
	_ = idioms(doc)
	return nil
}

func phonetics(doc *goquery.Document) error {
	fmt.Printf("%s\n\n", doc.Find(".phons_br").Text())

	return nil
}

func heading(doc *goquery.Document) error {
	doc.Find("div.webtop").Each(func(i int, webtop *goquery.Selection) {
		head := strings.Trim(webtop.Find("h1.headword").Text(), " ")
		pos := strings.Trim(webtop.Find("span.pos").Text(), " ")
		labels := strings.Trim(webtop.Find("span.labels").Text(), " ")
		variants := strings.Trim(webtop.Find("div.variants").Text(), " ")
		fmt.Printf("%s %s %s %s ", head, pos, labels, variants)
	})
	return nil
}

func definitions(doc *goquery.Document) error {
	doc.Find("[hclass=\"sense\"]").Each(func(i int, sense *goquery.Selection) {
		title := sense.Find("li.span.cf").Text()
		grammar := sense.Find("span.grammar").Text()
		label := sense.Find("span.dtxt").Text()
		labels := sense.Find("span.labels").Text()
		meaning := sense.Find("span.def").Text()
		fmt.Printf("%d: %s %s %s %s %s \n", i+1, grammar, title, label, labels, meaning)
		sense.Find("li[htag=\"li\"]").Each(func(i int, example *goquery.Selection) {
			if example.Find("span.x").Text() != "" {
				fmt.Printf("  • %s %s\n", example.Find("span.cf").Text(), example.Find("span.x").Text())
			}
		})
	})
	return nil
}

func idioms(doc *goquery.Document) error {
	fmt.Printf("\nIdioms:\n")
	doc.Find("span.idm-g").Each(func(i int, idiom *goquery.Selection) {
		meaning := idiom.Find("span.idm").Text()
		fmt.Printf("  %s:\n", meaning)
		idiom.Find("ul.examples li[htag=\"li\"]").Each(func(i int, example *goquery.Selection) {
			if example.Find("li span.x").Text() != "" {
				fmt.Printf("  - %s %s\n", example.Find("li span.cf").Text(), example.Find("li span.x").Text())
			}
		})
	})
	return nil
}

// Downloads and saves the fils from the given url in the given local path
func saveFile(filePath string, url string) error {
	basePath := filepath.Dir(filePath)
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		if err := os.MkdirAll(basePath, 0775); err != nil {
			return err
		}
		log.Logger.Debugf("Success in creating filePath: %v", filePath)
	}

	// Define the HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", uarand.GetRandom())

	// Run the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("URL: %v, StatusCode: %v", url, resp.StatusCode)
	}
	log.Logger.Debugf("statusCode %d", resp.StatusCode)

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
	fileName := word + ".html"
	filePath := makePath(fileName)
	if _, err = os.Stat(filePath); os.IsNotExist(err) {
		err = saveFile(filePath, baseURL+word)
	}

	return filePath, err
}

// Creates the local path of the file
func makePath(fileName string) string {
	partialPath := fmt.Sprintf("/%s/%s/%s/%s", fileName[:1], fileName[:3], fileName[:5], fileName)
	filePath, _ := homedir.Expand(basePath + partialPath)
	log.Logger.Debugf("Local filePath: %v", filePath)
	return filePath
}
