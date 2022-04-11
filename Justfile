fmt:  ## Format code
  treefmt

lint:  ## Lint code
  golangci-lint run

test:  ## Run tests
  go test ./...

build:  ## Build
  nix build
