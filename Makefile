
all: build 

.PHONY: test 
test:
	go test -p 1 ./...

.PHONY: testcov
testcov:
	go test -race -covermode=atomic -coverprofile=coverage.out -p 1 ./...
	go tool cover -html=coverage.out -o coverage.html

.PHONY: build
build:
	go build 

.PHONY: clean
clean:
	rm -rfv piot-server 
	rm coverage.html
	rm coverage.out
