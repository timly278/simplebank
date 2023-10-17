postgres:
	docker run --name postgrestulb2 -p 5432:5432 -e POSTGRES_PASSWORD=tulb -e  POSTGRES_USER=root -d postgres

createdb:
	docker exec -it postgrestulb2 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgrestulb2 createdb dropdb simple_bank

execdb:
	docker exec -it postgrestulb2 psql -U root simple_bank

migratecreate:
	migrate create -ext sql -dir db/migration -seq init_schema

migrateupdate:
	migrate create -ext sql -dir db/migration -seq update_schema

migrateup1:
	migrate -path db/migration -database "postgresql://root:tulb@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migratedown1:
	migrate -path db/migration -database "postgresql://root:tulb@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

.PHONY: postgres createdb dropdb migrateup1 migratedown1 sqlc migratecreate server

############################## TuLB noted ######################################
# docker exec -it postgrestulb2 psql -U root simple_bank
# if come across error when migration:
# select * from schema_migrations; // to get the version
# update schema_migrations set dirty =false where version=XXXX;

# force version of schema_migrations to 1
# migrate -path db/migration -database "postgresql://root:tulb@localhost:5432/simple_bank?sslmode=disable" -verbose force 1
