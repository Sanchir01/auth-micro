FROM golang:1.24.3-alpine3.20 AS builder

WORKDIR /app

RUN apk --no-cache add bash git make gcc gettext musl-dev

RUN go version

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN make build


FROM alpine:3.20 AS runner

WORKDIR /app

COPY --from=builder /app/.bin/main /

COPY config config

COPY .env.prod .env.prod

ENV ENV_FILE=.env.prod

CMD ["/main"]
