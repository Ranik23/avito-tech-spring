protoc:
	protoc \
  --proto_path=. \
  --proto_path=$(go list -f '{{ .Dir }}' -m google.golang.org/protobuf)/../.. \
  --go_out=. \
  --go-grpc_out=. \
  api/proto/pvz.proto
