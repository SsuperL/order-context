gen_proto:
	protoc --go_out=ohs/local/pl --go-grpc_out=ohs/local/pl ohs/local/pl/*.proto

gen_errors:
	protoc --go_out=ohs/local/pl/errors --go-grpc_out=ohs/local/pl/errors ohs/local/pl/errors/*.proto

server:
	go run main.go

client1:
	go run client/main.go
