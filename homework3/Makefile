export GO111MODULE=auto

VERSION=v2.0
BINARY=httpserver
BUILD=`date +%FT%T%z`

PACKAGES=`go list ./... | grep -v /vendor/`
VETPACKAGES=`go list ./... | grep -v /vendor/ | grep -v /examples/`
GOFILES=`find . -name "*.go" -type f -not -path "./vendor/*"`

.PHONY: build list fmt fmt-check vet install release push clean

build:
	@echo "building ${BINARY}"
	mkdir -p bin/amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/amd64/${BINARY} .
	
list:
	@echo ${PACKAGES}
	@echo "*******************************************"
	@echo ${VETPACKAGES}
	@echo "*******************************************"
	@echo ${GOFILES}
	
fmt:
	@gofmt -s -w ${GOFILES}
	
fmt-check:
	@diff=$$(gofmt -s -d ${GOFILES}); \
	if [ -n "$$diff" ]; then \
		echo "Please run 'make fmt' and commit the result:"; \
		echo "$${diff}"; \
		exit 1; \
    fi;

vet:
	@go vet ${VETPACKAGES}
	
install:
	@govendor sync -v
	
release: build
	@echo "building ${BINARY} container"
	@docker build -t yhjhero/${BINARY}:${VERSION} .
	
push: release
	@echo "pushing yhjhero/${BINARY}"
	@docker push yhjhero/${BINARY}:${VERSION}
	
clean:
	@rm -rf bin
