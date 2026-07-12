.PHONY: up
up: ## Build and start the entire stack (App, DB, Redis)
	@echo "Starting the entire Docker stack..."
	docker-compose up --build -d