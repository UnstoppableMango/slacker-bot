FIND      ?= find
GO        ?= go
GOMOD2NIX ?= gomod2nix
PODMAN    ?= podman

GO_SRC != $(FIND) . -path '*.go' -printf '%P\n'

build: bin/slacker-bot
docker: bin/image.tar.gz

load: bin/stream-image.sh
	${CURDIR}/$< | $(PODMAN) load

up: load
	$(PODMAN) compose up

format fmt:
	nix fmt

tidy: go.sum gomod2nix.toml

update:
	nix flake update
	$(GO) get -u ./...

bin/stream-image.sh: ${GO_SRC}
	nix build .#ctr --out-link $@

bin/image.tar.gz: bin/stream-image.sh
	${CURDIR}/$< >$@

bin/slacker-bot: result
	mkdir -p ${@D} && ln -s $(abspath $<)/bin/slacker-bot $@

gomod2nix.toml: go.sum
	$(GOMOD2NIX) generate

go.sum: go.mod ${GO_SRC}
	$(GO) mod tidy

result: ${GO_SRC}
	nix build .#slacker-bot
