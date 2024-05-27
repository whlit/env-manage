@echo off

if exist src\jvm.exe (
  del src\jvm.exe
)

if exist dist (
  rmdir /s /q dist
)

mkdir dist\unpack

cd .\src
go build jvm.go
cd ..\

move src\jvm.exe dist\unpack

copy bin\elevate.cmd dist\unpack\elevate.cmd
copy bin\elevate.vbs dist\unpack\elevate.vbs

copy LICENSE dist\unpack\LICENSE

buildtools\7zr.exe a dist\env-manage.7z .\dist\unpack\*
