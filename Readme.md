# Ussage
## Configure ANGIE
```
server {
    listen 8080;
    access_log off;

    location  /angie_status/ {
        api /status/;
        allow 127.0.0.1;
        deny all;
    }
}
```
## Run exporter
```bash
angie_exporter --listenaddr=0.0.0.0 --listenport 9197 --scrapeurl=http://localhost:8080/angie_status
```
Params:
* `--listenaddr`: Address, listened by exporter. Default `0.0.0.0`
* `--listenport`: Post, listened by exporter. Default `9197`
* `--scrapeurl`: URL, provides status of ANGIE. Default `http://localhost:8080/angie_status`

# Build
## Build specific
```bash
go build -ldflags "-X main.AngieExporterVersion=1.0.0" angie_exporter.go
```

## Build dist
```bash
make dist
```

