.PHONY: deps run-cpu run-gpu build-cpu build-gpu test gen build-llama build-llama-cublas

deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
	&& go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	$(MAKE) -C llama -f Makefile deps

run-cpu:
	go run -tags="llama" ./cmd/llm-runner serve

run-gpu:
	go run -tags="llama,nvidia" ./cmd/llm-runner serve

build-cpu:
	@mkdir -p build
	go build -tags="llama" -o build/llm-runner ./cmd/llm-runner

build-gpu:
	@mkdir -p build
	go build -tags="llama,nvidia" -o build/llm-runner ./cmd/llm-runner

test:
	go test ./...

gen:
	mkdir -p ./pb/llmrunnerpb
	protoc --proto_path=./ \
		--go_out=paths=source_relative:./pb/llmrunnerpb \
		--go-grpc_out=paths=source_relative:./pb/llmrunnerpb \
		./llmrunner.proto

build-llama:
	$(MAKE) -C llama libllama.a

build-llama-cublas:
	$(MAKE) -C llama libllama.a BUILD_TYPE=cublas
