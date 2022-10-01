DB_CONNECTION=sqlite://db/development.sqlite

migrate-up:
	migrate -database $(DB_CONNECTION) -path ./db/migrations/ up

migrate-down:
	migrate -database $(DB_CONNECTION) -path ./db/migrations/ down
	
