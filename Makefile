PWD?=$(shell pwd)
BUILDDIR := $(PWD)/build

run_simple:
	@go run ./examples/simple/main.go -file ./examples/simple/input.csv -s 3 -log ./examples/simple/log.txt

run_complex:
	@go run ./examples/complex/main.go -file ./examples/complex/input.csv -s 6

build: clean lint test
	@echo "+ $@"
	@echo "Building..."
	go build -o $(BUILDDIR)/example ./examples/simple/main.go
	cp ./examples/simple/input.csv $(BUILDDIR)/input.csv

clean:
	@echo "+ $@"
	@rm -rf $(BUILDDIR)
	@mkdir -p $(BUILDDIR)
test:
	@echo "+ $@"
	go test -v ./pkg/...

lint:
	@echo "+ $@"
	golangci-lint run ./pkg/...
	gofmt -l -s -w .

.PHONY: test build lint clean