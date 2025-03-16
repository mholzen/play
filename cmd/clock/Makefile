run:
	go run clock.go

build:
	go build -o clock clock.go
	chmod +x clock

symlink:
	ln -sf $(shell pwd)/clock $(DEVBIN)/clock

install:

get:
	(cd ../..; make get)

get-run: get run