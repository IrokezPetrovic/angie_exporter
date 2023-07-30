package angiescraper

type AngieStatus struct {
	Up            bool
	ScrapeSuccess bool
	Angie         Angie
	Http          Http
}

type Angie struct {
	Version string
}

type Http struct {
	ServerZonse map[string]ServerInfo   `json:"server_zones"`
	Upstreams   map[string]UpstreamInfo `json:"upstreams"`
}
