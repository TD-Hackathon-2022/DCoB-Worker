.PHONY: all
all: clean dev

dev: fmt test

fmt:
	go fmt ./...

build:
	GOOS=js GOARCH=wasm go build -o static/main.wasm

test:
	go test -tags=test ./...

clean:
	go clean -i -r -testcache -cache
	rm -rm static/main.wasm

package: build
	docker build -t dcob/works:latest .