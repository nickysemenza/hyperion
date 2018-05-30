#!/bin/bash
protoc \
    --go_out=plugins=grpc:backend \
    proto/*.proto
