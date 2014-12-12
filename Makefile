GOPATH := $(shell pwd)
BIN := siggi

.PHONY: all
all: $(BIN)

$(BIN): $(shell find src -name '*.go')
	go get -d -v $(BIN)/...
	go build $(BIN)

.PHONY: fmt
fmt:
	go fmt sig{gi,hub}

.PHONY: clean
clean:
	rm -f $(BIN)
