#!/bin/bash

set -x

# prepare
RELEASES_LATEST_URL=https://api.github.com/repos/legnoh/wlw-locate-kml/releases/latest
RELEASE_LATEST_ASSET_URL=$(curl -s ${RELEASES_LATEST_URL} | jq -r '.assets[].browser_download_url')
RELEASE_LATEST_ASSET_FILE=$(echo ${RELEASE_LATEST_ASSET_URL} | awk -F "/" '{ print $NF }')
RELEASE_NEW_NAME=$(date +%Y/%m/%d)
RELEASE_NEW_TAG=$(date +%Y%m%d)
RELEASE_NEW_ASSET_FILE=result-${RELEASE_NEW_TAG}.kml
RELEASE_NEW_DRAFT=draft-${RELEASE_NEW_TAG}.md
go get -u github.com/tcnksm/ghr

# get previous release
curl -sL ${RELEASE_LATEST_ASSET_URL} -o ${RELEASE_LATEST_ASSET_FILE}

# run
go run main.go

# make release drafts
cat << EOS > ${RELEASE_NEW_DRAFT}
# 新規出店・退店

\`\`\`diff
```
diff ${RELEASE_LATEST_ASSET_FILE} ${RELEASE_NEW_ASSET_FILE} \
--ignore-matching-lines=".*description.*" \
--ignore-matching-lines=".*name=\"住所\".*" \
--ignore-matching-lines=".*name=\"ランキング\"" \
--ignore-matching-lines=".*name=\"ランキング結果(5〜1位)\"" \
--ignore-matching-lines=".*name=\"ライブラリ設置\"" \
--ignore-matching-lines=".*coordinates.*" \
--ignore-matching-lines=".*styleUrl.*" \
-U 0
```
\`\`\`
EOS

# post new github prerelease
$GOPATH/bin/ghr \
-n ${RELEASE_NEW_NAME} \
-b "$(cat ${RELEASE_NEW_DRAFT})" \
-prerelease ${RELEASE_NEW_TAG} \
${RELEASE_NEW_ASSET_FILE}
