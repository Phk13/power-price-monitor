.PHONY: help
help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: run
run: ## Run the application
	go run ./cmd/main.go

.PHONY: lint
lint: ## Run linters
	$$(go env GOPATH)/bin/golangci-lint run --color always --timeout 5m

.PHONY: docker-build
docker-build: ## Build Docker image
	docker build -t power-price-monitor .

.PHONY: docker-run
docker-run: ## Run Docker container
	docker run --rm -p 8080:8080 --env-file .env power-price-monitor

.PHONY: docker-compose-up
docker-compose-up: ## Start services with docker-compose
	docker-compose up -d

.PHONY: docker-compose-down
docker-compose-down: ## Stop services with docker-compose
	docker-compose down

.PHONY: docker-compose-logs
docker-compose-logs: ## Show docker-compose logs
	docker-compose logs -f

.PHONY: docker-clean
docker-clean: ## Remove Docker images and containers
	docker-compose down --rmi all --volumes --remove-orphans