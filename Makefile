BUILD_VERSION ?= develop

slash-split = $(word $2,$(subst /, ,$1))

dist: clean build/linux/amd64 build/linux/386 
clean:
	@rm -rf build/	
	@rm -rf dist/

build/%: 
	@echo Build Version=$(BUILD_VERSION) Platform=$*
	GOOS=$(call slash-split,$*,1) GOARCH=$(call slash-split,$*,2) go build -ldflags "-X main.AngieExporterVersion=$(BUILD_VERSION)" -o build/$(call slash-split,$*,1)-$(call slash-split,$*,2)/angie_exporter angie_exporter.go
	@mkdir dist/ || true
	tar -czf dist/angie_exporter.$(call slash-split,$*,1)-$(call slash-split,$*,2).tgz -C "build/$(call slash-split,$*,1)-$(call slash-split,$*,2)" "."
	

	