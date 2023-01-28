
# include scripts/*.mk

.PHONY: run
run:
	@go run \
		cmd/main/main.go \
		--config config.yaml

.PHONY: mocks
mocks:
	@mockery

plugins:
	make plugins:proto

plugins\:proto:
	@protoc \
		--proto_path proto/ \
		--go_out . \
		--go-grpc_out . \
		--govalidators_out . \
		--doc_out ./proto \
		--doc_opt html,index.html \
		proto/plugin/*.proto
