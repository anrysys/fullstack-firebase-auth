# include .env
include backend/.env
# export $(shell sed 's/=.*//' backend/.env)
# .PHONY: clean test security build run

# APP_NAME = apiserver
# BUILD_DIR = $(PWD)/build

MIGRATIONS_FOLDER=./backend/platform/migrations
SEEDS_FOLDER=./backend/platform/seeds
name=1
DOCKER_RUN_COMMAND=docker compose -f $(DOCKER_COMPOSE_FILE) --env-file ./backend/.env --profile tools run --rm

################################################## 
#  ____             _             
# |  _ \  ___   ___| | _____ _ __ 
# | | | |/ _ \ / __| |/ / _ \ '__|
# | |_| | (_) | (__|   <  __/ |   
# |____/ \___/ \___|_|\_\___|_|   
#                                

docker.start:
	docker compose --env-file ./backend/.env up -d

docker.start.dev:
	docker compose --env-file ./backend/.env up

# make docker.build.start ver=1.005
docker.build.start:
	cd backend ; go mod tidy ; docker login ; docker build --no-cache -t api-hackstay:$(ver) -f Dockerfile.production . ; docker tag api-hackstay:$(ver) anrysys/api-hackstay:$(ver) ; docker push anrysys/api-hackstay:$(ver)
#	echo $(ver); echo $(ver); echo $(ver);		

##################################################
#  ____             _                  _ 
# | __ )  __ _  ___| | _____ _ __   __| |
# |  _ \ / _` |/ __| |/ / _ \ '_ \ / _` |
# | |_) | (_| | (__|   <  __/ | | | (_| |
# |____/ \__,_|\___|_|\_\___|_| |_|\__,_|
#                                      

start.server.backend.dev:
	cd backend ; air > /dev/null 2>&1 &

stop.server.backend.dev:
	cd backend ; killall air

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

##################################################
# 
#  __  __ _                 _   _                 
# |  \/  (_) __ _ _ __ __ _| |_(_) ___  _ __  ___ 
# | |\/| | |/ _` | '__/ _` | __| |/ _ \| '_ \/ __|
# | |  | | | (_| | | | (_| | |_| | (_) | | | \__ \
# |_|  |_|_|\__, |_|  \__,_|\__|_|\___/|_| |_|___/
#           |___/                                 
#
# https://github.com/golang-migrate/migrate
# 
# Tutorials: https://github.com/golang-migrate/migrate/blob/master/database/postgres/TUTORIAL.md
# Parameters: https://github.com/golang-migrate/migrate/tree/master/database/postgres
#
# make migrate.create name=create_posts_table
# make migrate up 
# make migrate up 2	# Run 2 migrations
# make migrate down 1	# Rollback 1 migration
# make migrate.force 20210901123456_create_posts_table
# 
migrate.create:
	migrate create -ext sql -dir $(MIGRATIONS_FOLDER) $(name)
# migrate create -ext sql -dir ./backend/platform/migrations NAME

## Run migrations UP/DOWN with ARGS (make migrate up 2 or make migrate down 1)
ARGS = $(filter-out $@,$(MAKECMDGOALS))
migrate.run:
	$(DOCKER_RUN_COMMAND) migrate -source file://migrations -database postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@postgres:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=$(POSTGRES_SSL_MODE) $(ARGS)
	@:

## Run migrations FORCE
# make migrate.force v=20240610070532
migrate.force:
	$(DOCKER_RUN_COMMAND) migrate -source file://migrations -database postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@postgres:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=$(POSTGRES_SSL_MODE) force $(v)
# docker compose -f docker-compose.yml --env-file ./backend/.env --profile tools run --rm migrate -source file://migrations -database postgres://yoko:8cC8900ac3a46a3cx4e1e1e64ce7492c4180e5@postgres:5432/yoko?sslmode=disable force 20240610070532	

shell.db:
	docker compose -f $(DOCKER_COMPOSE_FILE) exec postgres psql -U $(POSTGRES_USER) -d $(POSTGRES_DB)

# seed:
# 	PGPASSWORD=$(POSTGRES_PASSWORD) psql -h $(POSTGRES_HOST) -p $(POSTGRES_PORT) -U$(POSTGRES_USER) -d $(POSTGRES_DB) -a -f $(SEEDS_FOLDER)/001_seed_user_table.sql
# 	PGPASSWORD=$(POSTGRES_PASSWORD) psql -h $(POSTGRES_HOST) -p $(POSTGRES_PORT) -U$(POSTGRES_USER) -d $(POSTGRES_DB) -a -f $(SEEDS_FOLDER)/002_seed_book_table.sql
%:
	@:
