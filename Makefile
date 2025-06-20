install-tools:
	@echo installing tools && \
	go install \
	github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
	github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
	google.golang.org/protobuf/cmd/protoc-gen-go \
	google.golang.org/grpc/cmd/protoc-gen-go-grpc
	@echo done

generate:
	@echo running code generation
	go generate ./...
	@echo done


PUBLISH_VERSION ?= 1.0.1
PACT_BROKER_URL ?= http://host.docker.internal:9292

.PHONY: publish.pacts
publish.pacts:
	@echo "🔍 Searching for pact directories (excluding node_modules)..."
	@find modules \
		-path "*/node_modules" -prune -o \
		-type d -name pacts -print | while read dir; do \
			echo "📦 Publishing pacts in $$dir..."; \
			docker run --rm \
				-v $$PWD/$$dir:/pacts \
				pactfoundation/pact-cli:latest \
				publish /pacts \
				--consumer-app-version $(PUBLISH_VERSION) \
				--broker-base-url $(PACT_BROKER_URL); \
		done