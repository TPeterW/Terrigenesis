all:	clean build

build:
	go build -o tgserver fileserver/main.go
	go build -o tgclient fileclient/main.go

install:
	cp tgclient /usr/local/bin/tgclient

clean:
	rm -rf tgclient tgserver