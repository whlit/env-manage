# !/bin/bash

root=$(cd "$(dirname "$0")/../"; pwd)
download=$root/dist/unpack/download
cache=~/tmp/env-manage/downloads

if [[ -d $download ]]; then
    if [[ ! "$(ls -A $download)" ]]; then
        exit 0
    fi
fi

if [[ ! -d ~/tmp ]]; then
    mkdir ~/tmp
fi

if [[ ! -d ~/tmp/env-manage ]]; then
    mkdir ~/tmp/env-manage
fi

if [[ ! -d ~/tmp/env-manage/downloads ]]; then
    mkdir ~/tmp/env-manage/downloads
fi

mv $root/dist/unpack/download/* ~/tmp/env-manage/downloads
