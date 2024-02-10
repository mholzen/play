host=ubuntu-1
home=marchome

run:
	go run main.go server.go

on:
	curl -vvv http://$(host):1323/controls/dimmer/255

off:
	curl -vvv http://$(host):1323/controls/dimmer/0

ssh:
	ssh -A $(host) -l marc

pull:
	(cd ..; scp -r marchome@$(home):develop/mholzen/play-go .)

push:
	(cd ..; scp -r play-go/ marc@$(host):)
