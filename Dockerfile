
FROM golang:1.11 AS builder-server

# Copy the code from the host and compile it
WORKDIR /src/hyperion
COPY . ./
RUN make build
# move output binary to root so next stage can grab it more cleanly
RUN cp hyperion /

FROM node:8 as builder-ui

WORKDIR /app 
ADD ui/ ./

ENV NODE_PATH=/node_modules
ENV PATH=$PATH:/node_modules/.bin
RUN yarn
RUN yarn run build

FROM alpine
COPY --from=builder-ui /app/build ./ui/build
COPY --from=builder-server /hyperion ./hyperion
ENTRYPOINT ["./hyperion","server"]