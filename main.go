package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func parse(url string, selector string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return "", err
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Print(err)
	}
	var result string
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		title := s.Text()
		result += title + "#"
	})
	return result, nil
}

func parseHandler(w http.ResponseWriter, req *http.Request) {
	url := req.URL.Query().Get("url")
	if len(url) < 1 {
		fmt.Fprint(w, "Url must not be empty")
	}
	selector := req.URL.Query().Get("selector")
	if len(selector) < 1 {
		fmt.Fprint(w, "Selector must not be empty")
	}
	log.Printf("Extracting %s with selector %s", url, selector)
	resp, err := parse(url, selector)
	if err != nil {
		fmt.Fprint(w, "Error occured")
	}
	fmt.Fprint(w, resp)
}

func main() {
	http.HandleFunc("/", parseHandler)
	fmt.Println("Server started at port 4647")
	http.ListenAndServe(":4647", nil)
}
