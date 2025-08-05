include .env

export

.PHONY: help create-migration migrate-up migrate-down migrate-force

help: ## Show help
	@echo ""
	@echo "\033[1mAvailable commands:\033[0m"
	@awk 'BEGIN {FS = ":.*##";} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

create-migration: ## Create an empty migration
	@read -p "Enter the sequence name: " SEQ; \
	migrate create -ext sql -dir ./db/migrations -seq $${SEQ}

migrate-up: ## Migration up
	@migrate -path=./db/migrations -database $(DATABASE_URL) up

migrate-down: ## Migration down
	@read -p "Number of migrations you want to rollback (default: 1): " NUM; NUM=$${NUM:-1}; \
	migrate -path=./db/migrations -database $(DATABASE_URL) down $${NUM}

migrate-force: ## Migration force version
	@read -p "Enter the version to force: " VERSION; \
	migrate -path=./db/migrations -database $(DATABASE_URL) force $${VERSION}