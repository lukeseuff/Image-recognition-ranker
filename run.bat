@echo off
set dir="%~dp0"
echo %GOPATH%|find %dir% >nul
if errorlevel 1 (set GOPATH=%GOPATH%;%dir%) else (echo found)
go run main.go

