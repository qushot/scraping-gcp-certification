package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type Country int

const (
	English Country = iota
	Japanese
	Spanish
	Portuguese
	French
	German
	Indonesian
)

func (c Country) String() string {
	switch c {
	case English:
		return "English"
	case Japanese:
		return "Japanese"
	case Spanish:
		return "Spanish"
	case Portuguese:
		return "Portuguese"
	case French:
		return "French"
	case German:
		return "German"
	case Indonesian:
		return "Indonesian"
	default:
		return "Unknown"
	}
}

type Certification struct {
	Name              string
	Internationalized map[int]bool
}

func main() {
	url := "https://cloud.google.com/certification/register/?hl=en"

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("http get error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("response status code error")
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalf("goquery new document from reader error: %v", err)
	}

	var certs []*Certification
	doc.Find("tbody").Find("tr").Each(func(i int, tr *goquery.Selection) {
		var cert Certification
		i18ned := map[int]bool{}
		tr.Find("td").Each(func(j int, td *goquery.Selection) {
			if j == 0 {
				cert.Name = td.Find("div").Text()
			} else {
				_, exists := td.Find("img").Attr("src")
				i18ned[j-1] = exists
			}
		})

		cert.Internationalized = i18ned
		certs = append(certs, &cert)
	})

	fmt.Printf("=================================\n")
	for _, cert := range certs {
		fmt.Printf("- %s\n", cert.Name)
		for j := 0; j < len(cert.Internationalized); j++ {
			fmt.Printf("  - %s\t: %v\n", Country(j), cert.Internationalized[j])
		}
		fmt.Printf("=================================\n")
	}
}
