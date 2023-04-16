package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

var requestDesc = prometheus.NewDesc(
	prometheus.BuildFQName("httpproxy", "traffic", "request"),
	"Request information",
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
		"Timestamp",
	},
	nil,
)

type RequestCollector struct{}

type RequestMetadata struct {
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
	Timestamp        int64
}

var requestToPrometheus []RequestMetadata
var requestToDetector []RequestMetadata

func (u *RequestCollector) Describe(d chan<- *prometheus.Desc) {
	d <- requestDesc
}

func (u *RequestCollector) Collect(m chan<- prometheus.Metric) {
	log.Printf("Request Sum = %d", len(requestToPrometheus))

	for _, v := range requestToPrometheus {
		m <- prometheus.MustNewConstMetric(requestDesc, prometheus.GaugeValue, float64(len(requestToPrometheus)),
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
			fmt.Sprintf("%d", v.Timestamp),
		)
	}

	requestToPrometheus = nil
}
