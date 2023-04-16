package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/elazarl/goproxy"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	listenPort, serverPort := 8079, 8080
	flag.IntVar(&listenPort, "p", 8079, "Port to listen on")
	flag.IntVar(&serverPort, "m", 8080, "Port to listen on")
	flag.Parse()

	go func() {
		setHTTPServer(serverPort)
	}()

	setGoProxy(listenPort)
}

func setGoProxy(listenPort int) {
	proxy := goproxy.NewProxyHttpServer()
	// proxy.Verbose = true

	proxy.OnResponse().DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		response := ResponseMetadata{
			Status:           resp.Status,
			StatusCode:       resp.StatusCode,
			Proto:            resp.Proto,      // "HTTP/1.0"
			ProtoMajor:       resp.ProtoMajor, // 1
			ProtoMinor:       resp.ProtoMinor, // 0
			ContentLength:    resp.ContentLength,
			TransferEncoding: resp.TransferEncoding,
			Close:            resp.Close,
			Uncompressed:     resp.Uncompressed,
			Timestamp:        time.Now().UnixMicro(),
		}

		// TransferEncoding may be nil
		if resp.TransferEncoding == nil {
			response.TransferEncoding = []string{}
		}

		// traffic cache for detector
		responseToPrometheus = append(responseToPrometheus, response)

		return resp
	})

	proxy.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		request := RequestMetadata{
			UserAgent:        req.UserAgent(),
			Method:           req.Method,
			Proto:            req.Proto,      // "HTTP/1.0"
			ProtoMajor:       req.ProtoMajor, // 1
			ProtoMinor:       req.ProtoMinor, // 0
			ContentLength:    req.ContentLength,
			TransferEncoding: req.TransferEncoding,
			Close:            req.Close,
			RemoteAddr:       req.RemoteAddr,
			RequestURI:       req.RequestURI,
			Scheme:           req.URL.Scheme,
			Host:             req.URL.Host,
			Path:             req.URL.Path,
			IsMalicious:      -1.0,
			Timestamp:        time.Now().UnixMicro(),
		}

		// TransferEncoding may be nil
		if req.TransferEncoding == nil {
			request.TransferEncoding = []string{}
		}

		// traffic cache for detector
		requestToDetector = append(requestToDetector, request)
		return req, nil
	})

	// set timer to send traffic to detector
	go func() {
		for {
			time.Sleep(5 * time.Second)

			if len(requestToDetector) != 0 {
				requestToPrometheus = append(requestToPrometheus, trafficDetect(requestToDetector)...)
				requestToDetector = nil
			}
		}
	}()

	log.Printf("Proxy listening on port %d", listenPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", listenPort), proxy))
}

func setHTTPServer(serverPort int) {
	reg := prometheus.NewRegistry()

	reg.MustRegister(
		&RequestCollector{},
		&ResponseCollector{},
	)

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))

	log.Printf("Server serving on port %d", serverPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", serverPort), nil))
}

func trafficDetect(trafficList []RequestMetadata) []RequestMetadata {
	arg, err := json.Marshal(trafficList)
	if err != nil {
		log.Println("Marshal error: ", err)
		return nil
	}

	var resp *http.Response

	resp, err = http.Post("http://traffic-detector:8000/detect", "application/json", strings.NewReader(string(arg)))
	if err != nil {
		log.Println("Post to detector error: ", err)
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Read response body error: ", err)
		return nil
	}

	err = json.Unmarshal(body, &trafficList)
	if err != nil {
		log.Println("Unmarshal error: ", err)
		return nil
	}

	return trafficList
}
