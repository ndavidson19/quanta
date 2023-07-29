postgres:
	docker run --name postgres15 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:15-alpine

start:
	docker start postgres15

connect:
	docker exec -it postgres15 psql -U root -d go_client

createdb:
	docker exec -it postgres15 createdb --username=root --owner=root go_client

dropdb:
	docker exec -it postgres15 dropdb go_client

install:
	curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz sudo mv migrate /usr/bin/ which migrate 

migrateup:
	migrate -path db-migration -database "postgresql://root:secret@localhost:5432/go_client?sslmode=disable" -verbose up

migratedown:
	migrate -path db-migration -database "postgresql://root:secret@localhost:5432/go_client?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: createdb dropdb migrateup migratedown postgres sqlc test start connect