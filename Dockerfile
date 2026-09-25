FROM --platform=$BUILDPLATFORM golang:1.26 AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS=linux
ARG TARGETARCH=amd64
ARG VERSION=dev
ARG GIT_REF=unknown
ARG BUILD_DATE=unknown

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build \
    -ldflags "-s -w -X 'main.ProjectVersion=${VERSION}' -X 'main.GitRef=${GIT_REF}' -X 'main.BuildDate=${BUILD_DATE}'" \
    -o /app/altcha ./cmd/altcha

FROM busybox

COPY --from=build /app/altcha /app/altcha
RUN chown -R 1000:1000 /app

WORKDIR /app

CMD ["/app/altcha", "run"]
