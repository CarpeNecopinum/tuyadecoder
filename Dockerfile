FROM golang:latest AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN go build -tags netgo .

FROM scratch

WORKDIR /app
COPY --from=builder /src/main /app/main
CMD ["/app/main"]
