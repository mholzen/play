host=ubuntu-1
#home=marchome
# cannot for the life of me figure out how to change dns on ubuntu
home=192.168.50.89

cwd := $(shell pwd)

OS := $(shell uname -s |  tr '[:upper:]' '[:lower:]')

build:
	GOOS=$(OS) go build -o bin/server ./cmd/server

run:
	CompileDaemon --build="go build -o bin/server ./cmd/server" --command="./bin/server"

open:
	open http://localhost:1300/api/v2/root

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
	env ROOT=$(cwd) go test -timeout 5s ./... | grcat ~/.grc/conf.go-test

test-watch:
	CompileDaemon --build="make test || false" --command="echo 'Tests completed'" --pattern="(.+\.go|.+\.yaml)"

build-docker:
	docker buildx build --platform linux/amd64 -t play-go --load .

run-docker:
	docker run -p 8080:8080 --device=/dev/tty.usbserial-ENVVVCOF:/dev/tty.usbserial-ENVVVCOF play-go

run-docker-remote:
	ssh marc@$(host) "docker pull ubuntu-1:5000/play-go &&  docker run -p 1300:1300 --add-host=host.docker.internal:host-gateway --device=/dev/ttyUSB0:/dev/ttyUSB0 ubuntu-1:5000/play-go"

push-docker:
	docker tag play-go $(host):5000/play-go
	docker push $(host):5000/play-go
