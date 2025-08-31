# Billing Engine Makefile
# This Makefile provisions all necessary resources and sets up the database

# Variables
DOCKER_COMPOSE_FILE := docker-compose.yml
# Detect docker compose command (modern vs legacy)
DOCKER_COMPOSE := $(shell if command -v docker-compose >/dev/null 2>&1; then echo "docker-compose"; else echo "docker compose"; fi)
DB_HOST := localhost
DB_PORT := 3306
DB_NAME := billing_engine
DB_USER := billing_admin
DB_PASSWORD := billing_password
ROOT_PASSWORD := root_billing
GO_VERSION := 1.21
APP_NAME := billing-engine
APP_PORT := 9006

# Colors for output
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[1;33m
BLUE := \033[0;34m
NC := \033[0m # No Color

.PHONY: help setup clean build run test docker-up docker-down db-setup db-migrate db-verify db-reset dev-setup all

# Default target
all: setup

help: ## Show this help message
	@echo "$(YELLOW)Billing Engine - Available Commands:$(NC)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "$(BLUE)%-20s$(NC) %s\n", $$1, $$2}'
	@echo ""

setup: ## Complete setup - provisions all resources and sets up database
	@echo "$(YELLOW)🚀 Setting up Billing Engine...$(NC)"
	@$(MAKE) docker-up
	@sleep 10  # Wait for MySQL to be ready
	@$(MAKE) db-setup
	@$(MAKE) build
	@echo "$(GREEN)✅ Setup completed successfully!$(NC)"
	@echo ""
	@echo "$(YELLOW)Next steps:$(NC)"
	@echo "  • Run '$(BLUE)make run$(NC)' to start the application"
	@echo "  • Run '$(BLUE)make test$(NC)' to run tests"
	@echo "  • API will be available at http://localhost:$(APP_PORT)"

dev-setup: ## Setup development environment with dependencies
	@echo "$(YELLOW)🔧 Setting up development environment...$(NC)"
	@if ! command -v go >/dev/null 2>&1; then \
		echo "$(RED)❌ Go is not installed. Please install Go $(GO_VERSION) or later$(NC)"; \
		exit 1; \
	fi
	@echo "$(GREEN)✅ Go is installed: $$(go version)$(NC)"
	@if ! command -v docker >/dev/null 2>&1; then \
		echo "$(RED)❌ Docker is not installed. Please install Docker$(NC)"; \
		exit 1; \
	fi
	@echo "$(GREEN)✅ Docker is installed: $$(docker --version)$(NC)"
	@if command -v docker-compose >/dev/null 2>&1; then \
		echo "$(GREEN)✅ Docker Compose is installed: $$(docker-compose --version)$(NC)"; \
	elif docker compose version >/dev/null 2>&1; then \
		echo "$(GREEN)✅ Docker Compose is installed: $$(docker compose version)$(NC)"; \
	else \
		echo "$(RED)❌ Docker Compose is not installed. Please install Docker Compose$(NC)"; \
		exit 1; \
	fi
	@echo "$(YELLOW)📦 Installing Go dependencies...$(NC)"
	@go mod download
	@go mod vendor
	@echo "$(GREEN)✅ Development environment setup completed!$(NC)"

docker-up: ## Start Docker containers (MySQL database)
	@echo "$(YELLOW)🐳 Starting Docker containers...$(NC)"
	@$(DOCKER_COMPOSE) up -d
	@echo "$(GREEN)✅ Docker containers started$(NC)"
	@echo "$(YELLOW)⏳ Waiting for MySQL to be ready...$(NC)"
	@for i in $$(seq 1 30); do \
		if $(DOCKER_COMPOSE) exec mysql_db mysqladmin ping -h localhost --silent 2>/dev/null; then \
			break; \
		fi; \
		if [ $$i -eq 30 ]; then \
			echo "$(RED)❌ MySQL failed to start after 60 seconds$(NC)"; \
			exit 1; \
		fi; \
		sleep 2; \
	done
	@echo "$(GREEN)✅ MySQL is ready$(NC)"

docker-down: ## Stop and remove Docker containers
	@echo "$(YELLOW)🛑 Stopping Docker containers...$(NC)"
	@$(DOCKER_COMPOSE) down
	@echo "$(GREEN)✅ Docker containers stopped$(NC)"

docker-logs: ## Show Docker container logs
	@echo "$(YELLOW)📋 Docker container logs:$(NC)"
	@$(DOCKER_COMPOSE) logs -f

db-setup: ## Setup database, create tables, and populate initial data
	@echo "$(YELLOW)🗄️  Setting up database...$(NC)"
	@echo "$(YELLOW)Creating database and user...$(NC)"
	@$(DOCKER_COMPOSE) exec mysql_db mysql -u root -p$(ROOT_PASSWORD) -e "\
		CREATE DATABASE IF NOT EXISTS $(DB_NAME); \
		CREATE USER IF NOT EXISTS '$(DB_USER)'@'%' IDENTIFIED BY '$(DB_PASSWORD)'; \
		GRANT ALL PRIVILEGES ON $(DB_NAME).* TO '$(DB_USER)'@'%'; \
		FLUSH PRIVILEGES;" || (echo "$(RED)❌ Failed to create database and user$(NC)" && exit 1)
	@echo "$(GREEN)✅ Database and user created$(NC)"
	@$(MAKE) db-migrate
	@$(MAKE) db-verify

db-migrate: ## Run database migrations (create tables)
	@echo "$(YELLOW)📊 Creating database tables...$(NC)"
	@$(DOCKER_COMPOSE) exec -T mysql_db mysql -u $(DB_USER) -p$(DB_PASSWORD) $(DB_NAME) < mysql/deploy/0001-create-all-tables.sql || (echo "$(RED)❌ Failed to create tables$(NC)" && exit 1)
	@echo "$(GREEN)✅ Database tables created successfully$(NC)"

db-verify: ## Verify database setup
	@echo "$(YELLOW)🔍 Verifying database setup...$(NC)"
	@$(DOCKER_COMPOSE) exec mysql_db mysql -u $(DB_USER) -p$(DB_PASSWORD) $(DB_NAME) -e "\
		SELECT 'Database verification:' as Status; \
		SHOW TABLES;" || (echo "$(RED)❌ Database verification failed$(NC)" && exit 1)
	@echo "$(GREEN)✅ Database verification completed$(NC)"

db-reset: ## Reset database (drop and recreate)
	@echo "$(YELLOW)🔄 Resetting database...$(NC)"
	@$(DOCKER_COMPOSE) exec mysql_db mysql -u root -p$(ROOT_PASSWORD) -e "\
		DROP DATABASE IF EXISTS $(DB_NAME); \
		CREATE DATABASE $(DB_NAME);" || (echo "$(RED)❌ Failed to reset database$(NC)" && exit 1)
	@echo "$(GREEN)✅ Database reset completed$(NC)"
	@$(MAKE) db-migrate

db-shell: ## Open MySQL shell
	@echo "$(YELLOW)🐚 Opening MySQL shell...$(NC)"
	@$(DOCKER_COMPOSE) exec mysql_db mysql -u $(DB_USER) -p$(DB_PASSWORD) $(DB_NAME)

build: ## Build the Go application
	@echo "$(YELLOW)🔨 Building application...$(NC)"
	@go build -o bin/$(APP_NAME) . || (echo "$(RED)❌ Build failed$(NC)" && exit 1)
	@echo "$(GREEN)✅ Application built successfully$(NC)"

run: ## Run the application
	@echo "$(YELLOW)🚀 Starting $(APP_NAME)...$(NC)"
	@if [ ! -f bin/$(APP_NAME) ]; then \
		echo "$(YELLOW)Application not built. Building first...$(NC)"; \
		$(MAKE) build; \
	fi
	@echo "$(GREEN)✅ Starting application on port $(APP_PORT)$(NC)"
	@./bin/$(APP_NAME)

test: ## Run all tests
	@echo "$(YELLOW)🧪 Running tests...$(NC)"
	@go test ./... -v || (echo "$(RED)❌ Tests failed$(NC)" && exit 1)
	@echo "$(GREEN)✅ All tests passed$(NC)"

test-coverage: ## Run tests with coverage
	@echo "$(YELLOW)🧪 Running tests with coverage...$(NC)"
	@go test ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✅ Coverage report generated: coverage.html$(NC)"

lint: ## Run linter
	@echo "$(YELLOW)🔍 Running linter...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "$(YELLOW)⚠️  golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest$(NC)"; \
		go vet ./...; \
	fi

format: ## Format Go code
	@echo "$(YELLOW)📝 Formatting code...$(NC)"
	@go fmt ./...
	@goimports -w .
	@echo "$(GREEN)✅ Code formatted$(NC)"

clean: ## Clean up build artifacts and stop containers
	@echo "$(YELLOW)🧹 Cleaning up...$(NC)"
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@$(MAKE) docker-down
	@docker system prune -f
	@echo "$(GREEN)✅ Cleanup completed$(NC)"

status: ## Show status of all services
	@echo "$(YELLOW)📊 Service Status:$(NC)"
	@echo ""
	@echo "$(BLUE)Docker Containers:$(NC)"
	@$(DOCKER_COMPOSE) ps || echo "$(RED)Docker containers not running$(NC)"
	@echo ""
	@echo "$(BLUE)Database Connection:$(NC)"
	@$(DOCKER_COMPOSE) exec mysql_db mysqladmin ping -h localhost --silent && echo "$(GREEN)✅ MySQL is running$(NC)" || echo "$(RED)❌ MySQL is not accessible$(NC)"
	@echo ""
	@echo "$(BLUE)Application:$(NC)"
	@if [ -f bin/$(APP_NAME) ]; then \
		echo "$(GREEN)✅ Application binary exists$(NC)"; \
	else \
		echo "$(YELLOW)⚠️  Application not built$(NC)"; \
	fi

logs: ## Show application logs (if running in background)
	@echo "$(YELLOW)📋 Application logs:$(NC)"
	@if pgrep -f $(APP_NAME) > /dev/null; then \
		echo "$(GREEN)Application is running$(NC)"; \
	else \
		echo "$(YELLOW)Application is not running$(NC)"; \
	fi

api-test: ## Test API endpoints (requires running application)
	@echo "$(YELLOW)🔌 Testing API endpoints...$(NC)"
	@echo "$(BLUE)Testing health endpoint...$(NC)"
	@curl -s http://localhost:$(APP_PORT)/health || echo "$(RED)❌ Health endpoint failed$(NC)"
	@echo ""
	@echo "$(GREEN)✅ API test completed$(NC)"

# Development workflow targets
dev: ## Start development environment
	@echo "$(YELLOW)🚀 Starting development environment...$(NC)"
	@$(MAKE) docker-up
	@$(MAKE) build
	@$(MAKE) run

restart: ## Restart the application
	@echo "$(YELLOW)🔄 Restarting application...$(NC)"
	@pkill -f $(APP_NAME) || true
	@$(MAKE) build
	@$(MAKE) run

# Database utilities
db-backup: ## Backup database
	@echo "$(YELLOW)💾 Creating database backup...$(NC)"
	@mkdir -p backups
	@$(DOCKER_COMPOSE) exec mysql_db mysqldump -u $(DB_USER) -p$(DB_PASSWORD) $(DB_NAME) > backups/$(DB_NAME)_$$(date +%Y%m%d_%H%M%S).sql
	@echo "$(GREEN)✅ Database backup created$(NC)"

db-restore: ## Restore database from backup (usage: make db-restore BACKUP_FILE=path/to/backup.sql)
	@if [ -z "$(BACKUP_FILE)" ]; then \
		echo "$(RED)❌ Please specify BACKUP_FILE. Usage: make db-restore BACKUP_FILE=path/to/backup.sql$(NC)"; \
		exit 1; \
	fi
	@echo "$(YELLOW)📥 Restoring database from $(BACKUP_FILE)...$(NC)"
	@$(DOCKER_COMPOSE) exec -T mysql_db mysql -u $(DB_USER) -p$(DB_PASSWORD) $(DB_NAME) < $(BACKUP_FILE)
	@echo "$(GREEN)✅ Database restored$(NC)"

# Quick commands
quick-setup: dev-setup docker-up db-setup ## Quick setup for development
	@echo "$(GREEN)🎉 Quick setup completed! Run 'make run' to start the application.$(NC)"

full-reset: clean setup ## Complete reset and setup
	@echo "$(GREEN)🎉 Full reset completed!$(NC)"
# This Makefile provisions all necessary resources and sets up the database