build:
	go get ./...
	go build -o go-transport-queue

install:
	make build
	mv go-transport-queue /usr/local/bin/
	chmod +x /usr/local/bin/go-transport-queue
