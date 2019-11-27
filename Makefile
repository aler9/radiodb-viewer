
GO_BASE_IMAGE = amd64/golang:1.12-alpine3.10

help:
	@echo "usage: make [action]"
	@echo ""
	@echo "available actions:"
	@echo ""
	@echo "  mod-tidy"
	@echo "  format"
	@echo "  dev"
	@echo "  prod"
	@echo ""

mod-tidy:
	echo "FROM $(GO_BASE_IMAGE) \n\
	RUN apk add git \n\
	" | docker build - -t rdbviewer-modtidy
	docker run --rm -it -v $(PWD):/s rdbviewer-modtidy \
	sh -c "cd /s/back && go get -m ./... && go mod tidy"

format:
	docker run --rm -it -v $(PWD):/s $(GO_BASE_IMAGE) \
	sh -c "cd /s && find . -type f -name '*.go' | xargs gofmt -l -w -s"

BUILD = docker build . \
    --build-arg BUILD_MODE=$(1) \
    -t radiodb-viewer-$(1)

dev:
	docker run --rm -it -v radiodb:/out amd64/alpine:3.8 \
	sh -c "apk add curl && curl --compressed -o/out/radiodb.json https://radiodb.freeddns.org/dumpget"
	make dev-inner

dev-inner:
	$(if $(shell command -v inotifywait 2>&1),,$(error inotifywait is required))

	docker kill radiodb-viewer-dev >/dev/null 2>&1 || true

	$(call BUILD,dev) \
	&& ( docker run --rm \
	--read-only \
	-v radiodb:/data \
	-p 7446:7446 \
	--name radiodb-viewer-dev \
	radiodb-viewer-dev & ) || true

	inotifywait -qre close_write ./*/
	make dev-inner

prod:
	$(call BUILD,prod)
