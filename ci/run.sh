#!/bin/bash

set -e -u -x

apk add --no-cache git && go get -u github.com/golang/dep/cmd/dep

# prepare
JOB_DIR=$PWD
INPUT_REPO=$JOB_DIR/repo
GOPATH_REPO=$GOPATH/src/github.com/legnoh/wlw-locate-kml
mkdir -p $GOPATH/src/github.com/legnoh
cp -ar $INPUT_REPO $GOPATH_REPO
cd $GOPATH_REPO
dep ensure

# run
go run main.go

# output
cd $JOB_DIR
cp $GOPATH_REPO/result-*.kml out/


# make release info
date +%Y/%m/%d > out/name
date +%Y%m%d > out/tag

wget https://github.com/legnoh/wlw-locate-kml/releases/download/`date -d '-1 month' +%Y%m01`/result-`date -d '-1 month' +%Y%m01`.kml

diff result-`date -d '-1 month' +%Y%m01`.kml result-`date +%Y%m%d`.kml \
--ignore-matching-lines=".*name.*" \
--ignore-matching-lines=".*description.*" \
--ignore-matching-lines=".*coordinates.*" \
--ignore-matching-lines ".*ランキング.*" \
--ignore-matching-lines ".*styleUrl.*" \
-U 0 > out/diff

echo "# 新規出店・退店\n\n\`\`\`diff\n" > out/body
cat out/diff >> out/body
echo "\`\`\`" >> out/body
