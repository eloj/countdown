#!/bin/bash
go test ./...
if [ $? -ne 0 ]; then
	echo "Test failures. Build aborted."
	exit 1
fi
go build ./cmd/words-cli
