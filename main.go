package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func parse(url string, selector string) ([]string, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Print(err)
	}
	var result []string
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		title := s.Text()
		result = append(result, title)
	})
	return result, nil
}

func parseHandler(w http.ResponseWriter, req *http.Request) {
	url := req.URL.Query().Get("url")
	if len(url) < 1 {
		fmt.Fprint(w, "Url must not be empty")
		return
	}
	selector := req.URL.Query().Get("selector")
	if len(selector) < 1 {
		fmt.Fprint(w, "Selector must not be empty")
		return
	}
	log.Printf("Extracting %s with selector %s", url, selector)
	parsed, err := parse(url, selector)
	if err != nil {
		fmt.Fprint(w, "Error occured")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(parsed)
}

func main() {
	http.HandleFunc("/", parseHandler)
	fmt.Println("Server started at port 4648")
	http.ListenAndServe(":4648", nil)
}
