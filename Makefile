build:
	go build -o scalc scalc.go
	chmod +x scalc

test:
	go test .

.DEFAULT_GOAL := build
