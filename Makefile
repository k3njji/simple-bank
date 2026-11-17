
createdb:
	docker exec -e PGPASSWORD=secret postgres12 createdb -h 127.0.0.1 -U root simple_bank

dropdb:
	docker exec -e PGPASSWORD=secret postgres12 dropdb -h 127.0.0.1 -U root --if-exists simple_bank

migrateup:
	migrate -path db/migrate -database "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migrate -database "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: createdb dropdb migrateup migratedown sqlc test
