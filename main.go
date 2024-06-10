package main

import (
	"bufio"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"golang.org/x/net/html"
)

const (
	banner = `
	 _       ___   _   _   _  __   ____  
	| |     |_ _| | \ | | | |/ /  / ___| 
	| |      | |  |  \| | | ' /  | |  _  
	| |___   | |  | |\  | | . \  | |_| | 
	|_____| |___| |_| \_| |_|\_\  \____| 
									   
`
)

func robots(urlStr, output string) {
	robotsURL := urlStr + "/robots.txt"
	resp, err := http.Get(robotsURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Println("Failed to retrieve robots.txt from", urlStr)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading robots.txt:", err)
		return
	}

	re := regexp.MustCompile(`Disallow: (.+)`)
	matches := re.FindAllStringSubmatch(string(body), -1)

	file, err := os.Create(output)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, match := range matches {
		disallowedURL := urlStr + match[1]
		fmt.Fprintln(writer, disallowedURL)
	}
	writer.Flush()
	fmt.Println("Extracted URLs from robots.txt saved in", output)
}

func sitemap(urlStr, output string) {
	sitemapURL := urlStr + "/sitemap.xml"
	resp, err := http.Get(sitemapURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Println("Failed to retrieve sitemap.xml from", urlStr)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading sitemap.xml:", err)
		return
	}

	type Loc struct {
		Loc string `xml:"loc"`
	}
	type Urlset struct {
		URLs []Loc `xml:"url"`
	}

	var urlset Urlset
	err = xml.Unmarshal(body, &urlset)
	if err != nil {
		fmt.Println("Error parsing sitemap.xml:", err)
		return
	}

	file, err := os.Create(output)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, loc := range urlset.URLs {
		fmt.Fprintln(writer, loc.Loc)
	}
	writer.Flush()
	fmt.Println("Extracted URLs from sitemap.xml saved in", output)
}

func extractLinks(urlStr, output string) {
	resp, err := http.Get(urlStr)
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Println("Failed to retrieve page from", urlStr)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading page:", err)
		return
	}

	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		fmt.Println("Error parsing HTML:", err)
		return
	}

	file, err := os.Create(output)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					fmt.Fprintln(writer, attr.Val)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	writer.Flush()
	fmt.Println("Extracted links from page saved in", output)
}

func main() {
	fmt.Fprintln(os.Stdout, []any{banner}...)
	urlStr := flag.String("url", "", "Website URL")
	output := flag.String("output", "", "Output file name prefix")
	robotsFlag := flag.Bool("robots", false, "Extract URLs from robots.txt")
	sitemapFlag := flag.Bool("sitemap", false, "Extract URLs from sitemap.xml")
	linksFlag := flag.Bool("links", false, "Extract all links from the webpage")
	flag.Parse()

	if *urlStr == "" || *output == "" {
		fmt.Println("Website URL and output file name prefix are required.")
		flag.Usage()
		return
	}

	var wg sync.WaitGroup

	if *robotsFlag {
		wg.Add(1)
		go func() {
			defer wg.Done()
			outputFile := *output + "_robots.txt"
			robots(*urlStr, outputFile)
		}()
	}

	if *sitemapFlag {
		wg.Add(1)
		go func() {
			defer wg.Done()
			outputFile := *output + "_sitemap.txt"
			sitemap(*urlStr, outputFile)
		}()
	}

	if *linksFlag {
		wg.Add(1)
		go func() {
			defer wg.Done()
			outputFile := *output + "_links.txt"
			extractLinks(*urlStr, outputFile)
		}()
	}

	if !*robotsFlag && !*sitemapFlag && !*linksFlag {
		fmt.Println("At least one of --robots, --sitemap, or --links options must be specified.")
		flag.Usage()
		return
	}

	wg.Wait()
}
