###################################
FROM amd64/golang:1.12-alpine3.10 AS countrymeta

RUN apk add --no-cache \
    git

WORKDIR /s

RUN go mod init temp \
    && go get github.com/pariz/gountries@adb00f6

COPY back/gencountrymeta.go ./
RUN go run gencountrymeta.go

###################################
FROM amd64/golang:1.12-alpine3.10 AS back

RUN apk add --no-cache \
    git \
    protobuf-dev

WORKDIR /s

COPY back/go.mod back/go.sum ./
RUN go mod download

# grpc
RUN go install github.com/golang/protobuf/protoc-gen-go
COPY back/defs/*.proto ./defs/
COPY back/shared/*.proto ./shared/
RUN cd ./defs && protoc pdefs.proto \
    --go_out=paths=source_relative:.
RUN cd ./shared && protoc -I ../defs -I . db.proto \
    --go_out=paths=source_relative,plugins=grpc,Mpdefs.proto=rdbviewer/defs:.

ARG BUILD_MODE
RUN test -n "$BUILD_MODE"

COPY --from=countrymeta /s/countrymeta.go ./shared/
COPY back/defs/*.go ./defs/
COPY back/shared/*.go ./shared/

ENV CGO_ENABLED 0

COPY back/db/*.go ./db/
RUN go build -o /build/db ./db

COPY back/router/*.go ./router/
RUN go build -ldflags "-X main.BUILD_MODE=$BUILD_MODE" -o /build/router ./router

###################################
FROM amd64/golang:1.12-alpine3.10 AS countryflag

RUN apk add --no-cache \
    git

WORKDIR /s

RUN go mod init temp \
    && go get github.com/disintegration/imaging@v1.6.0

COPY image/gencountryflag.go ./
RUN git clone https://github.com/hjnilsson/country-flags \
    && cd country-flags*/ \
    && git checkout d5d1cc4 \
    && go run ../gencountryflag.go \
    && mkdir -p /build/static/countryflag \
    && mv png1000px/* /build/static/countryflag/

###################################
FROM amd64/node:12-alpine AS image

WORKDIR /s

COPY image/package*.json ./
RUN npm i

COPY image/*.png /build/static/

COPY image/favicons.js image/favicon.svg ./
COPY template/frame.tpl /build/template/
RUN node favicons.js

COPY image/*.svg ./
RUN find . -maxdepth 1 -name '*.svg' ! -name 'favicon.svg' \
    | xargs -n1 sh -c 'node_modules/.bin/svgo -i $0 -o /build/static/$(basename $0)'

###################################
FROM amd64/node:12-alpine AS script

WORKDIR /s

COPY script/package*.json ./
RUN npm i

ARG BUILD_MODE
RUN test -n "$BUILD_MODE"
COPY .eslintrc.browser.js ./.eslintrc.js
COPY script/babel.config.js \
    script/webpack.config.js \
    .browserslistrc script/*.jsx ./
RUN node_modules/.bin/webpack

###################################
FROM amd64/alpine:3.10 AS template

COPY template/*.tpl /build/template/
COPY --from=image /build/template/frame.tpl /build/template/

RUN RND=$(cat /dev/urandom | tr -dc 'a-z0-9' | fold -w 8 | head -n 1) \
    && sed -i "s/\(style\.css\)/\1?$RND/" /build/template/frame.tpl \
    && sed -i "s/\(script\.js\)/\1?$RND/" /build/template/frame.tpl \
    && sed -i "s/\(\.ico\)/\1?$RND/" /build/template/* \
    && sed -i "s/\(\.svg\)/\1?$RND/" /build/template/* \
    && sed -i "s/\(\.png\)/\1?$RND/" /build/template/*

###################################
FROM amd64/node:12-alpine AS style

WORKDIR /s

COPY style/package*.json ./
RUN npm i

COPY stylelint.config.js \
    style/postcss.config.js \
    .browserslistrc \
    style/*.scss ./
RUN mkdir -p /build/static \
    && node_modules/.bin/stylelint *.scss \
    && cat style.scss \
    | node_modules/.bin/node-sass \
    | node_modules/.bin/postcss \
    > /build/static/style.css

COPY style/googlefonts.js ./
RUN node googlefonts.js /build/static/style.css

###################################
FROM amd64/alpine:3.10

RUN echo $'#!/bin/sh\n\
/build/db &\n\
/build/router &\n\
wait\n\
' > /start.sh && chmod +x /start.sh

RUN adduser -D -H -s /bin/sh -u 1000 user

COPY --from=back /build /build
COPY --from=countryflag /build /build
COPY --from=image /build /build
COPY --from=script /build /build
COPY --from=template /build /build
COPY --from=style /build /build

USER user

ENTRYPOINT [ "/start.sh" ]
