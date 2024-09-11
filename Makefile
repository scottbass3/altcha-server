IMAGE_REPO ?= reg.cadoles.com/cadoles/altcha

GORELEASER_VERSION ?= v1.13.1
GORELEASER_ARGS ?= release --snapshot --rm-dist

ALTCHA_VERSION ?= 
GIT_COMMIT := $(shell git rev-parse --short HEAD)
DATE_VERSION := $(shell date +%Y.%-m.%-d)
FULL_VERSION := v$(DATE_VERSION)-$(GIT_COMMIT)$(if $(shell git diff --stat),-dirty,)

.PHONY: test
test: test-go ## Executing tests

test-go:
	go test -v -race ./...

build: build-altcha

build-altcha: ## Build executable
	CGO_ENABLED=0 go build \
		-v \
		-ldflags "\
			-X 'main.GitRef=$(shell git rev-parse --short HEAD)' \
			-X 'main.ProjectVersion=$(FULL_VERSION)' \
			-X 'main.BuildDate=$(shell date --utc --rfc-3339=seconds)' \
		" \
		-o ./bin/altcha \
		./cmd/altcha

run:
	bin/altcha $(ATLCHA_CMD)

.PHONY: goreleaser
goreleaser:
	curl -sfL https://goreleaser.com/static/run | VERSION=$(GORELEASER_VERSION) GORELEASER_CURRENT_TAG="$(FULL_VERSION)" bash /dev/stdin $(GORELEASER_ARGS)

build-image:
	docker build \
		-t "$(IMAGE_REPO):latest" \
		.

release-image: .mktools
	@[ ! -z "$(MKT_PROJECT_VERSION)" ] || ( echo "Just downloaded mktools. Please re-run command."; exit 1 )
	docker tag "$(IMAGE_REPO):latest" "$(IMAGE_REPO):$(MKT_PROJECT_VERSION)"
	docker tag "$(IMAGE_REPO):latest" "$(IMAGE_REPO):$(MKT_PROJECT_SHORT_VERSION)"
	docker push "$(IMAGE_REPO):latest"
	docker push "$(IMAGE_REPO):$(MKT_PROJECT_VERSION)"
	docker push "$(IMAGE_REPO):$(MKT_PROJECT_SHORT_VERSION)"

.PHONY: mktools
mktools:
	rm -rf .mktools
	curl -q https://forge.cadoles.com/Cadoles/mktools/raw/branch/master/install.sh | $(SHELL)

.mktools:
	$(MAKE) mktools

-include .mktools/*.mk