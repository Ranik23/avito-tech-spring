protoc:
	protoc \
  --proto_path=. \
  --proto_path=$(go list -f '{{ .Dir }}' -m google.golang.org/protobuf)/../.. \
  --go_out=. \
  --go-grpc_out=. \
  api/proto/pvz.proto


dto:
  docker run --rm -v ${PWD}:/local openapitools/openapi-generator-cli generate  \
   -i /local/api/openapi/backend.yaml \
    -g go  \
    -o /local/internal/models/dto \ 
    --global-property models,modelDocs=false,modelTests=false,supportingFiles= \ 
    --additional-properties=enumClassPrefix=true