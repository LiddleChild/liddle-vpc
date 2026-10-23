.PHONY: build
build:
	@echo "Building expenses"
	@cd services/expenses/scripts && ./build.sh

	@echo "Building backend"
	@cd services/backend/scripts && ./build.sh
