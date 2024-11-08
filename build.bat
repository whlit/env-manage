@echo off

SET ROOT_DIR=%CD%

if exist dist (
  rmdir /s /q dist
)

mkdir dist\unpack\bin

cd .\src


echo ----------------------------
echo Building vm.exe
echo ----------------------------

if exist .\vm.exe (
  del .\vm.exe
)
go build -o vm.exe .\bin\windows\main.go
move .\vm.exe %ROOT_DIR%\dist\unpack\bin

echo ----------------------------
echo Building install.exe
echo ----------------------------

if exist .\install.exe (
    del .\install.exe
)
go build -o install.exe .\bin\windows\install\install.go


echo ----------------------------
echo Building uninstall.exe
echo ----------------------------

if exist .\uninstall.exe (
    del .\uninstall.exe
)
go build -o uninstall.exe .\bin\windows\uninstall\uninstall.go

move .\*.exe %ROOT_DIR%\dist\unpack

cd ..\

echo ----------------------------
echo Package
echo ----------------------------

build\7zr.exe a dist\env-manage.7z .\dist\unpack\*
