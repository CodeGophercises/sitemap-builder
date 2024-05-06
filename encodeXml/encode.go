package encodeXml

import (
	"encoding/xml"
	"fmt"
	"log"
)

type URLSet struct {
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`
	URLs    []URL    `xml:"url"`
}

type URL struct {
	Loc string `xml:"loc"`
}

func EncodeSitemap(links []string) {

	urls := make([]URL, len(links))
	for i, link := range links {
		urls[i] = URL{link}
	}

	header := xml.Header
	urlSet := URLSet{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  urls,
	}

	xmlData, err := xml.MarshalIndent(urlSet, "", "  ")
	if err != nil {
		log.Fatalf("error encoding in xml")
	}

	fmt.Print(header)
	fmt.Println(string(xmlData))
}
