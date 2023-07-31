package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/IrokezPetrovic/angie_exporter/angiescraper"
	"github.com/IrokezPetrovic/angie_exporter/metricwriter"
)

var AngieExporterVersion string

type MetricsHandler struct {
	scraper angiescraper.AngieScraper
}

func NewMetricsHandler(scraper angiescraper.AngieScraper) *MetricsHandler {
	h := new(MetricsHandler)
	h.scraper = scraper
	return h
}

func (h *MetricsHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		writer.WriteHeader(405)
		writer.Write([]byte("Method not allowed\n"))
		return
	}

	if req.URL.Path != "/metrics" {
		writer.WriteHeader(404)
		writer.Write([]byte("Not found\n"))
		return
	}
	writer.Header().Add("Content-Type", "prometheus/metrics")
	writer.WriteHeader(200)

	status, err := h.scraper.Scrape()
	metricwriter.WriteMetrics(writer, &status)

	if err != nil {
		fmt.Printf("Scrape error: %s\n", err.Error())

	}

}

func main() {

	if AngieExporterVersion == "" {
		AngieExporterVersion = "develop"
	}

	var scrapeuri = flag.String("scrapeuri", "http://localhost:8080/angie_status", "URL for get status from Angie")
	var listenPort = flag.String("listenport", "9197", "Listen port")
	var listenAddr = flag.String("listenaddr", "0.0.0.0", "Listen address")
	var getVersion = flag.Bool("version", false, "Get version")

	listen := (*listenAddr) + ":" + (*listenPort)
	flag.Parse()

	if *getVersion {
		fmt.Printf("version %s\n", AngieExporterVersion)
		return
	}

	fmt.Printf("Scrape url=%s\n", *scrapeuri)
	fmt.Printf("Listen %s\n", listen)

	scraper := angiescraper.NewAngieScraper(*scrapeuri)
	handler := NewMetricsHandler(scraper)

	http.ListenAndServe(listen, handler)
	for {
		time.Sleep(time.Second)
	}
}
