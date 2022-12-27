FROM golang:1.19-alpine

# Set working directory
WORKDIR /go/src/target
COPY go.mod ./
COPY go.sum ./
RUN go mod tidy
RUN go mod download
COPY . ./
RUN go build -o /docker-go

# Run tests
ENV DATABASE_URL postgres://mwabtxlk:x4mWGDcSX0VqkVEugDsAkXesZOAazEwF@tiny.db.elephantsql.com/mwabtxlk
ENV PORT :2565
EXPOSE 2565
CMD ["/docker-go"]
# CMD ["/docker-go", "CGO_ENABLED=0 go test --tags=integration ./..."]

# docker-compose -f docker-compose.yml up --build --abort-on-container-exit --exit-code-from go_test