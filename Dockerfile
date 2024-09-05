FROM golang:1.22-alpine AS builder

ARG VERSION=dev

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN CGO_ENABLED=0 go build -trimpath --ldflags "-s -w -X ctx.sh/genie/pkg/build.Version=${VERSION}" -o ./dist/genie main.go

FROM alpine:latest
COPY --from=builder /usr/src/app/examples /etc/genie.d
COPY --from=builder /usr/src/app/dist/genie /bin/genie

RUN addgroup -S genie && adduser -S genie -G genie

RUN : \
  && apk update \
  && apk add ca-certificates \
  && rm -rf /var/cache/apk/* \
  && :
RUN update-ca-certificates

USER genie
CMD ["/bin/genie", "generate", "--config", "/etc/genie.d"]
