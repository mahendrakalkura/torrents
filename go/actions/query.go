package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/lestrrat/go-libxml2"
	"github.com/lestrrat/go-libxml2/types"
	"github.com/mahendrakalkura/torrents/go/settings"
)

// Item ...
type Item struct {
	Category string
	Title    string
	URL      string
	Seeds    int
	URLs     []string
	Magnet   string
}

// Items ...
type Items []Item

func (items Items) Len() int {
	return len(items)
}

func (items Items) Swap(one, two int) {
	items[one], items[two] = items[two], items[one]
}

func (items Items) Less(one, two int) bool {
	if items[one].Category < items[two].Category {
		return true
	}
	if items[one].Category > items[two].Category {
		return false
	}
	return items[one].Seeds < items[two].Seeds
}

func exists(items Items, message Item) bool {
	for _, item := range items {
		if item.URL == message.URL {
			return true
		}
	}
	return false
}

func consumer(waitGroup *sync.WaitGroup, count int, outgoing chan string, incoming chan Items) {
	defer waitGroup.Done()

	items := Items{}
	for messages := range incoming {
		for _, message := range messages {
			if exists(items, message) {
				continue
			}
			items = append(items, message)
		}
		count--
		if count == 0 {
			close(outgoing)
			close(incoming)
			break
		}
	}
	sort.Sort(Items(items))
	marshal, marshalErr := json.Marshal(items)
	if marshalErr != nil {
		log.Fatalln(marshalErr)
	}
	err := ioutil.WriteFile("torrents.json", marshal, 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

func producer(waitGroup *sync.WaitGroup, outgoing chan string, incoming chan Items) {
	defer waitGroup.Done()

	for url := range outgoing {
		items, itemsErr := getItems(url)
		if itemsErr != nil {
			outgoing <- url
			continue
		}
		incoming <- items
	}
}

func getCategory(node types.Node) string {
	xPath, xPathErr := node.Find(`.//center`)
	if xPathErr != nil {
		return ""
	}
	text := xPath.String()
	xPath.Free()
	regularExpression := regexp.MustCompile(`\s+`)
	replaceAllString := regularExpression.ReplaceAllString(text, " ")
	trimSpace := strings.TrimSpace(replaceAllString)
	return trimSpace
}

func getTitle(node types.Node) string {
	xPath, xPathErr := node.Find(`.//div/a/text()`)
	if xPathErr != nil {
		return ""
	}
	title := xPath.String()
	xPath.Free()
	return title
}

func getURL(baseURL string, node types.Node) string {
	parse, parseErr := url.Parse(baseURL)
	if parseErr != nil {
		return ""
	}
	xPath, xPathErr := node.Find(`.//div/a/@href`)
	if xPathErr != nil {
		return ""
	}
	text := xPath.String()
	xPath.Free()
	url := fmt.Sprintf("%s://%s%s", parse.Scheme, parse.Host, text)
	return url
}

func getSeeds(node types.Node) int {
	xPath, xPathErr := node.Find(`.//text()`)
	if xPathErr != nil {
		return 0
	}
	text := xPath.String()
	xPath.Free()
	seeds, seedsErr := strconv.Atoi(text)
	if seedsErr != nil {
		return 0
	}
	return seeds
}

func getURLs(document types.Document) []string {
	urls := []string{}

	aXPath, aXPathErr := document.Find(`//div[@class="nfo"]/pre/a`)
	if aXPathErr != nil {
		return urls
	}
	aNodes := aXPath.NodeList()
	aXPath.Free()
	for _, aNode := range aNodes {
		xPath, xPathErr := aNode.Find(`.//@href`)
		if xPathErr != nil {
			continue
		}
		text := xPath.String()
		xPath.Free()
		urls = append(urls, text)
	}

	bXPath, bXPathErr := document.Find(`//a`)
	if bXPathErr != nil {
		return urls
	}
	bNodes := bXPath.NodeList()
	for _, bNode := range bNodes {
		xPath, xPathErr := bNode.Find(`.//@href`)
		if xPathErr != nil {
			continue
		}
		text := xPath.String()
		xPath.Free()
		parse, parseErr := url.Parse(text)
		if parseErr != nil {
			continue
		}
		if parse.Host != "imdb.com" {
			continue
		}
		urls = append(urls, text)
	}

	sort.Strings(urls)

	return urls
}

func getMagnet(document types.Document) string {
	xPath, xPathErr := document.Find(`//div[@class="download"]/a/@href`)
	if xPathErr != nil {
		return ""
	}
	magnet := xPath.String()
	xPath.Free()
	return magnet
}

func getURLsAndMagnet(url string) ([]string, string, error) {
	timeout := time.Duration(15 * time.Second)

	client := http.Client{
		Timeout: timeout,
	}

	response, responseError := client.Get(url)
	if responseError != nil {
		return []string{}, "", errors.New("#1")
	}

	defer response.Body.Close()

	document, documentErr := libxml2.ParseHTMLReader(response.Body)
	if documentErr != nil {
		return []string{}, "", errors.New("#2")
	}

	urls := getURLs(document)

	magnet := getMagnet(document)

	document.Free()

	return urls, magnet, nil
}

func getItems(url string) (Items, error) {
	timeout := time.Duration(15 * time.Second)

	client := http.Client{
		Timeout: timeout,
	}

	response, responseError := client.Get(url)
	if responseError != nil {
		return Items{}, errors.New("#1")
	}

	defer response.Body.Close()

	document, documentErr := libxml2.ParseHTMLReader(response.Body)
	if documentErr != nil {
		return Items{}, errors.New("#2")
	}

	items := Items{}

	trXPath, trXPathErr := document.Find(`//table[@id="searchResult"]/tr`)
	if trXPathErr != nil {
		return Items{}, errors.New("#3")
	}
	trNodes := trXPath.NodeList()
	trXPath.Free()
	for _, trNode := range trNodes {
		tdXPath, tdXPathErr := trNode.Find(`.//td`)
		if tdXPathErr != nil {
			return Items{}, errors.New("#4")
		}
		tdNodes := tdXPath.NodeList()
		tdXPath.Free()
		if len(tdNodes) != 4 {
			continue
		}

		category := getCategory(tdNodes[0])

		title := getTitle(tdNodes[1])

		url := getURL(url, tdNodes[1])
		if url == "" {
			continue
		}

		seeds := getSeeds(tdNodes[2])
		if seeds < 100 {
			continue
		}

		urls, magnet, err := getURLsAndMagnet(url)
		if err != nil {
			continue
		}

		items = append(
			items,
			Item{
				Category: category,
				Title:    title,
				URL:      url,
				Seeds:    seeds,
				URLs:     urls,
				Magnet:   magnet,
			},
		)
	}

	document.Free()

	return items, nil
}

// Query ...
func Query() {
	waitGroup := &sync.WaitGroup{}

	urls := settings.Container.Spiders.URLs
	count := len(urls)

	outgoing := make(chan string)
	incoming := make(chan Items)

	waitGroup.Add(1)
	go consumer(waitGroup, count, outgoing, incoming)

	for index := 1; index <= count; index++ {
		waitGroup.Add(1)
		go producer(waitGroup, outgoing, incoming)
	}

	for _, url := range urls {
		outgoing <- url
	}

	waitGroup.Wait()
}
