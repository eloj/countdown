#!/bin/bash
YAMLLINTER=yamllint
root=$(realpath $(dirname "${BASH_SOURCE[0]}"))
cd "$root"

go test ./... 2>&1 | tee .gotest.log
if [ $? -ne 0 ]; then
	echo "Test failures. Build aborted."
	exit 1
fi
if [ -x "$(command -v $YAMLLINTER)" ]; then
	$YAMLLINTER config/*.yaml || (echo "Configuration file error, aborted." && exit 2)
	echo "Configuration file validation PASSED."
else
	echo "Configuration file validation skipped, '$YAMLLINTER' not available."
fi
echo "Building..."
GOBUILD="env CGO_ENABLED=0 go build"
$GOBUILD ./cmd/words-cli
$GOBUILD ./cmd/words-server
echo "Build done."
