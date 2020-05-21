COMMIT?=$(shell git rev-parse HEAD)
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

all: cli-test
	

cli-test: *.go
	go build -mod vendor -ldflags="-X main.COMMIT=$(COMMIT) -X main.BUILD_TIME=$(BUILD_TIME)" -o cli-test

clean:
	rm -f cli-test

check:
	go test -mod vendor
