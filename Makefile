

all: linetype_string.go bin
	go build -o bin/ ./...

test: all
	go test

linetype_string.go: log.go
	go generate

bin:
	mkdir bin

clean:
	rm -rf bin

realclean: clean
	rm -f linetype_string.go
