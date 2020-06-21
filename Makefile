# Go settings
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOLINT=golint
GOGET=$(GOCMD) get

# Build settings
BINARY_PATH=./bin/
BINARY_NAME=huekit

# Test Settings
TEST_FILES := $(shell $(GOCMD) list ./...)

all: deps tests build


build:
		$(GOBUILD) -o $(BINARY_PATH)$(BINARY_NAME) -v cmd/huekit/main.go

run:
		$(GOBUILD) -o $(BINARY_PATH)$(BINARY_NAME) -v cmd/huekit/main.go
		$(BINARY_PATH)$(BINARY_NAME) --config ./configs/config.yml

tests:
		mkdir -p report
		$(GOTEST) -v -short -covermode=count -coverprofile report/cover.out $(TEST_FILES)
		$(GOCMD) tool cover -html=report/cover.out -o report/cover.html
		$(GOLINT) -set_exit_status $(TEST_FILES)
		staticcheck $(TEST_FILES)

clean:
		$(GOCLEAN)
		rm -rf $(BINARY_PATH)
		rm -rf ./report/

deps:
		GO111MODULE=on $(GOCMD) mod vendor

release: clean
		mkdir -p $(BINARY_PATH)
		cp ./configs/config.yml.dist $(BINARY_PATH)config.yml
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_PATH)$(BINARY_NAME) -ldflags="-s -w" -a -installsuffix cgo -v cmd/huekit/main.go
		cd $(BINARY_PATH) && tar cvzf huekit_linux_amd64.tar.gz $(BINARY_NAME) config.yml
		rm -rf $(BINARY_PATH)$(BINARY_NAME)
		CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_PATH)$(BINARY_NAME) -a -installsuffix cgo -v cmd/huekit/main.go
		cd $(BINARY_PATH) && zip huekit_windows_amd64.zip $(BINARY_NAME) config.yml
		rm -rf $(BINARY_PATH)$(BINARY_NAME)
		CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_PATH)$(BINARY_NAME) -a -installsuffix cgo -v cmd/huekit/main.go
		cd $(BINARY_PATH) && tar cvzf huekit_macos_amd64.tar.gz $(BINARY_NAME) config.yml
		rm -rf $(BINARY_PATH)$(BINARY_NAME)
		rm -rf $(BINARY_PATH)config.yml
