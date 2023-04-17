package main

import (
	_ "embed"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	//go:embed normal.csv
	normal string
	//go:embed attack.csv
	attack string
	//go:embed malware.csv
	malware string
)

var (
	normalPath  []string
	attackPath  []string
	malwarePath []string
)

func splitCSV(s string) []string {
	var ret []string
	for _, i := range strings.Split(s, "\n") {
		if i == "" {
			continue
		}

		i = strings.ReplaceAll(i, "\r", "")
		i = strings.ReplaceAll(i, "\n", "")

		ret = append(ret, i)
	}
	return ret
}

func init() {
	normalPath = splitCSV(normal)
	attackPath = splitCSV(attack)
	malwarePath = splitCSV(malware)
}

var root = "http://172.20.62.49:8001"

func main() {
	rand.Seed(time.Now().UnixNano())

	for {
		time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)

		var err error
		var path string

		switch rand.Intn(3) {
		case 0:
			path = root + normalPath[rand.Intn(len(normalPath))]
		case 1:
			path = root + attackPath[rand.Intn(len(attackPath))]
		case 2:
			path = root + malwarePath[rand.Intn(len(malwarePath))]
		}

		proxy := func(_ *http.Request) (*url.URL, error) {
			return url.Parse("http://localhost:8079")
		}
		transport := &http.Transport{Proxy: proxy}
		client := &http.Client{Transport: transport}

		switch rand.Intn(4) {
		case 0:
			_, err = client.Get(path)
			fmt.Println("\033[38;5;82mGET\033[0m", path)
		case 1:
			_, err = client.Post(path, "application/json", nil)
			fmt.Println("\033[38;5;81mPOST\033[0m", path)
		case 2:
			req, _ := http.NewRequest("DELETE", path, nil)
			_, err = client.Do(req)
			fmt.Println("\033[38;5;196mDELETE\033[0m", path)
		case 3:
			req, _ := http.NewRequest("PUT", path, nil)
			_, err = client.Do(req)
			fmt.Println("\033[38;5;220mPUT\033[0m", path)
		}

		if err != nil {
			log.Println(err)
			continue
		}
	}
}
