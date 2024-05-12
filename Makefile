watch:
	@docker compose -f compose.dev.yml up --build

up:
	@docker compose -f compose.dev.yml up -d --build

down:
	@docker compose -f compose.dev.yml down -v

restart: down up

build-account:
	@docker compose -f compose.dev.yml build account_manager

build-payment:
	@docker compose -f compose.dev.yml build payment_manager