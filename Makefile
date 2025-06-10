VERSION := 1.3.3


# make bump-version VERSION=1.0.9

.PHONY: bump-version tidy up-version
bump-version:
	@echo "Bumping version to $(VERSION)"
	@OS=$$(uname); \
	if [ "$$OS" = "Darwin" ]; then \
		find . -name "go.mod" -type f -exec sed -i '' -E 's|(github.com/gone-io/goner/[^ ]*)[ \t]*v[0-9]+\.[0-9]+\.[0-9]+|\1 v$(VERSION)|g' {} \; ; \
	else \
		find . -name "go.mod" -type f -exec sed -i -E 's|(github.com/gone-io/goner/[^ ]*)[ \t]*v[0-9]+\.[0-9]+\.[0-9]+|\1 v$(VERSION)|g' {} \; ; \
	fi
	@echo "Version bump complete"


tidy:
	@echo "Running go mod tidy for all modules..."
	@find . -name go.mod  | xargs -n1 dirname | xargs -L1 bash -c 'cd "$$0" && echo "Processing directory: $$0" && go mod tidy'
	@echo "Tidy complete"


up:
	make bump-version
	make tidy

gone-up:
	@find . -name go.mod  | xargs -n1 dirname | xargs -L1 bash -c 'cd "$$0" && echo "Processing directory: $$0" && go get -u github.com/gone-io/gone/v2'
	make tidy