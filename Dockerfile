
FROM golang:1.10 AS builder-server

# Download and install the latest release of dep
ADD https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 /usr/bin/dep
RUN chmod +x /usr/bin/dep

# Copy the code from the host and compile it
WORKDIR $GOPATH/src/github.com/nickysemenza/hyperion
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure --vendor-only
COPY . ./
RUN make build
# move output binary to root so next stage can grab it more cleanly
RUN cp hyperion /

FROM alpine
COPY --from=builder-server /hyperion ./hyperion
COPY config-example.yaml config.yaml
COPY lights.yaml lights.yaml
EXPOSE 8080 
ENTRYPOINT ["./hyperion","server"]