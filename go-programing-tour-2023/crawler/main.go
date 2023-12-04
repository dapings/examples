package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	url := "https://www.thepaper.cn/"
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("fetch url(%s) error: %v", url, err)
		return
	}

	defer func(closer io.ReadCloser) {
		err := closer.Close()
		if err != nil {
			_ = closer.Close()
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Printf("error status code: %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("fail to read content: %v", err)
		return
	}

	log.Printf("body:%s", string(body))
}
