.PHONY: run docker protoc compose

protoc:
	protoc \
		--proto_path=. \
		--proto_path=$$(go list -f '{{ .Dir }}' -m google.golang.org/protobuf)/../.. \
		--go_out=. \
		--go-grpc_out=. \
		api/proto/pvz.proto

docker:
	docker build -t my-app .

run:
	-go run cmd/main/main.go || true

compose:
	docker compose up
