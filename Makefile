
ALPINE_BASE_IMAGE = amd64/alpine:3.12
GO_BASE_IMAGE = amd64/golang:1.14-alpine3.12

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
	sh -c "cd /s/back && go get ./... && go mod tidy"

format:
	docker run --rm -it -v $(PWD):/s $(GO_BASE_IMAGE) \
	sh -c "cd /s && find . -type f -name '*.go' | xargs gofmt -l -w -s"

BUILD = docker build . \
    --build-arg BUILD_MODE=$(1) \
    -t radiodb-viewer-$@

dev:
	$(if $(shell which inotifywait),,$(error inotify-tools non trovato. installare con apt install -y inotify-tools))

	docker run --rm -it -v radiodb:/out $(ALPINE_BASE_IMAGE) \
	sh -c "apk add curl && curl --compressed -o/out/radiodb.json https://radiodb.freeddns.org/dumpget"

	while true; do \
		$(call BUILD,development) \
		&& ( docker kill radiodb-viewer-dev >/dev/null 2>&1 || true ) \
		&& ( docker run --rm \
		--read-only \
		-v radiodb:/data \
		-p 7446:7446 \
		--name radiodb-viewer-dev \
		radiodb-viewer-dev & ); \
		inotifywait -qre close_write ./*/; \
	done
