@echo off

SET ORIG=%CD%
SET GOBIN=%CD%\bin

if exist src\jvm.exe (
  del src\jvm.exe
)

cd .\src
go build jvm.go

move jvm.exe "%GOBIN%"
cd ..\

