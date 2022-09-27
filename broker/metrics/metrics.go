package rpt

import "github.com/prometheus/client_golang/prometheus"

var PublishSuccessCount = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "Successful_Publish_Requests",
	Help: "number of publish requests which ended as expected",
})

var PublishFailedCount = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "Failed_Publish_Requests",
	Help: "number of publish requests which failed",
})

var SubscribeSuccessCount = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "Successful_Subscribe_Requests",
	Help: "number of subscribe requests that ended as expected",
})

var SubscribeFailedCount = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "Failed_Subscribe_Requests",
	Help: "number of subscribe requests that failed",
})

var FetchSuccessCount = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "Successful_Fetch_Requests",
	Help: "number of successful fetch requests",
})

var FetchFailedCount = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "Failed_Fetch_Requests",
	Help: "number of failed fetch requests",
})

var PublishTime = prometheus.NewSummary(prometheus.SummaryOpts{
	Name:       "Publish_Process_Time",
	Help:       "time spent processing publish request",
	Objectives: map[float64]float64{0.95: 0.01, 0.99: 0.001, 0.5: 0.01},
})

var SubscribeTime = prometheus.NewSummary(prometheus.SummaryOpts{
	Name:       "Subscribe_Process_Time",
	Help:       "time spent processing subscribe request",
	Objectives: map[float64]float64{0.95: 0.01, 0.99: 0.001, 0.5: 0.01},
})

var FetchTime = prometheus.NewSummary(prometheus.SummaryOpts{
	Name:       "Fetch_process_Time",
	Help:       "time spent processing fetch request",
	Objectives: map[float64]float64{0.95: 0.01, 0.99: 0.001, 0.5: 0.01},
})

var TotalSubs = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "Total_Subscribers",
	Help: "number of all subscribers, currently only gets reduced when we publish to the unsubscribed user",
})
