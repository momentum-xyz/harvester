.PHONY: gen
prerequisites:
	go install ariga.io/entimport/cmd/entimport@latest
	go install entgo.io/ent/cmd/ent@latest
	go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	
up:
	docker-compose up -d

db-reset:
	make db-recreate
	sleep 10
	make db-migrate-up

db-recreate:
	docker-compose rm --stop --force db
	docker-compose up -d db


run:
	go run ./cmd/harvester

test:
	rm -rf coverage; mkdir -p coverage
	go clean -testcache
	go test ./$(dir)... -v -coverpkg=./$(dir)... -coverprofile=./coverage/cover.out
	go tool cover -html=./coverage/cover.out -o ./coverage/cover.html

gen: 
	go generate ./...

db-migrations:
	mkdir -p ./migrations
	go run -mod=mod ./cmd/db_make_migrations/main.go $(name)

db-migrate-up:
	migrate \
		-source $(if $(source:-=),$(source),'file://migrations') \
		-database $(if $(database:-=),$(database),'mysql://root@tcp(localhost:3306)/harvester?multiStatements=true') up

db-migrate-down:
	migrate \
		-source $(if $(source:-=),$(source),'file://migrations') \
		-database $(if $(database:-=),$(database),'mysql://root@tcp(localhost:3306)/harvester') down

db-migrate-dc:
	 migrate \
		-source $(if $(source:-=),$(source),'file://migrations') \
		-database $(if $(database:-=),$(database),'mysql://${DB_USERNAME}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_DATABASE}?multiStatements=true') up
