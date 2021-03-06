
GLIDE = $(GOPATH)/bin/glide

all: build

$(GLIDE):
	# On travis the bin directory is not present by default
	mkdir -p $(GOPATH)/bin
	curl https://glide.sh/get | sh

install: $(GLIDE)
	$(GOPATH)/bin/glide install

build: install
	go clean
	go build
	go test -v . ./file

linux: install
	env GOOS=linux GOARCH=amd64 go build
