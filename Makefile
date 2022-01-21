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
	rm -rf static/main.wasm

package: build
	docker build -t hydezhao/worker:latest .

push: package
	docker push  hydezhao/worker:latest

run: package
	docker run --name=worker01 -d -p 8081:80 hydezhao/worker:latest

slim:
	tinygo build -opt z -o ./static/tinygo.wasm -target wasm ./main.go