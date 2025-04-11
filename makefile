all: gotool build

build:
	 CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ./bin/${BINARY}

gotool:
	go fmt ./
	go vet ./