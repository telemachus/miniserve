.DEFAULT_GOAL := build

fmt:
	golangci-lint run --disable-all --no-config -Egofmt --fix
	golangci-lint run --disable-all --no-config -Egofumpt --fix

lint: fmt
	staticcheck ./...
	revive -config revive.toml ./...
	golangci-lint run

golangci: fmt
	golangci-lint run

staticcheck: fmt
	staticcheck ./...

revive: fmt
	revive -config revive.toml ./...

build: lint
	go build .

install: build
	go install .

clean:
	$(RM) miniserve
	go clean -i -r -cache

.PHONY: fmt lint build install clean
