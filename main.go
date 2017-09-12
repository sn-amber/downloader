package main

import (
	"bufio"
	"flag"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"syscall"
)

var (
	url string
	to  string
)

func main() {
	flag.StringVar(&url, "url", "", "full output URL of your datacube processing")
	flag.StringVar(&to, "to", ".", "path where you want files to be downloaded")
	flag.Parse()

	if url == "" {
		log.Fatalln("Can't retrieve empty URL")
	}

	username, password := credentials()

	links := getFilelist(url, username, password)
	getFiles(url, links, to, username, password)
}

func getFiles(url string, links []string, to, username, password string) {

	for _, link := range links {

		// create output file
		filepath := path.Join(to, link)
		linkurl := url + "/" + link

		out, err := os.Create(filepath)
		if err != nil {
			log.Fatal("cannot create file", filepath)
		}
		defer out.Close()

		resp, err := createBasicAuthGetRequest(linkurl, username, password)
		if err != nil {
			log.Fatal("can't fetch ", linkurl, err)
		}
		// defer resp.Body.Close()

		fmt.Println("downloading", filepath)
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			log.Fatal("cannot download file", linkurl)
		}

		resp.Body.Close()
	}
}

func getLinksFromResponse(body io.Reader) (links []string) {
	z := html.NewTokenizer(body)

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			return links
		case tt == html.StartTagToken:
			t := z.Token()

			isAnchor := t.Data == "a"
			if isAnchor {
				for _, a := range t.Attr {
					if a.Key == "href" {
						links = append(links, a.Val)
						break
					}
				}
			}
		}
	}
}

func retrievalError() {
	log.Fatal("Couldn't retrieve ", url)
}

func getFilelist(url, username, password string) []string {

	resp, err := createBasicAuthGetRequest(url, username, password)
	if err != nil {
		retrievalError()
	}
	defer resp.Body.Close()

	links := getLinksFromResponse(resp.Body)

	return links
}

// don't forget to close the io.ReadCloser
func createBasicAuthGetRequest(url, username, password string) (*http.Response, error) {
	client := &http.Client{}
	getRequest, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		retrievalError()
	}
	getRequest.SetBasicAuth(username, password)

	resp, err := client.Do(getRequest)
	return resp, err
}

func credentials() (string, string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Username: ")
	username, _ := reader.ReadString('\n')

	fmt.Print("Enter Password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Fatal("couldn't read password")
	}
	password := string(bytePassword)
	fmt.Println()

	return strings.TrimSpace(username), strings.TrimSpace(password)
}
