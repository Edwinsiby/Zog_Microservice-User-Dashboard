proto:
	protoc -I . pb/service.proto --go_out=. --go-grpc_out=.
run:
	go run cmd/main.go