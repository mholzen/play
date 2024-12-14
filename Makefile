host=ubuntu-1
#home=marchome
# cannot for the life of me figure out how to change dns on ubuntu
home=192.168.50.89

cwd := $(shell pwd)

OS := $(shell uname -s |  tr '[:upper:]' '[:lower:]')

build:
	GOOS=$(OS) go build -o main main.go server.go handlers.go

run:
	# go run main.go server.go
	CompileDaemon --build="go build -o main main.go server.go handlers.go" --command="./main"

run-dev:
	~/go/bin/air

open:
	open http://$(host):1300

on:
	curl -vvv http://$(host):1300/controls/dimmer/255

off:
	curl -vvv http://$(host):1300/controls/dimmer/0

ssh:
	ssh -A $(host) -l marc -t "cd play; source .setup; zsh --login --interactive"

get:
	rsync -avz --delete --exclude .git --exclude node_modules --exclude main -e ssh marchome@$(home):develop/mholzen/play-go/ .

pull: get

push:
	rsync -avz --delete --exclude .git --exclude node_modules --exclude main -e ssh . marc@$(host):play-go

status:
	sudo systemctl status play-go.service

stop:
	sudo systemctl stop play-go-watcher.service
	sudo systemctl stop play-go.service

start:
	sudo systemctl start play-go.service
	sudo systemctl start play-go-watcher.service

restart: stop start

log:
	journalctl -u play-go.service -f

live:
	(cd cmd/live; go run live.go)

test:
	env ROOT=$(cwd) go test ./... | grcat ~/.grc/go.conf
