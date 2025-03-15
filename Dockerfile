FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN apk add --no-cache --virtual .build-deps \
        gcc \
        g++ \
        make \
        bash &&  \
    go mod download

COPY . .

RUN make build

ENTRYPOINT ["make", "run"]
