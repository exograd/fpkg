all: build

build: FORCE
	GOBIN=$(CURDIR) CGO_ENABLED=0 go install ./...

test:
	go test -race -count 1 ./...

vet:
	go vet ./...

clean:
	$(RM) $(wildcard $(BIN_DIR)/*)

FORCE:

.PHONY: all build test vet
