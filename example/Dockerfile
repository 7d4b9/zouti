FROM golang:1.16.6-alpine3.13 as builder
WORKDIR /project
COPY . .
RUN go build -o app .
FROM alpine
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
RUN update-ca-certificates
COPY --from=builder /project/app /app
ENTRYPOINT [ "/app" ]