call build.bat

%ROOT_DIR%\dist\unpack\bin\vm.exe jdk list

mkdir %ROOT_DIR%\dist\unpack\download\jdk\

robocopy D:\temp\jdk %ROOT_DIR%\dist\unpack\download\jdk\ /E

%ROOT_DIR%\dist\unpack\bin\vm.exe jdk add jdk11 D:\soft\java\jdk11

%ROOT_DIR%\dist\unpack\bin\vm.exe jdk list

%ROOT_DIR%\dist\unpack\bin\vm.exe node list

