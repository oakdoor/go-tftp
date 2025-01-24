#!/bin/bash

set -eu

TEST_PORT=10069
TEST_OUTPUT_DIR=/tmp/tftp-output/

rm -rf $TEST_OUTPUT_DIR
mkdir -p $TEST_OUTPUT_DIR

trap 'pkill -f "tftp-server"' exit

echo "starting server..."
./tftp-server --output-folder ${TEST_OUTPUT_DIR} --port ${TEST_PORT} &
sleep 1

echo "sending file..."
./tftp-client --file README.md --windowsize 42 --blocksize 1200 --single-port 10068 --retransmit 4 --timeout 2 tftp://localhost:${TEST_PORT}/sent-README.md

echo "checking file received..."
diff README.md ${TEST_OUTPUT_DIR}sent-README.md
echo "test passed"
