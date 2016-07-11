
build:
	go clean
	go build
	go test -v . ./file

linux:
	env GOOS=linux GOARCH=amd64 go build
