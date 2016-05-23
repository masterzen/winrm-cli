#!/bin/bash
set -e

export GO15VENDOREXPERIMENT=1

TESTS=$(go list ./... | grep -v vendor)
go test -v ${TESTS}

mkdir -p .cover
go list ./... | grep -v vendor/ | xargs -I% bash -c 'name="%"; go test -covermode=count % --coverprofile=.cover/${name//\//_} '
echo "mode: count" > profile.cov
cat .cover/* | grep -v mode >> profile.cov
rm -rf .cover

go tool cover -func=profile.cov
