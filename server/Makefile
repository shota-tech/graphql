.PHONY: install-tools
install-tools:
	cat tools.go | awk -F'"' '/_/ {print $$2}' | xargs -tI {} go install {}

.PHONY: codegen
codegen:
	go generate ./...
