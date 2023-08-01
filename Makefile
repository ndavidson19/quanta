postgres:
	docker run --name postgres15.4 -p 5433:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:15-alpine

start:
	docker start postgres15.4

connect:
	docker exec -it postgres15.4 psql -U root -d go_client

createdb:
	docker exec -it postgres15.4 createdb --username=root --owner=root go_client

dropdb:
	docker exec -it postgres15.4 dropdb go_client

migrateup:
	migrate -path db-migration -database "postgresql://root:secret@localhost:5433/go_client?sslmode=disable" -verbose up

migratedown:
	migrate -path db-migration -database "postgresql://root:secret@localhost:5433/go_client?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server: 
	go run main.go

.PHONY: createdb dropdb migrateup migratedown postgres sqlc test start connect server