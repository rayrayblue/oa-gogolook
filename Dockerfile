ARG GOLANG_VERSION=1.17.3
ARG ALPINE_VERSION=3.13

# Build stage
FROM golang:${GOLANG_VERSION}-alpine${ALPINE_VERSION} AS builder
WORKDIR /app
COPY . .

RUN go build -o main cmd/main.go

# Test stage

FROM builder AS tester
RUN apk --no-cache add make
RUN apk --no-cache add gcc
RUN apk --no-cache add build-base

RUN go test -v -count=1 -cover ./internal/...

# Run stage
FROM alpine:${ALPINE_VERSION} AS runner
WORKDIR /app
COPY --from=builder /app/main .
COPY configs ./configs

EXPOSE 8888
ENTRYPOINT [ "./main" ]
