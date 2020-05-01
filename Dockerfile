###################################
FROM amd64/golang:1.12-alpine3.10 AS countrymeta

RUN apk add --no-cache \
    git

WORKDIR /s

RUN go mod init temp \
    && go get github.com/pariz/gountries@adb00f6

COPY backend/gencountrymeta.go ./
RUN go run gencountrymeta.go

###################################
FROM amd64/golang:1.12-alpine3.10 AS backend

RUN apk add --no-cache \
    git \
    protobuf-dev

WORKDIR /s

COPY backend/go.mod backend/go.sum ./
RUN go mod download

# grpc
RUN go install github.com/golang/protobuf/protoc-gen-go
COPY backend/defs/*.proto ./defs/
COPY backend/shared/*.proto ./shared/
RUN cd ./defs && protoc pdefs.proto \
    --go_out=paths=source_relative:.
RUN cd ./shared && protoc -I ../defs -I . db.proto \
    --go_out=paths=source_relative,plugins=grpc,Mpdefs.proto=rdbviewer/defs:.

ARG BUILD_MODE
RUN test -n "$BUILD_MODE"

COPY --from=countrymeta /s/countrymeta.go ./shared/
COPY backend/defs/*.go ./defs/
COPY backend/shared/*.go ./shared/

ENV CGO_ENABLED 0

COPY backend/db/*.go ./db/
RUN go build -o /build/db ./db

COPY backend/router/*.go ./router/
RUN go build -ldflags "-X main.BUILD_MODE=$BUILD_MODE" -o /build/router ./router

###################################
FROM amd64/golang:1.12-alpine3.10 AS countryflag

RUN apk add --no-cache \
    git

WORKDIR /s

RUN go mod init temp \
    && go get github.com/disintegration/imaging@v1.6.0

COPY images/gencountryflag.go ./
RUN git clone https://github.com/hjnilsson/country-flags \
    && cd country-flags*/ \
    && git checkout d5d1cc4 \
    && go run ../gencountryflag.go \
    && mkdir -p /build/static/countryflag \
    && mv png1000px/* /build/static/countryflag/

###################################
FROM scratch AS templates

COPY templates/*.tpl /build/templates/

###################################
FROM amd64/node:12-alpine AS images

WORKDIR /s

COPY images/package.json images/yarn.lock ./
RUN yarn install

COPY images/*.png /build/static/

COPY images/favicons.js images/favicon.svg ./
COPY --from=templates /build/templates/frame.tpl /build/templates/
RUN node favicons.js

COPY images/*.svg ./
RUN find . -maxdepth 1 -name '*.svg' ! -name 'favicon.svg' \
    | xargs -n1 sh -c 'node_modules/.bin/svgo -i $0 -o /build/static/$(basename $0)'

###################################
FROM amd64/node:12-alpine AS scriptstyle

WORKDIR /s

COPY script/package.json script/yarn.lock ./
RUN yarn install

COPY script/webpack.config.js \
    script/.browserslistrc \
    script/babel.config.js \
    script/.eslintrc.js \
    style/postcss.config.js \
    style/stylelint.config.js \
    ./

COPY script/*.jsx ./
COPY style/*.scss ./

ARG BUILD_MODE
RUN node_modules/.bin/webpack

COPY --from=images /build/templates/frame.tpl /build/templates/
RUN sed -i "s/script\.js/$(ls /build/static/script* | xargs basename)/" /build/templates/frame.tpl \
    && sed -i "s/style\.css/$(ls /build/static/style* | xargs basename)/" /build/templates/frame.tpl

###################################
FROM amd64/alpine:3.10

RUN adduser -D -H -s /bin/sh -u 1078 user

COPY --from=backend /build /build
COPY --from=countryflag /build /build
COPY --from=templates /build /build
COPY --from=images /build /build
COPY --from=scriptstyle /build /build

COPY start.sh /
RUN chmod +x /start.sh

ENTRYPOINT [ "/start.sh" ]
