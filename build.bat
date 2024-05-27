@echo off

SET ROOT_DIR=%CD%

if exist dist (
  rmdir /s /q dist
)

mkdir dist\unpack\bin

cd .\src

echo ----------------------------
echo Building jvm.exe
echo ----------------------------

if exist .\jvm.exe (
  del .\jvm.exe
)
go build .\jvm\jvm.go
move .\jvm.exe %ROOT_DIR%\dist\unpack\bin


echo ----------------------------
echo Building install.exe
echo ----------------------------

if exist .\install.exe (
    del .\install.exe
)
go build .\tools\install\install.go
move .\install.exe %ROOT_DIR%\dist\unpack


echo ----------------------------
echo Building uninstall.exe
echo ----------------------------

if exist .\uninstall.exe (
    del .\uninstall.exe
)
go build .\tools\uninstall\uninstall.go
move .\uninstall.exe %ROOT_DIR%\dist\unpack

cd ..\


echo ----------------------------
echo Add other require files
echo ----------------------------

mkdir dist\unpack\lib

copy bin\elevate.cmd dist\unpack\lib\elevate.cmd
copy bin\elevate.vbs dist\unpack\lib\elevate.vbs
copy LICENSE dist\unpack\LICENSE

echo ----------------------------
echo Package
echo ----------------------------

buildtools\7zr.exe a dist\env-manage.7z .\dist\unpack\*
