#!/bin/bash
root=$(realpath $(dirname "${BASH_SOURCE[0]}"))
cd "$root"

go test ./... 2>&1 | tee .gotest.log
if [ $? -ne 0 ]; then
	echo "Test failures. Build aborted."
	exit 1
fi
if [ -x "$(command -v yamllint)" ]; then
	yamllint config/*.yaml || (echo "Configuration file error, aborted." && exit 2)
	echo "Configuration file validation PASSED."
else
	echo "Configuration file validation skipped, 'yamllint' not available."
fi
echo "Building..."
GOBUILD="env CGO_ENABLED=0 go build"
$GOBUILD ./cmd/words-cli
$GOBUILD ./cmd/words-server
echo "Build done."
