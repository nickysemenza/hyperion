GIT_COMMIT := $(shell git rev-parse --short HEAD)
GOBUILD= CGO_ENABLED=0 go build -ldflags "-X github.com/nickysemenza/hyperion/core/config.GitCommit=${GIT_COMMIT}" -o hyperion

build:
	$(GOBUILD)
dev: dev-server
dev-server: build
	./hyperion server
dev-client: build
	./hyperion client
dev-ui:
	cd ui && yarn run start
test-server: 
	go test -v -race -cover  ./...
test-ui:
	cd ui && CI=true yarn test
test: test-server test-ui

lint-server:
	revive -config revive.toml -formatter friendly -exclude=vendor/... ./...
generate-proto:
	protoc --go_out=plugins=grpc:api proto/*.proto
godepgraph:
	graphpkg -match 'nickysemenza/hyperion' github.com/nickysemenza/hyperion

IMAGE=nicky/hyperion

docker-build:
	docker build -t $(IMAGE) .
docker-run:
	docker run -p 8080:8080 -e "A=b" $(IMAGE) 

# from: https://goldfishtips.wordpress.com/2018/03/17/cool-makefile-target-for-golang-mocks-generation/
update-mocks:
	go list -f '{{.Dir}}' ./... \
    | grep -v "$(notdir $(CURDIR))$$" \
	| grep -v "api" \
    | xargs -n1 ${GOPATH}/bin/mockery \
	-inpkg -case "underscore" -all \
	-note "NOTE: run 'make update-mocks' to regenerate" \
	-dir

# start openlighting server with:
#	UI http://localhost:9092/ola.html
# 	Server on :9010
ola-dev: 
	olad -p 9092
ola-patch:
	ola_patch --device 10 --port 0 --universe 1