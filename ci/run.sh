#!/bin/bash

set -e -u -x

# prepare
mkdir -p $GOPATH/src/github.com/legnoh
cp -r repo $GOPATH/src/github.com/legnoh/wlw-locate-kml
cd $GOPATH/src/github.com/legnoh/wlw-locate-kml
dep ensure

# run
go run main.go

# output
cd -
cp $GOPATH/src/github.com/legnoh/wlw-locate-kml/result-*.kml out/


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

echo "@legnoh wlw-locate-kml is updated.\nhttps://github.com/legnoh/wlw-locate-kml/releases" > out/slack
