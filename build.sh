#! /bin/bash

root=$PWD

if [[ -d $root/dist ]]; then
    rm -rf $root/dist
fi

mkdir $root/dist
mkdir $root/dist/unpack
mkdir $root/dist/unpack/bin

cd ./src

echo Building vm

go build -o $root/dist/unpack/bin/vm $root/src/bin/linux/main.go
chmod +x $root/dist/unpack/bin/vm

echo copy install

cp $root/src/bin/linux/install/install.sh $root/dist/unpack/install.sh
chmod +x $root/dist/unpack/install.sh

export CGO_ENABLED=0
export GOOS=windows
export GOARCH=amd64

echo Building vm.exe

go build -o $root/dist/unpack/bin/vm.exe $root/src/bin/windows/main.go

echo Building install.exe

go build -o $root/dist/unpack/install.exe $root/src/bin/windows/install/install.go

echo Building uninstall.exe

go build -o $root/dist/unpack/uninstall.exe $root/src/bin/windows/uninstall/uninstall.go

echo Building env-manage.tar.gz

tar -zcf $root/dist/env-manage.tar.gz -C $root/dist/unpack .

echo Building env-manage.zip

cd $root/dist/unpack

zip -qr $root/dist/env-manage.zip .

