FROM golang:1.21 AS builder

ENV GO111MODULE=on \
  CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64

RUN apt-get -qq update && \
  apt-get -yqq install upx

WORKDIR /src
COPY . .

RUN go get main
RUN go build \
  -ldflags "-s -w -extldflags '-static'" \
  -o /bin/app \
  . \
  && strip /bin/app \
  && upx -q -9 /bin/app

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /bin/app /app

USER nobody
ENTRYPOINT ["/app"]