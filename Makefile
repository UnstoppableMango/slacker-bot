FIND      ?= find
GO        ?= go
GOMOD2NIX ?= gomod2nix

GO_SRC != $(FIND) . -path '*.go' -printf '%P\n'

build:
	nix build .#slacker-bot

format fmt:
	nix fmt

tidy: go.sum gomod2nix.toml

gomod2nix.toml: go.sum
	$(GOMOD2NIX) generate

go.sum: go.mod ${GO_SRC}
	$(GO) mod tidy
