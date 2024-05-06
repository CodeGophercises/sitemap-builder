package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/CodeGophercises/html-link-parser/parser"
	"github.com/CodeGophercises/sitemap-builder/encodeXml"
)

var rootPage = flag.String("url", "", "the domain url")
var depth = flag.Int("depth", 10, "max depth of link visiting")

func sameDomain(urlStr string, domain string) bool {
	u, err := url.Parse(urlStr)
	if err != nil {
		log.Fatalf("error parsing url %s", urlStr)
	}
	if !u.IsAbs() {
		return true
	}
	return u.Hostname() == domain
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("[+]%s took %s", name, elapsed)
}

func expandLink(link, scheme, domain string) string {

	link = strings.TrimSuffix(link, "/")

	// We have to expand the link to compare it to visited links
	u, err := url.Parse(link)
	if err != nil {
		log.Fatalf("error parsing url %s", link)
	}
	if !u.IsAbs() {
		link = scheme + "://" + domain + link
	}
	return link
}

func main() {
	// Lets benchmark the thing also
	defer timeTrack(time.Now(), "siteMapBuild")
	flag.Parse()
	if *rootPage == "" {
		log.Fatalf("Will need a url.")
	}

	u, err := url.Parse(*rootPage)
	if err != nil {
		log.Fatalf("error parsing url %s", *rootPage)
	}

	domain := u.Hostname()
	scheme := u.Scheme
	// Visiting each page and then visiting all links in that page wil make a tree. One with possible cycles.
	// To walk the tree , we will have to find a way. I will use BFS as we can use depth later for flexibility.
	// We will also have to maintain a map of visited links so that we don't get stuck in a cycle.

	// A queue to maintain the set of links to be visited
	links := make([]string, 0)

	// Map to maintain visited links
	visited := make(map[string]struct{})

	// Append root page to queue to get started
	rootLink := expandLink(*rootPage, scheme, domain)
	links = append(links, rootLink)
	visited[rootLink] = struct{}{}
	curDepth := 0
	for len(links) > 0 {
		curDepth += 1
		if curDepth > *depth {
			break
		}
		// Pop out the first element from queue
		link := links[0]
		links = links[1:]

		// Find other links on this page, find childrennnnn
		resp, err := http.Get(link)
		defer resp.Body.Close()

		if err != nil {
			log.Fatalf("error reading %s", link)
		}

		htmlData, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("error reading %s", link)
		}

		childLinks, err := parser.Parse(htmlData)
		if err != nil {
			log.Fatalf("error visiting %s", link)
		}

		// Push childLinks in queue
		for _, cl := range childLinks {
			if sameDomain(cl.Href, domain) {
				l := expandLink(cl.Href, scheme, domain)
				if _, ok := visited[l]; !ok {
					links = append(links, l)
					visited[l] = struct{}{}
				}
			}
		}
	}

	// Lets encode them in xml
	finalLinks := make([]string, 0)
	for l, _ := range visited {
		finalLinks = append(finalLinks, l)
	}
	encodeXml.EncodeSitemap(finalLinks)
}
