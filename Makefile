include .env
# .PHONY: clean test security build run

#========================#
#== DATABASE MIGRATION ==#
#========================#

# APP_NAME = apiserver
# BUILD_DIR = $(PWD)/build
MIGRATIONS_FOLDER = ./backend/platform/migrations
SEEDS_FOLDER = ./backend/platform/seeds
name = 1


# migrate create -ext sql -dir db/migrations create_posts_table
# make migrate.create name=create_new_migration
migrate.create:
# docker compose -f ${DOCKER_COMPOSE_FILE} --profile tools run --rm migrate create -ext sql -dir $(MIGRATIONS_FOLDER) ${name}
	migrate create -ext sql -dir $(MIGRATIONS_FOLDER) ${name}

## Run migrations UP
migrate.up:
	docker compose -f $(DOCKER_COMPOSE_FILE) --profile tools run --rm migrate up
# migrate -path $(MIGRATIONS_FOLDER) -database ${DATABASE_URL} up

## Run migrations DOWN (ROLLBACK)
migrate.down:
	docker compose -f $(DOCKER_COMPOSE_FILE) --profile tools run --rm migrate down
# migrate -path $(MIGRATIONS_FOLDER) -database $(DATABASE_URL) down

## Run migrations FORCE
migrate.force:
	docker compose -f $(DOCKER_COMPOSE_FILE) --profile tools run --rm migrate force $(name)
#migrate -path $(MIGRATIONS_FOLDER) -database $(DATABASE_URL) force $(name)

shell.db:
	docker compose -f ${DOCKER_COMPOSE_FILE} exec postgres psql -U $(POSTGRES_USER) -d $(POSTGRES_DB)

lang.start:
	cd resources/langs/ ; goi18n merge *.toml

lang.merge:
	cd resources/langs/ ; goi18n merge translate.*.toml active.*.toml ; cat active.uk.toml > uk.toml ; cat active.ru.toml > ru.toml

lang.del:
	find ./ -type f \( -iname translate.\* -o -iname active.\* \) -delete -print
	
# make docker.build.start ver=1.005
docker.build.start:
	go mod tidy ; docker login ; docker build --no-cache -t api-hackstay:$(ver) -f Dockerfile.production . ; docker tag api-hackstay:$(ver) anrysys/api-hackstay:$(ver) ; docker push anrysys/api-hackstay:$(ver)
#	echo $(ver); echo $(ver); echo $(ver);
# seed:
# 	PGPASSWORD=$(POSTGRES_PASSWORD) psql -h ${POSTGRES_HOST} -p ${POSTGRES_PORT} -U$(POSTGRES_USER) -d $(POSTGRES_DB) -a -f $(SEEDS_FOLDER)/001_seed_user_table.sql
# 	PGPASSWORD=$(POSTGRES_PASSWORD) psql -h ${POSTGRES_HOST} -p ${POSTGRES_PORT} -U$(POSTGRES_USER) -d $(POSTGRES_DB) -a -f $(SEEDS_FOLDER)/002_seed_book_table.sql

# dev: 
# 	docker-compose up -d

# dev-down:
# 	docker-compose down

# start-server:
# 	air

# install-modules:
# 	go get github.com/go-playground/validator/v10
# 	go get -u github.com/gofiber/fiber/v2
# 	go get -u github.com/golang-jwt/jwt/v4
# 	go get github.com/redis/go-redis/v9
# 	go get github.com/satori/go.uuid
# 	go get github.com/spf13/viper
# 	go get gorm.io/driver/postgres
# 	go get -u gorm.io/gorm