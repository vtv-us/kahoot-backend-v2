postgres:
	docker run --name kahoot -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres:12-alpine
createdb:
	docker exec -it kahoot createdb --username=postgres --owner=postgres kahoot
dropdb:
	docker exec -it kahoot dropdb --username=postgres kahoot
migrateup:
	migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/kahoot?sslmode=disable" -verbose up
migrateup1:
	migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/kahoot?sslmode=disable" -verbose up 1
migratedown:
	migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/kahoot?sslmode=disable" -verbose down
migratedown1:
	migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/kahoot?sslmode=disable" -verbose down 1
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
