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

# swagger:
# 	docker pull swaggerapi/swagger-ui
# 	docker run -p 8085:8080 -e SWAGGER_JSON=/foo/swagger.yaml -v ./api/openapi2/v1/backend.yaml:/foo/swagger.yaml swaggerapi/swagger-ui
