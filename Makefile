.PHONY: build
build:
	@echo "Building expenses"
	@cd services/expenses/scripts && ./build.sh

	@echo "Building gristguard"
	@cd services/gristguard/scripts && ./build.sh
