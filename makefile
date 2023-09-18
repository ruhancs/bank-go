postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=123456 -d postgres

createdb:
	docker exec -it postgres createdb --username=postgres --owner=postgres bank

dropdb:
	docker exec -it postgres dropdb bank

migrateup:
	migrate -path db/migration -database "postgresql://postgres:123456@localhost:5432/bank?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://postgres:123456@localhost:5432/bank?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://postgres:123456@localhost:5432/bank?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://postgres:123456@localhost:5432/bank?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

proto:
	rm -rf pb/*.go
	rm -rf doc/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
    --openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=bank \
    proto/*.proto

evans:
	evans --host localhost --port 8001 -r repl

redis:
	docker run --name redis -p 6379:6379 -d redis:7-alpine

.PHONY: postgres createdb dropdb migrateup migrateup1 migratedown migratedown1 sqlc test server proto evans redis