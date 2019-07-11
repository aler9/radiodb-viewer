###################################
FROM amd64/golang:1.12-alpine3.9 AS countrymeta

RUN apk add git

WORKDIR /countrymeta
RUN go mod init temp \
    && go get github.com/pariz/gountries@adb00f6
COPY gencountrymeta.go ./
RUN go run gencountrymeta.go

###################################
FROM amd64/golang:1.12-alpine3.9 AS back

RUN apk add \
    git \
    protobuf-dev

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

# grpc
RUN go install github.com/golang/protobuf/protoc-gen-go
COPY back/defs/*.proto ./back/defs/
COPY back/shared/*.proto ./back/shared/
RUN cd ./back/defs && protoc pdefs.proto \
    --go_out=paths=source_relative:.
RUN cd ./back/shared && protoc -I ../defs -I . db.proto \
    --go_out=paths=source_relative,plugins=grpc,Mpdefs.proto=rdbviewer/back/defs:.

ARG BUILD_MODE
RUN test -n "$BUILD_MODE"

COPY --from=countrymeta /countrymeta/countrymeta.go ./back/shared/
COPY back/defs/*.go ./back/defs/
COPY back/shared/*.go ./back/shared/

ENV CGO_ENABLED 0

COPY back/db/*.go ./back/db/
RUN go build -o /build/db ./back/db

COPY back/router/*.go ./back/router/
RUN go build -ldflags "-X main.BUILD_MODE=$BUILD_MODE" -o /build/router ./back/router

###################################
FROM amd64/golang:1.12-alpine3.9 AS countryflag

RUN apk add git

WORKDIR /countryflag
RUN go mod init temp \
    && go get github.com/disintegration/imaging@v1.6.0
COPY gencountryflag.go ./
RUN git clone https://github.com/hjnilsson/country-flags \
    && cd country-flags*/ \
    && git checkout d5d1cc4 \
    && go run ../gencountryflag.go \
    && mkdir -p /build/static/countryflag \
    && mv png1000px/* /build/static/countryflag/

###################################
FROM amd64/node:10-alpine AS nodebase

# phantomjs --> for favicon generator
# https://github.com/dustinblackman/phantomized/

RUN apk add --no-cache \
    curl \
    && rm -rf /var/cache/apk/*

RUN curl -Ls https://github.com/dustinblackman/phantomized/releases/download/2.1.1a/dockerized-phantomjs.tar.gz | tar xz

WORKDIR /src

COPY package*.json ./
RUN npm i
ENV PATH $PATH:/src/node_modules/.bin

###################################
FROM nodebase AS template

COPY template/*.tpl ./template/
RUN mkdir -p /build/template \
    && cp template/*.tpl /build/template/
COPY genfavicon.js ./
COPY image/favicon.svg ./image/
RUN node genfavicon.js
RUN RND=$(cat /dev/urandom | tr -dc 'a-z0-9' | fold -w 8 | head -n 1) \
    && sed -i "s/\(style\.css\)/\1?$RND/" /build/template/frame.tpl \
    && sed -i "s/\(script\.js\)/\1?$RND/" /build/template/frame.tpl \
    && sed -i "s/\(\.ico\)/\1?$RND/" /build/template/* \
    && sed -i "s/\(\.svg\)/\1?$RND/" /build/template/* \
    && sed -i "s/\(\.png\)/\1?$RND/" /build/template/*

###################################
FROM nodebase AS image

RUN mkdir -p /build/static
COPY image/*.png ./image/
RUN cp image/*.png /build/static/
COPY image/*.svg ./image/
RUN find ./image/ -maxdepth 1 -name '*.svg' ! -name 'favicon.svg' \
    | xargs -n1 sh -c 'svgo -i $0 -o /build/static/$(basename $0)'

###################################
FROM nodebase AS script

ARG BUILD_MODE
RUN test -n "$BUILD_MODE"

COPY .eslintrc.js babel.config.js webpack.config.js ./
COPY script/*.jsx ./script/
RUN webpack

###################################
FROM nodebase AS font

RUN get-google-fonts -o /build/static/fonts -p /static/fonts/ -c temp1.css \
    -i "https://fonts.googleapis.com/css?family=Muli:200,300,400,700"
RUN get-google-fonts -o /build/static/fonts -p /static/fonts/ -c temp2.css \
    -i "https://fonts.googleapis.com/css?family=Cutive+Mono"
RUN cat /build/static/fonts/temp*.css > ./fonts.scss && rm /build/static/fonts/temp*.css

###################################
FROM nodebase AS style

COPY --from=font /src/fonts.scss /src/fonts.scss
COPY stylelint.config.js postcss.config.js ./
COPY style/*.scss ./style/
WORKDIR ./style
RUN mkdir -p /build/static \
    && stylelint *.scss \
    && cat style.scss \
    | node-sass \
    | postcss \
    > /build/static/style.css

###################################
FROM amd64/alpine:3.9 AS run

RUN echo $'#!/bin/sh\n\
/build/db &\n\
/build/router &\n\
wait\n\
' > /start.sh && chmod +x /start.sh

RUN adduser -D -H -s /bin/sh -u 1000 user

COPY --from=back /build /build
COPY --from=template /build /build
COPY --from=countryflag /build /build
COPY --from=image /build /build
COPY --from=script /build /build
COPY --from=font /build /build
COPY --from=style /build /build

USER user

ENTRYPOINT [ "/start.sh" ]
