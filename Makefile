.PHONY: build
build:
	@echo "Building expenses"
	@cd services/expenses/scripts && ./build.sh

	@echo "Building nexus"
	@cd services/nexus/scripts && ./build.sh
