GO_ENV := CGO_ENABLED=0 GOOS=linux GOARCH=amd64
ROOT=$(shell pwd)
DOCKER_BUILD=$(shell pwd)/.docker_build
MAIN_BIN=$(DOCKER_BUILD)/midash

build: clean
	mkdir -p $(DOCKER_BUILD)
	$(GO_ENV) go build -v -o $(MAIN_BIN) ./cmd/midash

docker: build
	docker build -t $(USER)/midash ./

dockertest: build
	docker build -f test.dockerfile -t $(USER)/midash-test ./
	docker rmi $(USER)/midash-test

clean:
	rm -rf MAIN_BIN