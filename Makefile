.PHONY: all
all: build

.PHONY: build
build:
	go build -o notionctl

.PHONY: test
test:
	go test ./...

.PHONY: clean
clean:
	rm -f notionctl