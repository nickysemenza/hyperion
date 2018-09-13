dev:
	go build && ./hyperion server
dev-client:
	go build && ./hyperion client
dev-ui:
	cd ui && yarn start
test-backend: 
	go test -cover  ./...
test-ui:
	cd ui && CI=true yarn test
test: test-backend test-ui

proto:
	protoc --go_out=plugins=grpc:api proto/*.proto
godepgraph:
	graphpkg -match 'nickysemenza/hyperion' github.com/nickysemenza/hyperion