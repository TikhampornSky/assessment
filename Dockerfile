FROM golang:1.19-alpine as build-base

WORKDIR /app

COPY go.mod .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go test -v

RUN go build -o ./out/go-sample .

FROM alpine:3.16.2
COPY --from=build-base /app/out/go-sample /app/go-sample

ENV DATABASE_URL postgres://mwabtxlk:x4mWGDcSX0VqkVEugDsAkXesZOAazEwF@tiny.db.elephantsql.com/mwabtxlk
ENV PORT :2565
EXPOSE 2565

CMD ["/app/go-sample"]