dev:
	go build && ./hyperion server
test-backend: 
	go test -cover  ./...
test-ui:
	cd ui && CI=true yarn test
test: test-backend test-ui

proto:
	protoc --go_out=plugins=grpc:api proto/*.proto