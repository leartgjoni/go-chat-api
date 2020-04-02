start-local:
	CONFIG_PATH=local.env go run cmd/app/main.go
init-db:
	docker-compose -f scripts/env/docker-compose.yml up -d
	docker-compose -f scripts/env/docker-compose-test.yml up -d
init-ci-env:
	docker-compose -f scripts/env/docker-compose-test.yml up -d
test:
	ENV_FILE=test.env go test -v ./...
test-unit:
	go test -v -short ./...
test-integration:
	ENV_FILE=test.env go test -run Integration ./...
test-coverage:
	ENV_FILE=test.env go test -race -covermode atomic -coverprofile=profile.cov ./...
test-html-coverage:
	ENV_FILE=test.env go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
test-unit-coverage:
	go test -coverprofile=coverage-unit.out -short ./...
	go tool cover -html=coverage-unit.out
test-integration-coverage:
	ENV_FILE=test.env go test -coverprofile=coverage-integration.out -run Integration ./...
	go tool cover -html=coverage-integration.out