gen_proto:
	protoc --go_out=ohs/local/pl --go-grpc_out=ohs/local/pl ohs/local/pl/*.proto