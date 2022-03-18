SHELL = /bin/sh

GO = go
TESTFLAGS =

GOIMPORTS = goimports


.PHONY: all
all: fmt test

.PHONY: fmt
fmt:
	$(GOIMPORTS) -d -l -local github.com/now/x/ -w .

.PHONY: test
test:
	$(GO) test $(TESTFLAGS) ./...
