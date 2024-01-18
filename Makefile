export GOOSE_DRIVER=postgres
export GOOSE_DBSTRING=postgresql://postgres:example@127.0.0.1:5432/postgres?sslmode=disable

up:
	@goose -dir db/migrations up

down:
	@goose -dir db/migrations down

