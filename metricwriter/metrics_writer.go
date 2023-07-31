package metricwriter

import (
	"fmt"
	"io"

	"github.com/IrokezPetrovic/angie_exporter/angiescraper"
)

func boolToInt(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}

func WriteMetrics(writer io.Writer, status *angiescraper.AngieStatus) {
	fmt.Fprintf(writer, "angie_up %d\n", boolToInt(status.Up))

	if !status.Up {
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
