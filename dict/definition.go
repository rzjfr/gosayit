package dict

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/corpix/uarand"
	"github.com/mitchellh/go-homedir"
	"github.com/rzjfr/sayit/log"
	"jaytaylor.com/html2text"
)

const basePath = "~/.sayit/test/oxford/uk"
const baseURL = "https://www.oxfordlearnersdictionaries.com/definition/english/"

func Define(word string) error {
	// Get absolute path of the HTML file
	filePath, err := getFile(word)
	if err != nil {
		log.Logger.Debug(err)
		return err
	}
	// Open HTML file
	file, err := os.Open(filePath)
	if err != nil {
		log.Logger.Debug(err)
		return err
	}
	// Load the HTML document to be parsed
	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		return err
	}
	// Print
	_ = heading(doc)
	_ = definitions(doc)
	_ = idioms(doc)
	_ = phrasals(doc)

	return err
}

//TODO: set this method as alternative
func defHTML(doc *goquery.Document) error {
	main := doc.Find("div.responsive_entry_center_wrap")
	main.Not("div#ring-links-box")
	html, _ := main.Html()
	text, _ := html2text.FromString(html, html2text.Options{PrettyTables: true})
	fmt.Println(text)
	return nil
}

func heading(doc *goquery.Document) error {
	doc.Find("div.webtop").Each(func(i int, webtop *goquery.Selection) {
		if webtop.Find("h1.headword").Length() > 0 {
			// header
			head := webtop.Find("h1.headword").Text()
			pos := webtop.Find("span.pos").Text()
			phonetics := webtop.Find(".phons_br").Text()
			fmt.Printf("%s %s %s \n", head, pos, phonetics)
			//variants
			webtop.Find("div.variants").Each(func(i int, variants *goquery.Selection) {
				fmt.Printf("  %s\n", variants.Text())
			})
			// inflections
			inflections := webtop.Find("div.inflections")
			if inflections.Length() > 0 {
				fmt.Printf("  %s\n\n", inflections.Text())
			}
		}
	})
	//origin
	origin := doc.Find("[unbox=\"wordorigin\"] span.body")
	if origin.Length() > 0 {
		fmt.Printf("  Origin: %s\n\n", origin.Text())
	}
	return nil
}

func definitions(doc *goquery.Document) error {
	doc.Find("ol.senses_multiple span.shcut-g").Each(func(i int, sense *goquery.Selection) {
		title := sense.Find("[hclass=\"sense\"] li.span.cf").Text()
		part := sense.Find("h2").Text()
		grammar := sense.Find("[hclass=\"sense\"] span.grammar").Text()
		label := sense.Find("[hclass=\"sense\"] span.dtxt").Text()
		labels := sense.Find("[hclass=\"sense\"] span.labels").Text()
		meaning := sense.Find("[hclass=\"sense\"] span.def").Text()
		output := fmt.Sprintf("%d: %s %s %s %s %s %s \n", i+1, part, grammar, title, label, labels, meaning)
		fmt.Println(regexp.MustCompile(`\s{2}`).ReplaceAllString(output, " "))
		sense.Find("[hclass=\"sense\"] li[htag=\"li\"]").Each(func(i int, example *goquery.Selection) {
			if example.Find("span.x").Text() != "" {
				fmt.Printf("  • %s %s\n", example.Find("span.cf").Text(), example.Find("span.x").Text())
			}
		})
		fmt.Printf("  %s\n", sense.Find("ol.senses_multiple span.xrefs").Text())
	})
	return nil
}

func idioms(doc *goquery.Document) error {
	idioms := doc.Find("span.idm-g")
	if idioms.Length() > 0 {
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
	}
	return nil
}

func phrasals(doc *goquery.Document) error {
	verbs := doc.Find("aside.phrasal_verb_links span.xr-g")
	if verbs.Length() > 0 {
		fmt.Printf("\nPhrasal Verbs: [")
		verbs.Each(func(i int, verb *goquery.Selection) {
			fmt.Printf("%s, ", verb.Text())
		})
		fmt.Printf("]\n")
	}
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
