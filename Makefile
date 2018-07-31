dev:
	go build && ./hyperion server
test-backend: 
	go test ./...
test-ui:
	cd ui && CI=true yarn test
test: test-backend test-ui