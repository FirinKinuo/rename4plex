VERSION := $(file < VERSION)

clean:
	rm -rf _build/ release/

build:
	go mod download
	CGO_ENABLED=0 go build -tags release -o anime-replacer cmd/main/main.go

build-all:
	mkdir -p _build
	GOOS=darwin  GOARCH=amd64   CGO_ENABLED=0 go build -tags release -o _build/anime-replacer-$(VERSION)-darwin-amd64 cmd/main/main.go
	GOOS=linux   GOARCH=amd64   CGO_ENABLED=0 go build -tags release -o _build/anime-replacer-$(VERSION)-linux-amd64 cmd/main/main.go
	GOOS=linux   GOARCH=arm     CGO_ENABLED=0 go build -tags release -o _build/anime-replacer-$(VERSION)-linux-arm cmd/main/main.go
	GOOS=linux   GOARCH=arm64   CGO_ENABLED=0 go build -tags release -o _build/anime-replacer-$(VERSION)-linux-arm64 cmd/main/main.go
	GOOS=windows GOARCH=amd64   CGO_ENABLED=0 go build -tags release -o _build/anime-replacer-$(VERSION)-windows-amd64.exe cmd/main/main.go
	cd _build; sha256sum * > sha256sums.txt

install:
	make build
	mv ./anime-replacer /usr/bin/anime-replacer

	mkdir -p /etc/go-plex-anime
	cp ./configs/config.yaml.sample /etc/go-plex-anime/config.yaml

	mkdir -p /var/log/go-plex-anime

release:
	make clean
	make build-all
	mkdir release
	cp _build/* release
	cd release; sha256sum --quiet --check sha256sums.txt && \
	gh release create v$(VERSION) -d -t v$(VERSION) *
