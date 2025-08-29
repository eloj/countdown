#!/bin/bash
YAMLLINTER=yamllint
G_ENV="CGO_ENABLED=0"
GOBUILD="env ${G_ENV} go build"
GOTEST="env ${G_ENV} go test"
root=$(realpath $(dirname "${BASH_SOURCE[0]}"))
cd "$root"

if [ -x "$(command -v $YAMLLINTER)" ]; then
	$YAMLLINTER config/*.yaml
	if [ $? -ne 0 ]; then
		echo "Configuration file error, aborted."
		exit 1
	fi
	echo "Configuration file validation PASSED."
else
	echo "Configuration file validation skipped, '$YAMLLINTER' not available."
fi

echo "Testing..."
${GOTEST} ./... 2>&1 | tee .gotest.log
if [ $? -ne 0 ]; then
	echo "Test failures. Build aborted."
	exit 1
fi

echo "Building..."
${GOBUILD} ./cmd/words-cli
${GOBUILD} ./cmd/words-server

echo "Done."
