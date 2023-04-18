EXE  := tfselect
PKG  := github.com/basti0nz/tfselect
VER := 0.1.20
PATH := build:$(PATH)
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

ifneq (,$(wildcard ./version))
    include version
    export
endif


$(EXE): go.mod *.go lib/*.go
	go build -v -ldflags "-X main.version=$(VER)" -o ./dist/$@ $(PKG)

.PHONY: release
release: $(EXE) clean gorelease alpine 

.PHONY: darwin linux
darwin linux:
	GOOS=$@ go build -ldflags "-X main.version=$(VER)" -o ./dist/$(EXE)-$(VER)-$@-$(GOARCH) $(PKG)

.PHONY: clean
clean:
	rm -rf ./dist/

.PHONY: snap
snap:
	(mkdir ./dist && multipass stop snapcraft-tfswitch && multipass delete snapcraft-tfswitch && multipass purge && rm -f ./dist/tfselect_*.snap) || true  && snapcraft && mv tfselect*.snap ./dist

.PHONY: snap-stop
snap-stop:
	(multipass stop snapcraft-tfselect && multipass delete snapcraft-tfselect && multipass purge ) || true
.PHONY: alpine
alpine:
	cd ./alpine && bash ./build.sh

.PHONY: gorelease
gorelease:
	rm -rf ./dist/
	goreleaser

.PHONY: tag
tag:
	git tag -a $(VER) -m "New release"


