.PHONY: run docker protoc compose

protoc:
	protoc \
    --proto_path=. \
    --proto_path=api/proto/third_party/googleapis \
    --proto_path=$(shell go list -m -f '{{.Dir}}' github.com/grpc-ecosystem/grpc-gateway/v2) \
	--proto_path=api/proto/third_party/protobuf \
    --go_out=. \
    --go-grpc_out=. \
    --grpc-gateway_out=. \
    --openapiv2_out=. \
	--openapiv2_opt=logtostderr=true,file=pvz_service.swagger.json \
    api/proto/pvz.proto

docker:
	docker compose up --build

run:
	go run cmd/main/main.go || true

compose:
	docker compose up

k6:
    k6 run scripts/k6/script.js