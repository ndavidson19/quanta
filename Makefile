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

migrateup1:
	migrate -path db-migration -database "postgresql://root:secret@localhost:5433/go_client?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db-migration -database "postgresql://root:secret@localhost:5433/go_client?sslmode=disable" -verbose down

migratedown1:
	migrate -path db-migration -database "postgresql://root:secret@localhost:5433/go_client?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server: 
	go run main.go

mock: 
	mockgen -package mockdb -destination db/mock/store.go github.com/ndavidson19/quanta-backend/db Store 

.PHONY: createdb dropdb migrateup migratedown postgres sqlc test start connect server mock migratedown1 migrateup1