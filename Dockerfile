FROM golang:1.22.5-alpine3.20 as builder
RUN apk update && apk upgrade && apk add --no-cache ca-certificates
RUN update-ca-certificates
WORKDIR /app
COPY ./app .
RUN go mod init aigiscorp.dev/redata-api-consumer
RUN go get
RUN go generate ./...
RUN CGO_ENABLED=0 GOOS=linux go build -o redata-api-consumer

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/redata-api-consumer /app/
EXPOSE 8080
CMD ["/app/redata-api-consumer"]
