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
	failedJobs = promauto.NewCounter(prometheus.CounterOpts{
		Name: "failed_jobs",
		Help: "The total number of analysis jobs that failed",
	})
	unsuccessfulJobs = promauto.NewCounter(prometheus.CounterOpts{
		Name: "unsuccessful_jobs",
		Help: "The total number of analysis jobs that finished but were not successful",
	})
)
