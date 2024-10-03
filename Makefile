GO := go

GOFLAGS :=
LDFLAGS := -ldflags="-s -w"
GCFLAGS := -gcflags="-m"

SOURCES := ./cmd/cli/run.go
TARGET := cgnlog

build: clean
	$(GO) build $(GOFLAGS) $(LDFLAGS) $(GCFLAGS) -o $(TARGET) $(SOURCES)

clean:
	-@rm -f $(TARGET)

install: build
	-@mkdir -p $(HOME)/bin
	cp $(TARGET) $(HOME)/bin/

uninstall:
	-@rm -f $(HOME)/bin/$(TARGET)
	