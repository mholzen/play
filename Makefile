get:
	(cd ..; scp -r marchome@192.168.0.60:develop/mholzen/play-go .)

run:
	go run main.go server.go

on:
	curl -vvv http://192.168.0.80:1323/controls/dimmer/255

off:
	curl -vvv http://192.168.0.80:1323/controls/dimmer/0
