@echo off
set dir="%~dp0"
set GOPATH=%GOPATH%;%dir%
go run main.go

