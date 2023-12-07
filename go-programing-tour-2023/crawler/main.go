package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strings"
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

	numLinks := strings.Count(string(body), "<a")
	log.Printf("homepage has %d links!\n", numLinks)

	numLinks = bytes.Count(body, []byte("<a"))
	log.Printf("homepage has %d links!\n", numLinks)

	exist := strings.Contains(string(body), "疫情")
	log.Printf("是否存在疫情:%v\n", exist)

	exist = bytes.Contains(body, []byte("疫情"))
	log.Printf("是否存在疫情:%v\n", exist)

	// log.Printf("body:%s", string(body))
}
