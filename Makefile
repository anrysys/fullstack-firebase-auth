# include .env
include backend/.env
# export $(shell sed 's/=.*//' backend/.env)
# .PHONY: clean test security build run

#========================#
#== DATABASE MIGRATION ==#
#========================#

# APP_NAME = apiserver
# BUILD_DIR = $(PWD)/build

MIGRATIONS_FOLDER=./backend/platform/migrations
SEEDS_FOLDER=./backend/platform/seeds
name=1
DOCKER_RUN_COMMAND=docker compose -f $(DOCKER_COMPOSE_FILE) --env-file ./backend/.env --profile tools run --rm

########## DOCKER ##########

docker.start:
	docker compose --env-file ./backend/.env up -d

docker.start.dev:
	docker compose --env-file ./backend/.env up

# make docker.build.start ver=1.005
docker.build.start:
	cd backend ; go mod tidy ; docker login ; docker build --no-cache -t api-hackstay:$(ver) -f Dockerfile.production . ; docker tag api-hackstay:$(ver) anrysys/api-hackstay:$(ver) ; docker push anrysys/api-hackstay:$(ver)
#	echo $(ver); echo $(ver); echo $(ver);		

########## BACKEND ##########

start.server.backend:
	cd backend ; air 

start.server.backend.prod:
	cd backend ; go run main.go


backend.init:
	cd backend ; go mod init github.com/anrysys/fullstack-firebase-auth
	go mod tidy
	go mod vendor

backend.install:
	cd backend ; go mod download
	go mod download

backend.update:
	cd backend ; go get -u
	go mod tidy
	go mod vendor

lang.start:
	cd backend/resources/langs/ ; goi18n merge *.toml

lang.merge:
	cd backend/resources/langs/ ; goi18n merge translate.*.toml active.*.toml ; cat active.uk.toml > uk.toml ; cat active.ru.toml > ru.toml

lang.del:
	cd backend ; find ./ -type f \( -iname translate.\* -o -iname active.\* \) -delete -print
	
lang.extract:
	cd backend ; goi18n extract -outdir=resources/langs -sourceLanguage=en

lang.update:
	cd backend ; goi18n -sourceLanguage=en -outdir=resources/langs active.en.toml active.ru.toml active.uk.toml

lang.init:
	cd backend ; goi18n -sourceLanguage=en -outdir=resources/langs active.en.toml active.ru.toml active.uk.toml

lang.init.en:
	cd backend ; goi18n -sourceLanguage=en -outdir=resources/langs active.en.toml

lang.init.ru:
	cd backend ; goi18n -sourceLanguage=ru -outdir=resources/langs active.ru.toml

lang.init.uk:
	cd backend ; goi18n -sourceLanguage=uk -outdir=resources/langs active.uk.toml

lang.init.all:
	cd backend ; goi18n -sourceLanguage=en -outdir=resources/langs active.en.toml active.ru.toml active.uk.toml

########## MIGRATIONS ##########
#
# https://github.com/golang-migrate/migrate
# 
# Tutorials: https://github.com/golang-migrate/migrate/blob/master/database/postgres/TUTORIAL.md
# Parameters: https://github.com/golang-migrate/migrate/tree/master/database/postgres
#
# migrate create -ext sql -dir db/migrations create_posts_table
# make migrate.create create_posts_table
migrate.create:
	migrate create -ext sql -dir $(MIGRATIONS_FOLDER) $@

## Run migrations UP/DOWN with ARGS (make migrate up 2 or make migrate down 1)
ARGS = $(filter-out $@,$(MAKECMDGOALS))
migrate:
	$(DOCKER_RUN_COMMAND) migrate -source file://migrations -database postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@postgres:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=$(POSTGRES_SSL_MODE) $(ARGS)
	@:

## Run migrations FORCE
migrate.force:
	$(DOCKER_RUN_COMMAND) migrate force $(name)

shell.db:
	docker compose -f $(DOCKER_COMPOSE_FILE) exec postgres psql -U $(POSTGRES_USER) -d $(POSTGRES_DB)

# seed:
# 	PGPASSWORD=$(POSTGRES_PASSWORD) psql -h $(POSTGRES_HOST) -p $(POSTGRES_PORT) -U$(POSTGRES_USER) -d $(POSTGRES_DB) -a -f $(SEEDS_FOLDER)/001_seed_user_table.sql
# 	PGPASSWORD=$(POSTGRES_PASSWORD) psql -h $(POSTGRES_HOST) -p $(POSTGRES_PORT) -U$(POSTGRES_USER) -d $(POSTGRES_DB) -a -f $(SEEDS_FOLDER)/002_seed_book_table.sql
%:
	@:
