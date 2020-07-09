##################################
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

COPY back/frame.html /build/

##################################
FROM amd64/golang:1.12-alpine3.10 AS countryflags

RUN apk add --no-cache \
    git

WORKDIR /s

RUN go mod init temp \
    && go get github.com/disintegration/imaging@v1.6.0

COPY front/images/gencountryflag.go ./
RUN git clone https://github.com/hjnilsson/country-flags \
    && cd country-flags*/ \
    && git checkout d5d1cc4 \
    && go run ../gencountryflag.go \
    && mkdir /countryflags \
    && mv png1000px/* /countryflags/

###################################
FROM amd64/node:12-alpine AS front

WORKDIR /s

COPY front/package.json front/yarn.lock ./
RUN yarn install

COPY --from=back /build/frame.html /build/
COPY front/images/favicon.svg ./
COPY front/favicons.js ./
RUN node favicons.js

COPY --from=countryflags /countryflags ./images/countryflags
COPY front ./

ARG BUILD_MODE
RUN node_modules/.bin/webpack

RUN sed -i "s/script\.js/$(ls /build/static/script* | xargs basename)/" /build/frame.html \
    && sed -i "s/style\.css/$(ls /build/static/style* | xargs basename)/" /build/frame.html

###################################
FROM amd64/alpine:3.10

RUN adduser -D -H -s /bin/sh -u 1078 user

COPY --from=back /build /build
COPY --from=front /build /build

COPY start.sh /
RUN chmod +x /start.sh

ENTRYPOINT [ "/start.sh" ]
