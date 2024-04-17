migration-up:
	migrate -path ./migrations/postgres -database 'postgres://muhiddin:1@localhost:5432/postgres?sslmode=disable' up
	
migration-down:
	migrate -path ./migrations/postgres -database 'postgres://muhiddin:1@localhost:5432/postgres?sslmode=disable' down
	
migration-force-1v:
	migrate -path ./migrations/postgres -database 'postgres://muhiddin:1@localhost:5432/postgres?sslmode=disable' force 1
	

