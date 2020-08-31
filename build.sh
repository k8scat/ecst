#!/bin/bash
set -ev

function build() {
  os=$1
  out=build/"${os}"/vss
  rm -f build/"${os}"/vss
  CGO_ENABLED=0 GOOS="${os}" GOARCH=amd64 go build -o "${out}" main.go
  chmod a+x "${out}"
}

platform=$1
if [ -z "${platform}" ];then
  build darwin
  build linux
  build windows
elif [ "${platform}" != "darwin" ] && [ "${platform}" != "linux" ] && [ "${platform}" != "windows" ];then
  echo "usage: ./build.sh darwin|linux|windows"
  exit 1
else
  build "${platform}"
fi


