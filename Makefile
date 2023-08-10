
# help
.PHONY: help
help:
	@echo 'make fmt|test|upgrade-deps'

# fmt
.PHONY: fmt
fmt:
	@bash scripts/fmt.sh

# test
.PHONY: test
test:
	@bash scripts/test.sh

# upgrade_deps
.PHONY: upgrade-deps
upgrade-deps:
	@go get -d -u -t ./...


