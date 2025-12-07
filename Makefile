format:
	go fmt .

vet:
	go vet .

build: vet
	go build .

run: build
	./keystats
