#!/bin/bash

set -ex

hack/dind hack/test/unit

# prevent legacy test suites from running. they take a very long time
sed -i -e '/run_test_integration_legacy_suites$/ s/^/#/' hack/make/.integration-test-helpers

hack/dind hack/make.sh dynbinary test-integration
