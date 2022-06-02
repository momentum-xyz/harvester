.PHONY: gen
req:
	go get entgo.io/ent/cmd/ent
	go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	
up:
	docker-compose up -d

db-reset:
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

db-make-migrations:
	mkdir -p ./ent/migrations
	go run -mod=mod ./cmd/db_make_migrations/main.go $(name)

db-migrate-up:
	migrate \
		-source $(if $(source:-=),$(source),'file://ent/migrations') \
		-database $(if $(database:-=),$(database),'mysql://root@tcp(localhost:3306)/harvester_dev?multiStatements=true') up

db-migrate-down:
	migrate \
		-source $(if $(source:-=),$(source),'file://ent/migrations') \
		-database $(if $(database:-=),$(database),'mysql://root@tcp(localhost:3306)/harvester_dev') down
