PROJECT    		= dummy-bc-platform
DOCKER_IMAGE 	= wesovilabs/$(PROJECT)
DATE      		?= $(shell date +%FT%T%z)
GOPATH     		= $(CURDIR)
GOLANG_VERSION	= 1.10

# Tools
GO         		= GOPATH=$(GOPATH) go
GOBUILD 		= $(GO) build -a -installsuffix nocgo $(LDFLAGS)
GOTEST	 		= $(GO) test -p=1
GORUN	 		= $(GO) run
GOLINT     		= $(GO) run $(DIR_TOOLS)/vendor/github.com/golang/lint/golint/import.go $(DIR_TOOLS)/vendor/github.com/golang/lint/golint/golint.go
GOVET     		= $(GO) tool vet
GO2XUNIT   		= $(GO) run $(DIR_TOOLS)/vendor/github.com/tebeka/go2xunit/main.go  $(DIR_TOOLS)/vendor/github.com/tebeka/go2xunit/cmdline.go

# Directories
PACKAGES 		= $(shell $(GO) list -f '{{.Dir}}' ./src/$(PROJECT)/... | grep -v /vendor/)
PACKAGES_TEST 	= $(shell $(GO) list -f '{{ if or .TestGoFiles .XTestGoFiles }}{{.ImportPath}}{{ end }}' ./src/$(PROJECT)/... | grep -v /vendor/)
DIR_TESTREPORT  = test
DIR_BUILD 		= build
DIR_RESOURCES 	= resources
DIR_TOOLS		= $(CURDIR)/src/tools


# Docker
DOCKER_COMPOSE  	= docker-compose.yml
DOCKER_DEPS 		= mongodb initializer
DOCKER_COMPOSE_UP 	= docker-compose -f $(DOCKER_COMPOSE) up -d
DOCKER_COMPOSE_DOWN = docker-compose -f $(DOCKER_COMPOSE) down -v --remove-orphans
DOCKER_COMPOSE_LOGS = docker-compose -f $(DOCKER_COMPOSE) logs

STACK_RUN           = docker stack deploy --compose-file $(DOCKER_COMPOSE)

# Misc
COMMIT 			= $(shell git log -1 --format="%h" 2>/dev/null || echo "0")
VERSION			= $(shell git describe --tags --always)
BUILD_DATE 		= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS	 		= -ldflags "\
  -X $(PROJECT)/constants.COMMIT=$(COMMIT) \
  -X $(PROJECT)/constants.RELEASE_VERSION=$(VERSION) \
  -X $(PROJECT)/constants.BUILD_DATE=$(BUILD_DATE) \
  "

.PHONY: all
all: fmt lint test build docker-build docker-publish ; @ ## [ fmt lint test build docker-build docker-publish ]

.PHONY: development
development: lint test build ; @ ## [ lint test build ]

.PHONY: release
release: docker-build docker-publish ; @ ## [ docker-build docker-publish ]

.PHONY: clean
clean: ; @ ## Delete temporary directories
	rm -rf bin
	rm -rf pkg
	rm -rf $(DIR_BUILD)
	rm -rf $(DIR_DIST)
	rm -rf test
	rm -rf src/$(PROJECT)/vendor/
	rm -rf $(DIR_TOOLS)/vendor



.PHONY: deps
deps: ; @ ## Download required Golang libraries for project & Tools
	cd $(CURDIR)/src/$(PROJECT) ; GOPATH=$(GOPATH)  glide update; GOPATH=$(GOPATH)  glide install --force;
	cd $(DIR_TOOLS)		; GOPATH=$(GOPATH)  glide update; GOPATH=$(GOPATH)  glide install --force;


.PHONY: fmt
fmt: ; @ ## Code formatter
	for pkg in $(PACKAGES); do \
		gofmt -l -w  -e $$pkg/*.go; \
	done

.PHONY: lint
lint: ; @ ## Code analysis
	for pkg in $(PACKAGES); do \
    	$(GOVET) $$pkg/*.go; \
    done;\
    $(GOLINT) -set_exit_status $(PACKAGES);

.PHONY: test
junit ?=0
testArgs := -v -short -cover
ifeq ($(junit),1)
testArgs += | tee $(DIR_TESTREPORT)/$(PROJECT)-test.output
endif
test: ; @ ## Run unit tests [junit:0|1 ]
ifeq ($(junit),1)
	mkdir -p test
endif
	$(GOTEST) $(PACKAGES_TEST)  ${testArgs}; \
	status=$$?; \
	exit $$status
ifeq ($(junit),1)
	$(GO2XUNIT) -fail -input $(DIR_TESTREPORT)/$(PROJECT)-test.output -output $(DIR_TESTREPORT)/$(PROJECT)-test.xml
endif


.PHONY: build
build: ; @ ## Build mollydb executables for linux and darwin
	rm -rf $(CURDIR)/$(DIR_BUILD)/*; \
	GOARCH=amd64 CGO_ENABLED=0  GOOS=linux $(GOBUILD) -o $(DIR_BUILD)/$(PROJECT).linux $(PROJECT);
	chmod +x $(DIR_BUILD)/$(PROJECT).linux

.PHONY: build-all
build-all: ; @ ## Build mollydb executables for all architectures
	rm -rf $(CURDIR)/$(DIR_BUILD)/*; \
	GOARCH=amd64 GOOS=linux $(GOBUILD) -o $(DIR_BUILD)/$(PROJECT).linux $(PROJECT);\
	chmod +x $(DIR_BUILD)/$(PROJECT).linux; \
	GOOS=windows CGO_ENABLED=0  GOARCH=386 $(GOBUILD) -o $(DIR_BUILD)/$(PROJECT).exe $(PROJECT); \
	GOOS=darwin CGO_ENABLED=0  GOARCH=amd64 $(GOBUILD) -o $(DIR_BUILD)/$(PROJECT).darwin $(PROJECT);


.PHONY: run
run: ; @ ## Run the application
	$(GORUN) $(LDFLAGS) src/$(PROJECT)/main.go ${ARGS} -config=$(CURDIR)/resources/config/mollydb.json


.PHONY: docker-build
docker-build: build ; @ ## Build Docker images for mollydb version
	docker build -t $(DOCKER_IMAGE):local -f Dockerfile .

.PHONY: docker-publish
docker-publish:  ; @ ## Build Docker images for mollydb version
	docker login -u "$(DOCKER_USERNAME)" -p "$(DOCKER_PASSWORD)"; \
        docker push $(DOCKER_IMAGE):$(VERSION);

.PHONY: docker-run
docker-run: build;
	$(DOCKER_COMPOSE_UP)

.PHONY: docker-stop
docker-stop: ;
	docker-compose -f $(DOCKER_COMPOSE) down -v --remove-orphans;


.PHONY: version
version: ; @ ## Current version
	@printf "version $(VERSION)\n"

.PHONY: help
help:
	@printf "\n\033[0;31m---------------------------------\n"
	@printf "\033[0;37m    The Golang Makefile\n"
	@printf "\033[0;31m---------------------------------\n\n"
	@printf "\033[0;31mCommands\n\n"
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[32m%-25s\033[0m %s\n", $$1,$$2}'
	@printf "\n"
