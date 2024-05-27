@echo off

SET ROOT_DIR=%CD%

if exist dist (
  rmdir /s /q dist
)

mkdir dist\unpack

cd .\src\jvm
if exist jvm.exe (
  del jvm.exe
)
go build jvm.go
move jvm.exe %ROOT_DIR%\dist\unpack
cd ..\..\

copy bin\elevate.cmd dist\unpack\elevate.cmd
copy bin\elevate.vbs dist\unpack\elevate.vbs

copy LICENSE dist\unpack\LICENSE

buildtools\7zr.exe a dist\env-manage.7z .\dist\unpack\*
