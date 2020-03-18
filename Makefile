start-local:
	CONFIG_PATH=local.env go run cmd/app/main.go
init-db:
	docker-compose -f scripts/env/docker-compose.yml up -d