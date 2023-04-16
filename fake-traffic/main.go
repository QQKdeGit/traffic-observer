package main

import (
	_ "embed"
	"fmt"
	"log"
	"math/rand"
	"net/http"
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

func main() {
	rand.Seed(time.Now().UnixNano())

	for {
		time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)

		var err error
		var path string

		switch rand.Intn(3) {
		case 0:
			path = "http://localhost:8000" + normalPath[rand.Intn(len(normalPath))]
		case 1:
			path = "http://localhost:8000" + attackPath[rand.Intn(len(attackPath))]
		case 2:
			path = "http://localhost:8000" + malwarePath[rand.Intn(len(malwarePath))]
		}

		switch rand.Intn(4) {
		case 0:
			_, err = http.Get(path)
			fmt.Println("GET", path)
		case 1:
			_, err = http.Post(path, "application/json", nil)
			fmt.Println("POST", path)
		case 2:
			req, _ := http.NewRequest("DELETE", path, nil)
			_, err = http.DefaultClient.Do(req)
			fmt.Println("DELETE", path)
		case 3:
			req, _ := http.NewRequest("PUT", path, nil)
			_, err = http.DefaultClient.Do(req)
			fmt.Println("PUT", path)
		}

		if err != nil {
			log.Println(err)
			continue
		}
	}
}
