@echo off

rem # Check if powershell is in path
where /q pwsh
IF ERRORLEVEL 1 (
    where /q powershell
    IF ERRORLEVEL 1 (
        echo Neither pwsh.exe nor powershell.exe was found in your path.
        echo Please install powershell it is required
        exit /B
    ) ELSE (
        set ps=powershell
    )
) ELSE (
    set ps=pwsh
)

rem ps is the installed powershell
%ps% -executionpolicy remotesigned -File "%~dp0/java/jvm.ps1" %* 
