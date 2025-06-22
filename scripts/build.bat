@echo off

REM Set the input name for the server and strip quotes
set SERVER_NAME=%1
set SERVER_NAME=%SERVER_NAME:"=%
if "%SERVER_NAME%"=="" set SERVER_NAME=login-server

REM Set the Go environment variables for building for Linux (64-bit)
set GOARCH=amd64
set GOOS=linux

REM Build for Linux
echo Building %SERVER_NAME% for Linux...
go build -ldflags="-w -s" -o bin\%SERVER_NAME%\%SERVER_NAME% cmd\%SERVER_NAME%\main.go

REM Reset Go environment variables to their defaults
set GOARCH=
set GOOS=

REM Build for Windows
echo Building %SERVER_NAME% for Windows...
go build -ldflags="-w -s" -o bin\%SERVER_NAME%\%SERVER_NAME%.exe cmd\%SERVER_NAME%\main.go

echo %SERVER_NAME% build complete!
