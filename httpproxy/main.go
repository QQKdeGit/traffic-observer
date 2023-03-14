package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/elazarl/goproxy"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// var (
// 	testTraffic = TrafficMetadata{
// 		UserAgent:        "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/110.0",
// 		Method:           "GET",
// 		Proto:            "HTTP/1.1",
// 		ProtoMajor:       1,
// 		ProtoMinor:       1,
// 		ContentLength:    0,
// 		TransferEncoding: []string{},
// 		Close:            false,
// 		RemoteAddr:       "172.20.48.1:50743",
// 		RequestURI:       "http://www.yiidian.com/questions/30562",
// 		Scheme:           "http",
// 		Host:             "www.yiidian.com",
// 		Path:             "/questions/30562",
// 		IsMalicious:      0,
// 	}
// )

var TEST_MODE = true

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

	proxy.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		traffic := TrafficMetadata{
			UserAgent:     req.UserAgent(),
			Method:        req.Method,
			Proto:         req.Proto,      // "HTTP/1.0"
			ProtoMajor:    req.ProtoMajor, // 1
			ProtoMinor:    req.ProtoMinor, // 0
			ContentLength: req.ContentLength,
			// TransferEncoding: req.TransferEncoding,
			TransferEncoding: []string{},
			Close:            req.Close,
			RemoteAddr:       req.RemoteAddr,
			RequestURI:       req.RequestURI,
			Scheme:           req.URL.Scheme,
			Host:             req.URL.Host,
			Path:             req.URL.Path,
			IsMalicious:      -1.0,
		}

		fmt.Printf("%#v\n", traffic)

		// trafficArray = append(trafficArray, trafficDetect(traffic))
		trafficArray = append(trafficArray, traffic)

		return req, nil
	})

	log.Printf("Proxy listening on port %d", listenPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", listenPort), proxy))
}

var trafficDesc = prometheus.NewDesc(
	prometheus.BuildFQName("httpproxy", "traffic", "info"),
	"Labeled system information as provided by the uname system call.",
	[]string{
		"UserAgent",
		"Method",
		"Proto",
		"ProtoMajor",
		"ProtoMinor",
		"ContentLength",
		"TransferEncoding",
		"Close",
		"RemoteAddr",
		"RequestURI",
		"Scheme",
		"Host",
		"Path",
		"IsMalicious",
	},
	nil,
)

type TrafficCollector struct{}

type TrafficMetadata struct {
	UserAgent        string
	Method           string
	Proto            string // "HTTP/1.0"
	ProtoMajor       int    // 1
	ProtoMinor       int    // 0
	ContentLength    int64
	TransferEncoding []string
	Close            bool
	RemoteAddr       string
	RequestURI       string
	Scheme           string
	Host             string
	Path             string
	IsMalicious      float64 // -1 = undefined, 0 = benign, 1 = malicious,
}

var trafficArray []TrafficMetadata

func (u *TrafficCollector) Describe(d chan<- *prometheus.Desc) {
	d <- trafficDesc
}

func (u *TrafficCollector) Collect(m chan<- prometheus.Metric) {
	log.Printf("Traffic Sum = %d", len(trafficArray))

	for _, v := range trafficArray {
		m <- prometheus.MustNewConstMetric(trafficDesc, prometheus.GaugeValue, float64(len(trafficArray)),
			v.UserAgent,
			v.Method,
			v.Proto,
			fmt.Sprintf("%d", v.ProtoMajor),
			fmt.Sprintf("%d", v.ProtoMinor),
			fmt.Sprintf("%d", v.ContentLength),
			strings.Join(v.TransferEncoding, ","),
			fmt.Sprintf("%v", v.Close),
			v.RemoteAddr,
			v.RequestURI,
			v.Scheme,
			v.Host,
			v.Path,
			fmt.Sprintf("%f", v.IsMalicious),
		)
	}

	trafficArray = nil
}

func setHTTPServer(serverPort int) {
	reg := prometheus.NewRegistry()

	reg.MustRegister(&TrafficCollector{})

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))

	log.Printf("Server serving on port %d", serverPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", serverPort), nil))
}

func trafficDetect(trafficList []TrafficMetadata) []TrafficMetadata {
	arg, err := json.Marshal(trafficList)
	if err != nil {
		log.Fatal("error:", err)
	}

	var resp *http.Response

	if TEST_MODE {
		resp, err = http.Post("http://localhost:8000/test", "application/json", strings.NewReader(string(arg)))
		if err != nil {
			log.Fatal("error:", err)
		}
	} else {
		resp, err = http.Post("http://localhost:8000/detect", "application/json", strings.NewReader(string(arg)))
		if err != nil {
			log.Fatal("error:", err)
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("error:", err)
	}

	err = json.Unmarshal(body, &trafficList)
	if err != nil {
		log.Fatal("error:", err)
	}

	return trafficList
}