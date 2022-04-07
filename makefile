gen_proto:
	protoc --go_out=ohs/local/pl --go-grpc_out=ohs/local/pl --grpc-gateway_out=ohs/local/pl ohs/local/pl/*.proto

server:
	go run main.go

client1:
	go run client/main.go

migrate.create:
	migrate create -ext sql -dir ./migrations -seq init_table_schema

migrate.up:
	migrate -path ./migrations -database "postgresql://postgres:driver@localhost:5432/order?sslmode=disable" -verbose up

migrate.down:
	migrate -path ./migrations -database "postgresql://postgres:driver@localhost:5432/order?sslmode=disable" -verbose down

test:
	go test -v -coverprofile cover.out ./...
	go tool cover -func cover.out
	go tool cover -html=cover.out -o cover.html

.PHONY: gen_proto gen_errors server client1 migrate.up migrate.down migrate.create get