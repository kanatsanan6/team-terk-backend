migrateup:
	migrate -path db/migrations -database "postgresql://postgres:@localhost:5432/go_backend_development?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migrations -database "postgresql://postgres:@localhost:5432/go_backend_development?sslmode=disable" -verbose down

.PHONY: migrateup migratedown
