package main

import (
	"bareos_exporter/DataAccess/Queries"
	"bareos_exporter/Error"
	"github.com/prometheus/client_golang/prometheus"
)

//Define a struct for you collector that contains pointers
//to prometheus descriptors for each metric you wish to expose.
//Note you can also include fields of other types if they provide utility
//but we just won't be exposing them as metrics.
type BareosMetrics struct {
	TotalFiles *prometheus.Desc
	TotalBytes *prometheus.Desc
	LastJob *prometheus.Desc
	LastFullJob *prometheus.Desc
}

func BareosCollector() *BareosMetrics {
	return &BareosMetrics{
		TotalFiles: prometheus.NewDesc("total_files_saved_for_server",
			"Total files saved for server during all backups combined",
			[]string{"hostname"}, nil,
		),
		TotalBytes: prometheus.NewDesc("total_bytes_saved_for_server",
			"Total bytes saved for server during all backups combined",
			[]string{"hostname"}, nil,
		),
		LastJob: prometheus.NewDesc("last_job_executed_for_server",
			"Last job that was executed for ",
			[]string{"hostname"}, nil,
		),
		LastFullJob: prometheus.NewDesc("last_full_job",
			"Total bytes saved by server",
			[]string{"hostname"}, nil,
		),
	}
}

func (collector *BareosMetrics) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.TotalFiles
	ch <- collector.TotalBytes
	ch <- collector.LastJob
	ch <- collector.LastFullJob
}

func (collector *BareosMetrics) Collect(ch chan<- prometheus.Metric) {
	results, err := db.Query("SELECT DISTINCT Name FROM job WHERE SchedTime LIKE '2019-07-24%'")

	Error.Check(err)

	var servers []Queries.Server

	for results.Next() {
		var server Queries.Server
		err = results.Scan(&server.Name)

		servers = append(servers, server)

		Error.Check(err)
	}

	for _, server := range servers {
		serverFiles := server.TotalFiles(db)

		ch <- prometheus.MustNewConstMetric(collector.TotalFiles, prometheus.CounterValue, float64(serverFiles.Files), server.Name)
	}

	for _, server := range servers {
		serverBytes := server.TotalBytes(db)

		ch <- prometheus.MustNewConstMetric(collector.TotalBytes, prometheus.CounterValue, float64(serverBytes.Bytes), server.Name)
	}
}
