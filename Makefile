GIT_COMMIT := $(shell git rev-parse --short HEAD)
GOBUILD=go build -ldflags "-X github.com/nickysemenza/hyperion/core/config.GitCommit=${GIT_COMMIT}"

build:
	$(GOBUILD)
dev: dev-server
dev-server: build
	./hyperion server
dev-client:
	./hyperion client
dev-ui:
	cd ui && yarn start
test-server: 
	go test -cover  ./...
test-ui:
	cd ui && CI=true yarn test
test: test-server test-ui

proto:
	protoc --go_out=plugins=grpc:api proto/*.proto
godepgraph:
	graphpkg -match 'nickysemenza/hyperion' github.com/nickysemenza/hyperion