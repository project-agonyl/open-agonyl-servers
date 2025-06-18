@echo off

REM Set the input name for the server
set SERVER_NAME=%1
if "%SERVER_NAME%"=="" set SERVER_NAME=login-server

go run .\cmd\%SERVER_NAME%\main.go
