postgres:
	docker run --name postgres88 -p 5432:5432 -e POSTGRES_USER=nader -e POSTGRES_PASSWORD=nader123 -e POSTGRES_DB=billing_system -d postgres:latest

migrateup:
	migrate -path db/migrations -database "postgresql://nader:nader123@localhost:5432/billing_system?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migrations -database "postgresql://nader:nader123@localhost:5432/billing_system?sslmode=disable" -verbose down

sqlc:
	sqlc generate

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/naderSameh/billing_system/db/sqlc Store
	# mockgen -package mockwk -destination worker/mock/distributor.go github.com/naderSameh/billing_system/worker TaskDistributor

test:
	go test -cover -short ./...

swag:
	swag init --parseDependency  --parseInternal -g main.go

	
new_migration:
	migrate create -ext sql -dir db/migrations -seq $(name)

redis:
	docker run --name redis -p 6379:6379 -d redis:7-alpine

.PHONY:
	migrateup migratedown postgres sqlc mock swag redis test new_migration