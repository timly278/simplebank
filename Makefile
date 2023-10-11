postgres:
	docker run --name postgrestulb2 -p 5432:5432 -e POSTGRES_PASSWORD=tulb -e  POSTGRES_USER=root -d postgres

createdb:
	docker exec -it postgrestulb2 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgrestulb2 createdb dropdb simple_bank

execdb:
	docker exec -it postgrestulb2 psql -U root simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:tulb@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:tulb@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrateup migratedown sqlc


# docker exec -it postgrestulb2 psql -U root simple_bank
# if come across error when migration:
# select * from schema_migrations; // to get the version
# update schema_migrations set dirty =false where version=XXXX;