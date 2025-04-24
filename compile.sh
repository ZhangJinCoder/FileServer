#!/bin/bash
rm -rf go.mod
go mod init fileserver
go mod tidy
go mod vendor
go build -o fileserver cmd/fileserver/main.go


