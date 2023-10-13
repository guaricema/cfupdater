#!/bin/bash

RELEASE="0.0.2"
OS=("windows" "darwin" "linux" "linux")
ARCH=("amd64" "amd64" "amd64" "arm64")

rm -rf dist/

for key in "${!OS[@]}"
do
    echo Building compressed archive for "${OS[$key]}-${ARCH[$key]}"
    mkdir -p "dist/${OS[$key]}-${ARCH[$key]}"
    cp config.json.example README.md LICENSE "dist/${OS[$key]}-${ARCH[$key]}"
    env GOOS="${OS[$key]}" GOARCH="${ARCH[$key]}" go build -o dist/"${OS[$key]}-${ARCH[$key]}"/cfupdater
    cd "dist/${OS[$key]}-${ARCH[$key]}"
    if [ ${OS[$key]} == "windows" ]; then
        mv cfupdater cfupdater.exe
    fi
    zip "../cfupdater-"${OS[$key]}-${ARCH[$key]}"-$RELEASE.zip" * && cd ../..
done