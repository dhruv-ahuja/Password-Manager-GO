include .env 
export 

migrateup:
	migrate -path db/migration -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(HOST):$(PORT)/$(DB_NAME)?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(HOST):$(PORT)/$(DB_NAME)?sslmode=disable" -verbose down -all