.PHONY: proto
proto:
	protoc -I=./proto \
	       --go_out=./proto --go_opt=paths=source_relative \
	       --go-grpc_out=./proto --go-grpc_opt=paths=source_relative \
	       --js_out=import_style=commonjs:./client \
	       --grpc-web_out=import_style=commonjs,mode=grpcwebtext:./client \
	       proto/pingpong.proto
