#!/bin/bash

cd `dirname $0`
cd ..

go build -o "example/exec/fakedriver" "vendor/github.com/cloudfoundry-incubator/volman/fakedriver/cmd/fakedriver/main.go"

#=======================================================================================================================
# fakedriver runs in 4 different modes to test the 4 different transports we support.  This script tests all 4
#=======================================================================================================================

# UNIX SOCKET TESTS
export FIXTURE_FILENAME=example/fixtures/certification_unix.json
/bin/bash example/scripts/startdriver_unix.sh
ginkgo
/bin/bash example/scripts/stopdriver.sh

# TCP TESTS
export FIXTURE_FILENAME=example/fixtures/certification_tcp.json
/bin/bash example/scripts/startdriver_tcp.sh
ginkgo
/bin/bash example/scripts/stopdriver.sh

# JSON SPEC TESTS
export FIXTURE_FILENAME=example/fixtures/certification_json.json
/bin/bash example/scripts/startdriver_json.sh
ginkgo
/bin/bash example/scripts/stopdriver.sh

# JSON TLS SPEC TESTS
export FIXTURE_FILENAME=example/fixtures/certification_json.json
/bin/bash example/scripts/startdriver_json_tls.sh
ginkgo
/bin/bash example/scripts/stopdriver.sh

rm example/exec/fakedriver
