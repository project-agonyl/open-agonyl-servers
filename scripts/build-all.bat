@echo off
setlocal enabledelayedexpansion

REM Get the directory where this script is located
set SCRIPT_DIR=%~dp0
set PROJECT_ROOT=%SCRIPT_DIR%..

REM Change to the project root directory
cd /d "%PROJECT_ROOT%"

echo Building all servers...

REM Iterate through all directories in the cmd folder
for /d %%i in (cmd\*) do (
    if exist "%%i" (
        REM Extract the server name from the directory path
        for %%j in ("%%i") do set server_name=%%~nj
        
        echo Building !server_name!...
        
        REM Run the build script for this server
        call scripts\build.bat "!server_name!"
        
        if !errorlevel! equ 0 (
            echo !server_name! build completed successfully!
        ) else (
            echo ERROR: !server_name! build failed!
            exit /b 1
        )
        echo.
    )
)

echo All servers built successfully!
