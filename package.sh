#!/bin/bash -xe

# This usually looks like "v4"
BUNDLE_VERSION=${BUNDLE_VERSION:=latest}

SCRATCH_DIR="$(mktemp -d)"
mkdir -p ${SCRATCH_DIR}/sample-app
cp -a * ${SCRATCH_DIR}/sample-app
tar -zcvC ${SCRATCH_DIR} -f sample-app-${BUNDLE_VERSION}.tgz sample-app
gsutil cp -a public-read sample-app-${BUNDLE_VERSION}.tgz gs://gke-spinnaker
