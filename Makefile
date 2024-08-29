DB_URL=postgresql://root:root@localhost:5432/bank?sslmode=disable

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover -short ./...

server:
	go run main.go

mockgen:
	mockgen -package mockdb -destination /home/sayed/go/src/github.com/sayedppqq/banking-backend/db/mock/store.go -build_flags=--mod=mod github.com/sayedppqq/banking-backend/db/sqlc Store

.PHONY: migrateup migratedown sqlc test server mockgen