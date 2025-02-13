# base image
FROM golang:1.22-alpine as base
WORKDIR /learn-go

ENV CGO_ENABLED=0

COPY go.mod go.sum /learn-go/
RUN go mod download

ADD . .
RUN go build -o /usr/local/bin/learn-go ./cmd/learn-go

# runner image
FROM gcr.io/distroless/static:latest
WORKDIR /app
COPY --from=base /usr/local/bin/learn-go learn-go

EXPOSE 5476
ENTRYPOINT ["/app/learn-go"]
