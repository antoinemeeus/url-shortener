all:: help

help ::
	@grep -E '^[a-zA-Z_-]+\s*:.*?## .*$$' ${MAKEFILE_LIST} | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[32m%-30s\033[0m %s\n", $$1, $$2}'

start :: ## Run service
	@echo "  > Starting service"
	@docker-compose up -d
	@go run cmd/server/main.go
