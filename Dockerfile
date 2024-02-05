FROM golang:1.21 AS builder

ENV GO111MODULE=on \
  CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64

WORKDIR /src
COPY . .

RUN go get main
RUN go build \
  -ldflags "-s -w -extldflags '-static'" \
  -o /bin/app \
  .

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /bin/app /app

USER nobody
ENTRYPOINT ["/app"]