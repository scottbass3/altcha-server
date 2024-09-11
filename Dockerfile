FROM reg.cadoles.com/proxy_cache/library/golang:1.23.1 AS build

RUN apt-get update && apt-get install make

COPY . /src

WORKDIR /src

RUN go mod download && make GORELEASER_ARGS="build --rm-dist --single-target --snapshot" goreleaser

FROM reg.cadoles.com/proxy_cache/library/busybox

COPY --from=build /src/dist/altcha_linux_amd64_v1 /app
RUN chown -R 1000:1000 /app

WORKDIR /app

CMD ["bin/altcha", "run"]
