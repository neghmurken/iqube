.PHONY: build run

build:
	go build -o iqube .

run: build
	./iqube
