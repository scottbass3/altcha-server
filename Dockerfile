FROM golang:1.23.1 AS build

RUN apt-get update && apt-get install -y make && rm -rf /var/lib/apt/lists/*

COPY . /src

WORKDIR /src

ARG TARGETOS=linux
ARG TARGETARCH=amd64

RUN go mod download && \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -o /app/altcha ./cmd/altcha

FROM busybox

COPY --from=build /app/altcha /app/altcha
RUN chown -R 1000:1000 /app

WORKDIR /app

CMD ["/app/altcha", "run"]
