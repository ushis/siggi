GOPATH := $(shell pwd)
BIN := siggi

.PHONY: all
all: $(BIN)

$(BIN): $(shell find src -name '*.go') deps
	go build $(BIN)

.PHONY: deps
deps:
	go get -d -v $(BIN)/...

.PHONY: fmt
fmt:
	go fmt sig{gi,hub}

.PHONY: clean
clean:
	rm -f $(BIN)
