package angiescraper

type UpstreamInfo struct {
	Peers map[string]UpstreamPeep
}

type UpstreamPeep struct {
	Server    string
	Backup    bool
	Weight    uint64
	State     string
	Responses map[string]uint64
	Data      UpstreamPeerDataStats
	Health    UpstreamPeerHealth
}

type UpstreamPeerDataStats struct {
	Sent     uint64
	Received uint64
}

type UpstreamPeerHealth struct {
	Fails       uint64
	Unavailable uint64
	Downtime    uint64
}
