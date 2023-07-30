package angiescraper

type ServerInfo struct {
	Ssl       ServerInfoSsl
	Requests  ServerInfoRequests
	Responses map[string]uint64
	Data      ServerInfoData
}

type ServerInfoSsl struct {
	Handshaked uint64
	Reuses     uint64
	Timedout   uint64
	Failed     uint64
}

type ServerInfoRequests struct {
	Total      uint64
	Processing uint64
	Discarded  uint64
}

type ServerInfoData struct {
	Received uint64
	Sent     uint64
}
