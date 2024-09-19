package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	reqsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "processed_requests",
		Help: "The total number of processed http requests",
	})
)
