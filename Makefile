APP = ztca
COV = coverage.out
TAG = v$(shell cat VERSION)

GIN_PORT ?= 5000
APP_PORT ?= 3000

.PHONY: dep
dep:
	go mod tidy && go mod vendor

.PHONY: dev
dev: dep
	gin \
		--all \
		--immediate \
		--bin $(APP) \
		--excludeDir $(APP).sqlite3 \
		--excludeDir testdata \
		--excludeDir vendor \
		--excludeDir web \
		--port $(GIN_PORT) \
		--appPort $(APP_PORT) \
		run \
		--

.PHONY: release
release:
	git tag -a $(TAG) -m "$(TAG) release"
	git push origin $(TAG)

.PHONY: test
test: unittest gosec trufflehog

.PHONY: unittest
unittest:
	go test -v -race -coverprofile=$(COV) ./... \
		&& go tool cover -func $(COV)

.PHONY: gosec
gosec:
	gosec ./...

.PHONY: trufflehog
trufflehog:
	trufflehog3
