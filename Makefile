all: build start

build:
	@echo "🔨 Building..."
	@go build -o ./bin/benchmq main.go
	@echo "✅ Built Successful."

start:
	@echo "🚀 Starting..."
	@if [ ! -f ./bin/benchmq ]; then \
		echo "👀 Benchmq not found, building..."; \
		$(MAKE) build; \
	fi
	@./bin/benchmq
	@echo "✅ Started Successfully."

run:
	@echo "🚀 Running..."
	@go run main.go
	@echo "✅ Ran Successful."

clean:
	@echo "🧹 Cleaning..."
	@rm -rf ./bin/*
	@echo "✅ Cleaned Successful."

conn:
	@echo "🔌 Connecting..."
	@if [ ! -f ./bin/benchmq ]; then \
		echo "👀 Benchmq not found, building..."; \
		$(MAKE) build; \
	fi
	@./bin/benchmq conn
	@echo "✅ Connected Successful."

.PHONY: all build start run clean conn
