#! /bin/bash

root=$(cd "$(dirname "$0")"; pwd)

if [[ -d $root/dist ]]; then
    rm -rf $root/dist
fi

mkdir $root/dist
mkdir $root/dist/unpack
mkdir $root/dist/unpack/bin

echo Building vm

go build -o $root/dist/unpack/bin/vm $root/main.go
chmod +x $root/dist/unpack/bin/vm

echo copy install

cp $root/bin/install/install.sh $root/dist/unpack/install.sh
chmod +x $root/dist/unpack/install.sh

export CGO_ENABLED=0
export GOOS=windows
export GOARCH=amd64

echo Building vm.exe

go build -o $root/dist/unpack/bin/vm.exe $root/main.go

echo Building install.exe

go build -o $root/dist/unpack/install.exe $root/bin/install/install_windows.go

echo Building uninstall.exe

go build -o $root/dist/unpack/uninstall.exe $root/bin/uninstall/uninstall_windows.go

echo Building env-manage.tar.gz

tar -zcf $root/dist/env-manage.tar.gz -C $root/dist/unpack .

echo Building env-manage.zip

cd $root/dist/unpack

zip -qr $root/dist/env-manage.zip .

