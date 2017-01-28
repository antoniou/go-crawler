LINTIGNOREDOT='awstesting/integration.+should not use dot imports'
LINTIGNOREDOC='service/[^/]+/(api|service|waiters)\.go:.+(comment on exported|should have comment or be unexported)'
LINTIGNORECONST='service/[^/]+/(api|service|waiters)\.go:.+(type|struct field|const|func) ([^ ]+) should be ([^ ]+)'
LINTIGNORESTUTTER='service/[^/]+/(api|service)\.go:.+(and that stutters)'
LINTIGNOREINFLECT='service/[^/]+/(api|service)\.go:.+method .+ should be '

TOOL_ONLY_PKGS=$(shell go list ./... | grep -v "/vendor/")
SDK_GO_1_4=$(shell go version | grep "go1.4")
SDK_GO_1_5=$(shell go version | grep "go1.5")
SDK_GO_VERSION=$(shell go version | awk '''{print $$3}''' | tr -d '''\n''')

all: get-deps

help:
	@echo "Please use \`make <target>' where <target> is one of"
	@echo "  get-deps                to go get the SDK dependencies"
	@echo "  get-deps-tests          to get the SDK's test dependencies"
	@echo "  get-deps-verify         to get the SDK's verification dependencies"

get-deps: generate get-deps-tests get-deps-verify
	@echo "go get dependencies"
	@go get -v $(TOOL_ONLY_PKGS)

get-deps-tests:
	@echo "go get testing dependencies"
	# go get github.com/gucumber/gucumber/cmd/gucumber
	go get github.com/stretchr/testify
	# go get github.com/smartystreets/goconvey
	# go get golang.org/x/net/html

get-deps-verify:
	@echo "go get verification utilities"
	@if [ \( -z "${SDK_GO_1_4}" \) -a \( -z "${SDK_GO_1_5}" \) ]; then  go get github.com/golang/lint/golint; else echo "skipped getting golint"; fi

generate: gen-bindata

gen-bindata:
	go get -u github.com/jteeuwen/go-bindata/...
	go generate main.go

build:
	@echo "go build SDK and vendor packages"
	@go build ${TOOL_ONLY_PKGS}

verify: lint

lint:
	@echo "go lint SDK and vendor packages"
	@lint=`if [ \( -z "${SDK_GO_1_4}" \) -a \( -z "${SDK_GO_1_5}" \) ]; then  golint ./...; else echo "skipping golint"; fi`; \

unit: get-deps-tests build verify
		@echo "go test SDK and vendor packages"
		# @go test -tags ${UNIT_TEST_TAGS} $(SDK_UNIT_TEST_ONLY_PKGS)
		@go test ./domain -cover -v
