# !/bin/bash

root=$(cd "$(dirname "$0")/../"; pwd)
cache=~/tmp/env-manage/downloads
download=$root/dist/unpack/download

if [[ ! -d $cache ]]; then
    exit 0
fi

if [[ ! -d $root/dist ]]; then
    mkdir $root/dist
fi

if [[ ! -d $root/dist/unpack ]]; then
    mkdir $root/dist/unpack
fi

if [[ ! -d $root/dist/unpack/download ]]; then
    mkdir $root/dist/unpack/download
fi

cp -r $cache/* $download
