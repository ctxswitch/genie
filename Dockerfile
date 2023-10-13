FROM golang:1.21-alpine AS builder
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -v -o genie main.go

FROM alpine:latest
COPY --from=builder /usr/src/app/genie.d /etc/genie.d
COPY --from=builder /usr/src/app/genie /genie
RUN : \
  && apk update \
  && apk add ca-certificates \
  && rm -rf /var/cache/apk/* \
  && :
RUN update-ca-certificates
CMD ["/genie", "generate", "--config", "/etc/genie.d"]
