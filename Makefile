# PROTOS := $(shell find ./divar_interface -name *.proto)
TEST_PROTOS := $(shell find ./test/divar_interface -name *.proto)

build:
	go build -o $(GOPATH)/bin/protoc-gen-divar-doc ./cmd/protoc-gen-divar-doc/...

generate:
	@protoc --divar-doc_out=. $(PROTOS)

test-generate: build
	@protoc --divar-doc_out=. --divar-doc_opt=exclude=__ALL__ --proto_path=test $(TEST_PROTOS)