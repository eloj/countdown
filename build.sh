#!/bin/bash
echo "Running tests"
go test ./...
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
echo "Building cmds..."
go build ./cmd/words-cli
go build ./cmd/words-server
echo "Build done."
