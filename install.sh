#!/bin/bash

go build -o tgserver fileserver/main.go
go build -o tgclient fileclient/main.go

cp ./tgclient /usr/local/bin/tgclient