package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type Element struct {
	Tag   string            `json:"tag"`
	Text  string            `json:"text"`
	Attrs map[string]string `json:"attrs"`
}

func fetchHTML(url string) (*goquery.Document, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func extractElements(doc *goquery.Document) []Element {
	var elements []Element

	doc.Find("*").Each(func(i int, s *goquery.Selection) {
		element := Element{
			Tag:   goquery.NodeName(s),
			Text:  s.Text(),
			Attrs: make(map[string]string),
		}

		// Get attributes
		for _, attr := range s.Nodes[0].Attr {
			element.Attrs[attr.Key] = attr.Val
		}

		elements = append(elements, element)
	})

	return elements
}

func main() {
	url := "https://www.example.com" // Replace with your target URL

	doc, err := fetchHTML(url)
	if err != nil {
		log.Fatal(err)
	}

	elements := extractElements(doc)

	jsonData, err := json.MarshalIndent(elements, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(jsonData))
}
