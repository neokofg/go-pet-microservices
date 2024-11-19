PROTO_DIR = internal/proto
GO_OUT_DIR = internal/proto

.PHONY: proto
proto:
    # Установка необходимых инструментов
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

    # Генерация Go кода из proto файлов
    protoc --proto_path=$(PROTO_DIR) \
           --go_out=$(GO_OUT_DIR) --go_opt=paths=source_relative \
           --go-grpc_out=$(GO_OUT_DIR) --go-grpc_opt=paths=source_relative \
           $(PROTO_DIR)/**/*.proto

.PHONY: build-gateway
build-gateway:
    go build -o ./bin/api-gateway ./internal/api-gateway/cmd/main.go

.PHONY: build-catalog
build-catalog:
    go build -o ./bin/catalog-service ./internal/catalog-service/cmd/main.go

.PHONY: build
build: proto build-gateway build-catalog

.PHONY: run
run:
    docker-compose up --build