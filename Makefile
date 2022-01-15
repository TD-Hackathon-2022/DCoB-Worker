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
	docker build -t hydezhao/worker:latest .

push: package
	docker push  hydezhao/worker:latest

run: package
	docker run --name=worker01 -d -p 8080:80 hydezhao/worker:latest