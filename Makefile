PWD=$(shell pwd)
TARGET=github_statistics

build:
	go build -o ${TARGET} cmd/main.go

run:
	go run cmd/main.go

clean:
	rm -f ${TARGET}

rebuild: clean build
