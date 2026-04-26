.PHONY: deps gen build-libs-cpu build-libs-gpu build-cpu build-gpu run-cpu run-gpu test-llama-cpu test-llama-gpu test clean

LLAMA_DIR := llama
RUN_ENV := LD_LIBRARY_PATH="$(PWD)/$(LLAMA_DIR):$(LD_LIBRARY_PATH)" LIBRARY_PATH="$(PWD)/$(LLAMA_DIR):$(LIBRARY_PATH)"
deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
	&& go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	$(MAKE) -C $(LLAMA_DIR) -f Makefile deps

gen:
	mkdir -p ./pb/llmrunnerpb
	protoc --proto_path=./ \
		--go_out=paths=source_relative:./pb/llmrunnerpb \
		--go-grpc_out=paths=source_relative:./pb/llmrunnerpb \
		./llmrunner.proto

test:
	go test ./...

test-llama-cpu: build-libs-cpu
	$(RUN_ENV) go test -tags="llama" ./provider ./service

test-llama-gpu:
	@if command -v nvidia-smi >/dev/null 2>&1 && $(MAKE) build-libs-gpu >/dev/null 2>&1; then \
		$(RUN_ENV) go test -tags="llama,nvidia" ./provider ./service; \
	else \
		echo "Пропуск test-llama-gpu: библиотеки GPU/CUDA недоступны в этой среде"; \
	fi

build-libs-cpu:
	$(MAKE) -C $(LLAMA_DIR) libbinding.a
	ln -sf libllama.so $(LLAMA_DIR)/libllama.so.0
	ln -sf libggml.so $(LLAMA_DIR)/libggml.so.0
	ln -sf libggml-base.so $(LLAMA_DIR)/libggml-base.so.0
	ln -sf libggml-cpu.so $(LLAMA_DIR)/libggml-cpu.so.0

build-libs-gpu:
	$(MAKE) -C $(LLAMA_DIR) libbinding.a BUILD_TYPE=cublas
	ln -sf libllama.so $(LLAMA_DIR)/libllama.so.0
	ln -sf libggml.so $(LLAMA_DIR)/libggml.so.0
	ln -sf libggml-base.so $(LLAMA_DIR)/libggml-base.so.0
	ln -sf libggml-cpu.so $(LLAMA_DIR)/libggml-cpu.so.0
	ln -sf libggml-cuda.so $(LLAMA_DIR)/libggml-cuda.so.0

build-cpu: build-libs-cpu
	@mkdir -p build
	go build -tags="llama" -o build/gen-runner ./cmd/gen-runner

build-gpu: build-libs-gpu
	@mkdir -p build
	go build -tags="llama,nvidia" -o build/gen-runner ./cmd/gen-runner

run-cpu: build-libs-cpu
	$(RUN_ENV) go run -tags="llama" ./cmd/gen-runner serve

run-gpu: build-libs-gpu
	$(RUN_ENV) go run -tags="llama,nvidia" ./cmd/gen-runner serve

clean:
	rm -rf build
