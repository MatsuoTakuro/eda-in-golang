install.tools:
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
PACT_BROKER_URL ?= http://localhost:9292
PACT_BROKER_USERNAME ?= pactuser
PACT_BROKER_PASSWORD ?= pactpass

.PHONY: publish.pacts
publish.pacts:
	@echo "🔍 Searching for pact directories (excluding node_modules)..."
	@find modules \
		-path "*/node_modules" -prune -o \
		-type d -name pacts -print | while read dir; do \
			echo "📦 Publishing pacts in $$dir..."; \
			pact-broker publish $$dir \
				--consumer-app-version $(PUBLISH_VERSION) \
				--broker-base-url $(PACT_BROKER_URL) \
				--broker-username $(PACT_BROKER_USERNAME) \
				--broker-password $(PACT_BROKER_PASSWORD); \
		done

test.all: publish.pacts
	@echo "Running all tests..."
	go test -tags="integration e2e" ./...
	@echo "Tests completed."

build.monolith:
	docker compose --profile monolith build

up.monolith:
	docker compose --profile monolith up -d

build.microservices:
	docker compose --profile microservices build

up.microservices:
	docker compose --profile microservices up -d