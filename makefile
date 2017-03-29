MAKEFLAGS+=--ignore-errors
MAKEFLAGS+=--no-print-directory
SHELL:=/bin/bash

.PHONY: clean
.PHONY: clean_css
.PHONY: clean_go
.PHONY: clean_js
.PHONY: lint
.PHONY: lint_css
.PHONY: lint_go
.PHONY: lint_js
.PHONY: build
.PHONY: build_css
.PHONY: build_js
.PHONY: build_go
.PHONY: bootstrap
.PHONY: populate
.PHONY: serve
.PHONY: watch
.PHONY: watch_css
.PHONY: watch_js
.PHONY: watch_go

clean:
	@$(MAKE) --jobs=3 clean_css clean_go clean_js

clean_css:
	@/usr/bin/rm -rf ./assets/*.css

clean_go:
	@/usr/bin/rm -rf ./torrents

clean_js:
	@/usr/bin/rm -rf ./assets/*.js

lint:
	@$(MAKE) --jobs=3 lint_css lint_go lint_js

lint_css:
	@./node_modules/.bin/csslint ./resources/css/**

lint_go:
	@/usr/bin/go fmt
	@${GOPATH}/bin/golint go/actions
	@${GOPATH}/bin/golint go/routes
	@${GOPATH}/bin/golint go/settings
	@${GOPATH}/bin/golint go/views
	@${GOPATH}/bin/golint .
	@/usr/bin/go vet github.com/mahendrakalkura/torrents/go/actions
	@/usr/bin/go vet github.com/mahendrakalkura/torrents/go/routes
	@/usr/bin/go vet github.com/mahendrakalkura/torrents/go/settings
	@/usr/bin/go vet github.com/mahendrakalkura/torrents/go/views
	@/usr/bin/go vet github.com/mahendrakalkura/torrents

lint_js:
	@./node_modules/.bin/eslint ./resources/js/**

build:
	@$(MAKE) --jobs=3 build_css build_go build_js

build_css:
	@./node_modules/.bin/concat-cli                          \
		--files                                              \
			./node_modules/bootstrap/dist/css/bootstrap.css  \
			./node_modules/font-awesome/css/font-awesome.css \
			./resources/css/all.css                          \
		--output ./assets/compressed.css > /dev/null
	@/usr/bin/sed --expression='s/..\/fonts\//fonts\//g' --in-place ./assets/compressed.css
	@./node_modules/.bin/cleancss --output=./assets/compressed.min.css ./assets/compressed.css

build_go:
	@/usr/bin/go build

build_js:
	@./node_modules/.bin/concat-cli                       \
		--files                                           \
			./node_modules/jquery/dist/jquery.js          \
			./node_modules/bootstrap/dist/js/bootstrap.js \
			./resources/js/all.js                         \
		--output ./assets/compressed.js > /dev/null
	@./node_modules/.bin/uglifyjs --output=./assets/compressed.min.js ./assets/compressed.js

watch:
	@$(MAKE) --jobs=3 watch_css watch_go watch_js

watch_css:
	@${GOPATH}/bin/reflex --regex=resources/css/.* $(MAKE) build_css

watch_go:
	@${GOPATH}/bin/reflex --regex=.*\.\(go\|html\) --start-service $(MAKE) serve

watch_js:
	@${GOPATH}/bin/reflex --regex=resources/js/.* $(MAKE) build_js

serve:
	@/usr/bin/pkill ./torrents || true
	@/usr/bin/go build
	@./torrents --action=serve
