FROM golang:1.25.2 AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS
ARG TARGETARCH
ENV CGO_ENABLED=0

RUN GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} \
    go build -buildvcs=false -trimpath -ldflags="-s -w" -o /out/server ./cmd/server
RUN mkdir -p /out/data

FROM gcr.io/distroless/static:nonroot

WORKDIR /app

COPY --from=builder --chown=65532:65532 /out/server /app/server
COPY --from=builder --chown=65532:65532 /out/data /app/data

EXPOSE 8080
USER 65532:65532

ENTRYPOINT ["/app/server"]
