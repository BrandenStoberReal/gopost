package main

import (
	"bytes"
	_ "embed"
	"flag"
	"gopost/src/helpers"
	"io"
	"log"
	"net/http"
	url2 "net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var wg sync.WaitGroup

//go:embed res/proxies.txt
var proxies string
var proxiesList []string

func check(e error) {
	if e != nil {

		log.Fatal(e)
	}
}

func logThread(num int, msg string) {
	log.Printf("[THREAD %d]: %s\n", num, msg)
}

func sendPostDataThreaded(url string, payload string, threadNum int, verbose bool, aggressive bool, headerLines []string) {
	defer wg.Done()
	if verbose && !aggressive {
		logThread(threadNum, "New thread created and starting work.")
	}

	for {
		if verbose && !aggressive {
			logThread(threadNum, "URL:>"+url)
		}

		if verbose && !aggressive {
			logThread(threadNum, "JSON:>"+payload)
		}
		req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
		check(err)

		// Set up headers
		for _, header := range headerLines {
			keyValue := strings.Split(header, "=")
			req.Header.Set(keyValue[0], keyValue[1])
		}

		// Set up proxy stuff
		var proxy = proxiesList[helpers.RandomRange(0, len(proxiesList)-1)]
		proxy = strings.Replace(proxy, "\r", "", -1)
		if verbose && !aggressive {
			logThread(threadNum, "Proxy:>"+proxy)
		}
		proxyUrl, err := url2.Parse(proxy)
		check(err)

		// Init client
		client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
		resp, err := client.Do(req)
		if err != nil {
			continue // Normally we would check this, but this usually happens due to bad proxies that are user error. We will just cycle them instead.
		}

		// Log results
		if !aggressive {
			if verbose {
				body, _ := io.ReadAll(resp.Body)
				logThread(threadNum, "Response Status:>"+resp.Status)
				logThread(threadNum, "Response Body:>"+string(body))
			} else {
				logThread(threadNum, "Response Status:>"+resp.Status)
			}
		}

		err = resp.Body.Close()
		check(err)
		client.CloseIdleConnections()
	}

}
func main() {
	urlPtr := flag.String("url", "", "URL to flood with POST requests. This should be the full URL to your endpoint.")
	payloadPtr := flag.String("payload", "", "Path to a payload JSON file to send via POST flooding.")
	threadsPtr := flag.Int("threads", 16, "Number of goroutines to use.")
	headersPtr := flag.String("headers", "", "Path to a txt file of headers to use. Format should be HEADER_NAME=HEADER_VALUE with newlines. Quotes are not added automatically.")
	verbosePtr := flag.Bool("v", false, "Verbose output.")
	speedModePtr := flag.Bool("aggressive", false, "Remove all console logging in order to minimize latency. Not recommended unless you want maximal efficiency.")

	headerLines := make([]string, 0) // Hehe slice

	flag.Parse()

	if *urlPtr == "" {
		log.Fatal("[ERROR] URL is required.")
	}

	if *payloadPtr == "" {
		log.Fatal("[ERROR] Payload is required.")
	}

	// Payload
	absPayloadPath, err := filepath.Abs(*payloadPtr)
	check(err)

	readPayload, err := os.ReadFile(absPayloadPath)
	check(err)

	// Headers
	if *headersPtr != "" {
		headerLines = strings.Split(*headersPtr, "\n")
	}

	log.Println("GoPOST v1.0.1 starting...")
	log.Printf("URL: %s\n", *urlPtr)
	log.Printf("Payload: %s\n", *payloadPtr)
	log.Printf("Threads: %d\n", *threadsPtr)
	log.Printf("Verbose: %t\n", *verbosePtr)
	log.Printf("Aggressive: %t\n", *speedModePtr)

	// Proxies
	proxiesList = strings.Split(proxies, "\n")

	for i := 0; i < *threadsPtr; i++ {
		wg.Add(1)
		go sendPostDataThreaded(*urlPtr, string(readPayload), i+1, *verbosePtr, *speedModePtr, headerLines)
	}
	wg.Wait()
}
