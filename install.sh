#!/bin/bash

go build -o tgserver fileserver/main.go
go build -o tgclient fileclient/main.go

# cp ./tgserver /usr/local/bin/tgserver
# cp ./tgclient /usr/local/bin/tgclient