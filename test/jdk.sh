#! /bin/bash

root=$(cd "$(dirname "$0")/../"; pwd)

if [[ ! -d $root/dist/unpack ]]; then
    $root/build.sh
fi

$root/dist/unpack/bin/vm jdk init

mkdir $root/dist/unpack/tmp

for i in `seq 1 10`; do
    mkdir $root/dist/unpack/tmp/jdk-$i
    $root/dist/unpack/bin/vm jdk add jdk-$i $root/dist/unpack/tmp/jdk-$i
    touch $root/dist/unpack/tmp/jdk-$i/jdk-$i.txt
done

$root/dist/unpack/bin/vm jdk list
