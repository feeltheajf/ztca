APP = ztca

GIN_PORT ?= 5000
APP_PORT ?= 3000

.PHONY: dep
dep:
	go mod tidy && go mod vendor

.PHONY: dev
dev:
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

.PHONY: test
test: gosec

.PHONY: gosec
gosec: 
	gosec ./...