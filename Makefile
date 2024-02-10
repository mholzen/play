host=ubuntu-1
#home=marchome
# cannot for the life of me figure out how to change dns on ubuntu
home=192.168.50.89

build:
	go build -o main main.go server.go

run:
	go run main.go server.go

on:
	curl -vvv http://$(host):1323/controls/dimmer/255

off:
	curl -vvv http://$(host):1323/controls/dimmer/0

ssh:
	ssh -A $(host) -l marc

get:
	(cd ..; scp -r marchome@$(home):develop/mholzen/play-go .)

pull: get

push:
	(cd ..; scp -r play-go/ marc@$(host):)
