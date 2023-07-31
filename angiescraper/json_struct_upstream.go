package angiescraper

type UpstreamInfo struct {
	Peers map[string]UpstreamPeer
}

type UpstreamPeer struct {
	Server    string
	Backup    bool
	Weight    uint64
	State     string
	Selected  UpstreamPeerSelected
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

type UpstreamPeerSelected struct {
	Current uint64
	Total   uint64
}
