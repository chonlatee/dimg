package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func getImgURLs() []string {

	res, err := http.Get("https://pokemondb.net/pokedex/national")
	if err != nil {
		log.Fatalf("getImgURLs: http.Get err: %v", err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalf("getImgURLs: goquery.NewDocumentFromReader err: %v", err)
	}

	var urls []string

	doc.Find(".img-fixed").Each(func(i int, s *goquery.Selection) {
		url, ok := s.Attr("src")
		if ok {
			urls = append(urls, url)
		}
	})

	return urls
}

func downloadImg(url string) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatalf("downloadImg: http.Get err: %v", err)
	}
	defer res.Body.Close()

	urs := strings.Split(url, "/")
	imgName := urs[len(urs)-1]

	f, err := os.Create("img/" + imgName)
	if err != nil {
		log.Fatalf("downloadImg: os.Create err: %v", err)
	}

	_, err = io.Copy(f, res.Body)
	if err != nil {
		log.Fatalf("downloadImg: io.Copy err: %v", err)
	}

	log.Printf("downloadImg: download %s success\n", imgName)
}

func main() {
	tick := time.NewTicker(1 * time.Second)

	urls := getImgURLs()

	i := 0

outterLoop:
	for {
		select {
		case <-tick.C:
			downloadImg(urls[i])
			i += 1
			if i >= 5 {
				tick.Stop()
				break outterLoop
			}
		}
	}

	log.Println("Done")
}
