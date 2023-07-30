package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/IrokezPetrovic/angie_exporter/angiescraper"
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

func boolToInt(b bool) uint8 {
	if b {
		return 1
	}
	return 0
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
	fmt.Fprintf(writer, "angie_up=%d\n", boolToInt(status.Up))

	if err != nil {
		fmt.Printf("Scrape error: %s\n", err.Error())
		return
	}

	for upstreamName, upstream := range status.Http.Upstreams {
		for peerName, peer := range upstream.Peers {
			fmt.Fprintf(writer, "angie_upstream_peer_up{upstream=\"%s\",peer=\"%s\",peer_server=\"%s\"} %d\n",
				upstreamName, peerName, peer.Server, boolToInt(peer.State == "up"))
			fmt.Fprintf(writer, "angie_upstream_peer_data_sent{upstream=\"%s\",peer=\"%s\",peer_server=\"%s\"} %d\n",
				upstreamName, peerName, peer.Server, peer.Data.Sent)
			fmt.Fprintf(writer, "angie_upstream_peer_data_received{upstream=\"%s\",peer=\"%s\",peer_server=\"%s\"} %d\n",
				upstreamName, peerName, peer.Server, peer.Data.Received)

			fmt.Fprintf(writer, "angie_upstream_peer_health_fails{upstream=\"%s\",peer=\"%s\",peer_server=\"%s\"} %d\n",
				upstreamName, peerName, peer.Server, peer.Health.Fails)
			fmt.Fprintf(writer, "angie_upstream_peer_health_unavailable{upstream=\"%s\",peer=\"%s\",peer_server=\"%s\"} %d\n",
				upstreamName, peerName, peer.Server, peer.Health.Unavailable)
			fmt.Fprintf(writer, "angie_upstream_peer_health_downtime{upstream=\"%s\",peer=\"%s\",peer_server=\"%s\"} %d\n",
				upstreamName, peerName, peer.Server, peer.Health.Downtime)

			for statusCode, count := range peer.Responses {
				fmt.Fprintf(writer, "angie_upstream_peer_responses{upstream=\"%s\",peer=\"%s\",peer_server=\"%s\",status=\"%s\"} %d\n",
					upstreamName, peerName, peer.Server, statusCode, count)
			}

			fmt.Fprintf(writer, "\n\n")

		}

	}

	for serverName, serverInfo := range status.Http.ServerZonse {
		fmt.Fprintf(writer, "angie_server_ssl_handshaked{server=\"%s\"} %d\n",
			serverName, serverInfo.Ssl.Handshaked)
		fmt.Fprintf(writer, "angie_server_ssl_reuses{server=\"%s\"} %d\n",
			serverName, serverInfo.Ssl.Reuses)
		fmt.Fprintf(writer, "angie_server_ssl_timedout{server=\"%s\"} %d\n",
			serverName, serverInfo.Ssl.Timedout)
		fmt.Fprintf(writer, "angie_server_ssl_failed{server=\"%s\"} %d\n",
			serverName, serverInfo.Ssl.Failed)
		fmt.Fprintln(writer, "")

		fmt.Fprintf(writer, "angie_server_requests_total{server=\"%s\"} %d\n",
			serverName, serverInfo.Requests.Total)
		fmt.Fprintf(writer, "angie_server_requests_processing{server=\"%s\"} %d\n",
			serverName, serverInfo.Requests.Processing)
		fmt.Fprintf(writer, "angie_server_requests_discarded{server=\"%s\"} %d\n",
			serverName, serverInfo.Requests.Discarded)
		fmt.Fprintln(writer, "")

		fmt.Fprintf(writer, "angie_server_data_received{server=\"%s\"} %d\n",
			serverName, serverInfo.Data.Received)
		fmt.Fprintf(writer, "angie_server_data_sent{server=\"%s\"} %d\n",
			serverName, serverInfo.Data.Sent)
		fmt.Fprintln(writer, "")

		for statusCode, count := range serverInfo.Responses {
			fmt.Fprintf(writer, "angie_server_responses{server=\"%s\",status=\"%s\"} %d\n",
				serverName, statusCode, count)
		}

		fmt.Fprintln(writer, "")
		fmt.Fprintln(writer, "")
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
