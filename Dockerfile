FROM golang:1.24-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/pqcweb ./cmd/pqcweb

FROM alpine:3.20

RUN apk add --no-cache git ca-certificates && adduser -D -H -u 10001 pqc
WORKDIR /app
COPY --from=build /out/pqcweb /app/pqcweb
COPY rules ./rules
COPY sample_repo ./sample_repo
USER pqc
EXPOSE 8080
ENV PORT=8080
ENTRYPOINT ["/app/pqcweb"]
