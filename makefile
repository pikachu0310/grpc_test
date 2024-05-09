.PHONY: proto

SERVER_OUTPUT_DIR=./server/proto
CLIENT_OUTPUT_DIR=./client

proto:
	@echo "Generating gRPC code..."
	@protoc --proto_path=proto pingpong.proto \
		--go_out=${SERVER_OUTPUT_DIR} --go_opt=paths=source_relative \
		--go-grpc_out=${SERVER_OUTPUT_DIR} --go-grpc_opt=paths=source_relative \
		--js_out=import_style=commonjs:${CLIENT_OUTPUT_DIR} \
		--grpc-web_out=import_style=commonjs,mode=grpcweb:${CLIENT_OUTPUT_DIR}


client-build:
	npm run --prefix ./client build

client-run:
	npx http-server ./client -p 8080 & echo $$! > client-server.pid

client-stop:
	kill `cat client-server.pid` && rm client-server.pid

server-run:
	go run main.go & echo $$! > server.pid

server-stop:
	kill `cat server.pid` && rm server.pid
