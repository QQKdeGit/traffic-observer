package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

var responseDesc = prometheus.NewDesc(
	prometheus.BuildFQName("httpproxy", "traffic", "response"),
	"Response information",
	[]string{
		"Status",
		"StatusCode",
		"Proto",
		"ProtoMajor",
		"ProtoMinor",
		"ContentLength",
		"TransferEncoding",
		"Close",
		"Uncompressed",
		"Timestamp",
	},
	nil,
)

type ResponseCollector struct{}

type ResponseMetadata struct {
	Status           string
	StatusCode       int
	Proto            string // "HTTP/1.0"
	ProtoMajor       int    // 1
	ProtoMinor       int    // 0
	ContentLength    int64
	TransferEncoding []string
	Close            bool
	Uncompressed     bool
	Timestamp        int64
}

var responseToPrometheus []ResponseMetadata

func (u *ResponseCollector) Describe(d chan<- *prometheus.Desc) {
	d <- responseDesc
}

func (u *ResponseCollector) Collect(m chan<- prometheus.Metric) {
	log.Printf("Response Sum = %d", len(responseToPrometheus))

	for _, v := range responseToPrometheus {
		m <- prometheus.MustNewConstMetric(
			responseDesc,
			prometheus.GaugeValue,
			1,
			v.Status,
			fmt.Sprintf("%d", v.StatusCode),
			v.Proto,
			fmt.Sprintf("%d", v.ProtoMajor),
			fmt.Sprintf("%d", v.ProtoMinor),
			fmt.Sprintf("%d", v.ContentLength),
			strings.Join(v.TransferEncoding, ","),
			fmt.Sprintf("%v", v.Close),
			fmt.Sprintf("%v", v.Uncompressed),
			fmt.Sprintf("%d", v.Timestamp),
		)
	}

	responseToPrometheus = []ResponseMetadata{}
}
