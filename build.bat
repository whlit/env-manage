@echo off

set root=%~dp0

if exist %root%\dist (
  rmdir /s /q dist
)

mkdir %root%\dist\unpack\bin

echo Building vm.exe

go build -o %root%\dist\unpack\bin\vm.exe %root%\main.go

echo Building install.exe

go build -o %root%\dist\unpack\install.exe %root%\bin\install\install_windows.go

echo Building uninstall.exe

go build -o %root%\dist\unpack\uninstall.exe %root%\bin\uninstall\uninstall_windows.go

