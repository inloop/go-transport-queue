install:
	go get ./...
	go build -o go-transport-queue

deploy-local:
	make install
	mv go-transport-queue /usr/local/bin/
	chmod +x /usr/local/bin/go-transport-queue
