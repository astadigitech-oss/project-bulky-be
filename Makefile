.PHONY: audit

# Run dependency integrity, vulnerability, and static-analysis checks.
audit:
	go mod verify
	govulncheck ./...
	go vet ./...
