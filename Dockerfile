FROM golang:1.22-alpine

# install psql
RUN apk update && apk add postgresql-client curl

# install goose
RUN curl -fsSL \
        https://raw.githubusercontent.com/pressly/goose/master/install.sh |\
        sh -s v3.19.2

WORKDIR /go/src/app

COPY . .

# make start.sh executable
RUN chmod +x start.sh

# build go app
RUN go mod download && go build -o banners-service ./cmd/server

CMD ["./banners-service"]