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

type pokemonDownloader struct {
	closeCh chan bool
	urlCh   chan string
}

func (p *pokemonDownloader) Close() {
	log.Println("Close channel")
	close(p.closeCh)
}

func (p *pokemonDownloader) GetImgURL() string {
	return <-p.urlCh
}

func (p *pokemonDownloader) getAllImgURL() {

	res, err := http.Get("https://pokemondb.net/pokedex/national")
	if err != nil {
		log.Fatalf("getImgURLs: http.Get err: %v", err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalf("getImgURLs: goquery.NewDocumentFromReader err: %v", err)
	}

	doc.Find(".img-fixed").Each(func(i int, s *goquery.Selection) {
		url, ok := s.Attr("src")
		if ok {
			p.urlCh <- url
		}
	})

	p.Close()
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

	timeout := make(chan bool)

	go func() {
		time.Sleep(10 * time.Second)
		timeout <- true
	}()

	p := &pokemonDownloader{
		closeCh: make(chan bool),
		urlCh:   make(chan string),
	}

	go func() {
		time.Sleep(6 * time.Second)
		p.Close()
	}()

	go p.getAllImgURL()

outterLoop:
	for {
		select {
		case <-tick.C:
			url := p.GetImgURL()
			downloadImg(url)
		case <-timeout:
			log.Println("timeout")
			tick.Stop()
			break outterLoop
		case <-p.closeCh:
			log.Println("channel closed")
			tick.Stop()
			break outterLoop
		}
	}

	log.Println("Done")
}
