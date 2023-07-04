BINSUFFIX := $(shell if [ "${GOOS}" -a "${GOARCH}" ]; then echo "-${GOOS}-${GOARCH}"; else echo ""; fi)

all: build

## init: Install required apps
init:
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest

## build: Build core (default)
build:
	CGO_ENABLED=${CGO_ENABLED} GOOS=${GOOS} GOARCH=${GOARCH} go build -o corecli${BINSUFFIX} -trimpath

## test: Run all tests
test:
	go test -race -coverprofile=/dev/null -v ./...

## vet: Analyze code for potential errors
vet:
	go vet ./...

## fmt: Format code
fmt:
	go fmt ./...

## vulncheck: Check for known vulnerabilities in dependencies
vulncheck:
	govulncheck ./...

## update: Update dependencies
update:
	go get -u
	@-$(MAKE) tidy
	@-$(MAKE) vendor

## tidy: Tidy up go.mod
tidy:
	go mod tidy

## vendor: Update vendored packages
vendor:
	go mod vendor

## run: Build and run core
run: build
	./core

## lint: Static analysis with staticcheck
lint:
	staticcheck ./...

## commit: Prepare code for commit (vet, fmt, test)
commit: vet fmt lint test build
	@echo "No errors found. Ready for a commit."

## release: Build a release binary of core
release:
	CGO_ENABLED=${CGO_ENABLED} GOOS=${GOOS} GOARCH=${GOARCH} go build -o corecli${BINSUFFIX} -trimpath -ldflags="-s -w"

.PHONY: help init build test vet fmt vulncheck vendor commit coverage lint release update

## help: Show all commands
help: Makefile
	@echo
	@echo " Choose a command:"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
