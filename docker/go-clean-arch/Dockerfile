# Builder
FROM golang:1.14.2-alpine3.11 as builder

RUN apk update && apk upgrade && \
    apk --update add git make

WORKDIR /app

COPY . .

RUN go build -o engine cmd/go-clean-arch/main.go

# Distribution
FROM alpine:latest

RUN apk update && apk upgrade && \
    apk --update --no-cache add tzdata && \
    mkdir -p /cmd/go-clean-arch 

WORKDIR /app 

EXPOSE 8080

COPY --from=builder /app/engine /cmd/go-clean-arch

ENV ENVIRONMENT="DEVELOP"
ENV DB_HOST="mysql"
ENV DB_PORT="3306"
ENV DB_USER="user"
ENV DB_PASS="password"
ENV DB_NAME="go_clean_arch"
ENV SECRET_KEY="sample_secret_key"

CMD /cmd/go-clean-arch/engine