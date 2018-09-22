GIT_COMMIT := $(shell git rev-parse --short HEAD)
GOBUILD= CGO_ENABLED=0 go build -ldflags "-X github.com/nickysemenza/hyperion/core/config.GitCommit=${GIT_COMMIT}" -o hyperion

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

IMAGE=nicky/hyperion

docker-build:
	docker build -t $(IMAGE) .
docker-run:
	docker run -p 8080:8080 -e "A=b" $(IMAGE) 