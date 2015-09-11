#!/bin/bash
REV=$(git rev-parse HEAD | cut -c1-6)
if [ -z "$CIRCLE_ARTIFACTS" ]; then
  DEST="./pkg"
else
  DEST="${CIRCLE_ARTIFACTS}"
fi
goxc -wd=. -build-ldflags="-X main.Version ${REV}" -d=${DEST}

if [ ! -z "$CIRCLE_ARTIFACTS" ]; then
	cp -r $CIRCLE_ARTIFACTS ./pkg
fi