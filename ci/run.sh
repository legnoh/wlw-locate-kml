#!/bin/sh

set -e -u -x

# prepare
go get -u github.com/golang/dep/cmd/dep
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

set +e
touch out/diff
diff result-`date -d '-1 month' +%Y%m01`.kml out/result-`date +%Y%m%d`.kml \
--ignore-matching-lines=".*description.*" \
--ignore-matching-lines=".*name=\"住所\".*" \
--ignore-matching-lines ".*name=\"ランキング\"" \
--ignore-matching-lines ".*name=\"ランキング結果(5〜1位)\"" \
--ignore-matching-lines=".*name=\"ライブラリ設置\"" \
--ignore-matching-lines=".*coordinates.*" \
--ignore-matching-lines ".*styleUrl.*" \
-U 0 >> out/diff
set -e

touch out/body
echo "# 新規出店・退店" >> out/body
echo '```diff' >> out/body
cat out/diff >> out/body
echo '```' >> out/body

cp -ar $GOPATH_REPO/vendor $INPUT_REPO/vendor
