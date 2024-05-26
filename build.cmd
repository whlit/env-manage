@echo off

if exist src\jvm.exe (
  del src\jvm.exe
)

if not exist dist (
  mkdir dist
)

if not exist dist\unpack (
  mkdir dist\unpack
) else (
  rmdir /s /q dist\unpack
  mkdir dist\unpack
)

cd .\src
go build jvm.go
cd ..\

move src\jvm.exe dist\unpack

copy bin\elevate.cmd dist\unpack\elevate.cmd
copy bin\elevate.vbs dist\unpack\elevate.vbs
