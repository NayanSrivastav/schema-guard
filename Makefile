.PHONY: test test-go test-ui test-python

# The primary orchestration target running all suites globally seamlessly natively
test: test-go test-ui test-python
	@echo "============================================================"
	@echo "✅ All platform tests across Go, React, and Python passed!"
	@echo "============================================================"

# Executes all Golang architectural logic and BDD server endpoints securely inside an isolated Alpine container natively
test-go:
	@echo "------------------------------------------------------------"
	@echo "🧪 Running Go Backend Core & API Tests (Dockerized)..."
	@echo "------------------------------------------------------------"
	@docker run --rm -v $(PWD):/app -w /app golang:alpine sh -c "go mod tidy && go test -v ./..."

# Executes Frontend component verification mapping strictly natively inside NodeJS Docker environments smoothly
test-ui:
	@echo "------------------------------------------------------------"
	@echo "🧪 Running React/Vite Dashboard UI Tests (Dockerized)..."
	@echo "------------------------------------------------------------"
	@docker run --rm -v $(PWD)/dashboard:/app -w /app node:20-alpine sh -c "npm install && npm install --no-save vitest @testing-library/react @testing-library/jest-dom jsdom && npx vitest run --environment jsdom"

# Executes Python abstraction hooks cleanly natively
test-python:
	@echo "------------------------------------------------------------"
	@echo "🧪 Running Python SDK BDD Verification Maps (Dockerized)..."
	@echo "------------------------------------------------------------"
	@docker run --rm -v $(PWD)/sdks/python:/app -w /app python:3.9-slim sh -c "pip install pytest jsonschema requests && PYTHONPATH=. pytest -v tests/"
