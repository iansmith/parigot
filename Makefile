API_PROTO=$(shell find api/proto -type f -regex ".*\.proto")

EXAMPLE_PROTO=$(shell find example -type f -regex ".*\.proto")

TEST_PROTO=$(shell find test -type f -regex ".*\.proto")

## we just use a single representative file for all the generated code
g/file/file.pb.go: $(API_PROTO) $(EXAMPLE_PROTO) $(TEST_PROTO)
	@echo
	@echo "api -->" $(API_PROTO)
	@echo "test -->" $(TEST_PROTO)
	@echo "example -->" $(EXAMPLE_PROTO)
	@echo
	buf lint
	buf generate

.PHONY: protoclean
protoclean: 
	rm -rf g/*
