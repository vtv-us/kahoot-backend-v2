postgres:
	docker run --name postgres-bank -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres:12-alpine
createdb:
	docker exec -it postgres-bank createdb --username=postgres --owner=postgres simple_bank
dropdb:
	docker exec -it postgres-bank dropdb --username=postgres simple_bank
migrateup:
	migrate -path db/migration -database "postgresql://postgres:postgres@localhost:2345/simple_bank?sslmode=disable" -verbose up
migrateup1:
	migrate -path db/migration -database "postgresql://postgres:postgres@localhost:2345/simple_bank?sslmode=disable" -verbose up 1
migratedown:
	migrate -path db/migration -database "postgresql://postgres:postgres@localhost:2345/simple_bank?sslmode=disable" -verbose down
migratedown1:
	migrate -path db/migration -database "postgresql://postgres:postgres@localhost:2345/simple_bank?sslmode=disable" -verbose down 1
sqlc:
	sqlc generate
	sed -i 's/repositories/entities/g' ./internal/entities/models.go
test:
	go test -v -cover ./...
server: 
	go run main.go
mock:
	mockgen -package mockdb -destination mock/store.go github.com/vtv-us/kahoot-backend/internal/repositories Store

.PHONY: postgres createdb dropdb migrateup migratedown migrateup1 migratedown1 sqlc test server mock
